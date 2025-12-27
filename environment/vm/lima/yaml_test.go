package lima

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/abiosoft/colima/config"
	"github.com/abiosoft/colima/environment/vm/lima/limaconfig"
	"github.com/abiosoft/colima/util"
	"github.com/abiosoft/colima/util/fsutil"
)

func Test_checkOverlappingMounts(t *testing.T) {
	type args struct {
		mounts []string
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{args: args{mounts: []string{"/User", "/User/something"}}, wantErr: true},
		{args: args{mounts: []string{"/User/one", "/User/two"}}, wantErr: false},
		{args: args{mounts: []string{"/User/one", "/User/one_other"}}, wantErr: false},
		{args: args{mounts: []string{"/User/one_other", "/User/one"}}, wantErr: false},
		{args: args{mounts: []string{"/User/one", "/User/one/other"}}, wantErr: true},
		{args: args{mounts: []string{"/User/one/", "/User/one"}}, wantErr: true},
		{args: args{mounts: []string{"/User/one/", "/User/two", "User/one"}}, wantErr: true},
		{args: args{mounts: []string{"/home/a/b/c", "/home/b/c/a", "/home/c/a/b"}}, wantErr: false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			mounts := func(mounts []string) (mnts []config.Mount) {
				for _, m := range mounts {
					mnts = append(mnts, config.Mount{Location: m})
				}
				return
			}(tt.args.mounts)
			if err := checkOverlappingMounts(mounts); (err != nil) != tt.wantErr {
				t.Errorf("checkOverlappingMounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_config_Mounts(t *testing.T) {
	fsutil.FS = fsutil.FakeFS
	tests := []struct {
		mounts        []string
		isDefault     bool
		includesCache bool
	}{
		{mounts: []string{"/User/user", "/tmp/another"}},
		{mounts: []string{"/User/another", "/User/something", "/User/else"}},
		{isDefault: true},
		{mounts: []string{util.HomeDir()}, includesCache: true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			mounts := func(mounts []string) (mnts []config.Mount) {
				for _, m := range mounts {
					mnts = append(mnts, config.Mount{Location: m})
				}
				return
			}(tt.mounts)
			conf, err := newConf(context.Background(), config.Config{Mounts: mounts})
			if err != nil {
				t.Error(err)
				return
			}

			expectedLocations := tt.mounts
			if tt.isDefault {
				expectedLocations = []string{"~", "/tmp/colima"}
			} else if !tt.includesCache {
				expectedLocations = append([]string{config.CacheDir()}, tt.mounts...)
			}

			sameMounts := func(expectedLocations []string, mounts []limaconfig.Mount) bool {
				sanitize := func(s string) string { return strings.TrimSuffix(s, "/") + "/" }
				for i, m := range mounts {
					if sanitize(m.Location) != sanitize(expectedLocations[i]) {
						return false
					}
				}
				return true
			}(expectedLocations, conf.Mounts)
			if !sameMounts {
				foundLocations := func() (locations []string) {
					for _, m := range conf.Mounts {
						locations = append(locations, m.Location)
					}
					return
				}()
				t.Errorf("got: %+v, want: %v", foundLocations, expectedLocations)
			}
		})
	}
}

func Test_ingressDisabled(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{args: []string{"--flag=f", "--another", "flag"}, want: false},
		{args: []string{"--disable=traefik", "--version=3"}, want: true},
		{args: []string{}, want: false},
		{args: []string{"--disable", "traefik", "--one=two"}, want: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			if got := ingressDisabled(tt.args); got != tt.want {
				t.Errorf("ingressDisabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sshAgentForwarding(t *testing.T) {
	fsutil.FS = fsutil.FakeFS

	tests := []struct {
		name             string
		forwardAgent     bool
		hostSSHAuthSock  string
		wantEnvSet       bool
		wantPortForward  bool
		wantProvision    bool
		wantGuestSocket  string
	}{
		{
			name:             "ForwardAgent enabled with SSH_AUTH_SOCK set",
			forwardAgent:     true,
			hostSSHAuthSock:  "/tmp/ssh-agent.sock",
			wantEnvSet:       true,
			wantPortForward:  true,
			wantProvision:    true,
			wantGuestSocket:  "/run/host-services/ssh-auth.sock",
		},
		{
			name:             "ForwardAgent enabled without SSH_AUTH_SOCK",
			forwardAgent:     true,
			hostSSHAuthSock:  "",
			wantEnvSet:       false,
			wantPortForward:  false,
			wantProvision:    false,
		},
		{
			name:             "ForwardAgent disabled with SSH_AUTH_SOCK set",
			forwardAgent:     false,
			hostSSHAuthSock:  "/tmp/ssh-agent.sock",
			wantEnvSet:       false,
			wantPortForward:  false,
			wantProvision:    false,
		},
		{
			name:             "ForwardAgent disabled without SSH_AUTH_SOCK",
			forwardAgent:     false,
			hostSSHAuthSock:  "",
			wantEnvSet:       false,
			wantPortForward:  false,
			wantProvision:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset SSH_AUTH_SOCK environment variable
			// t.Setenv automatically restores the original value after the test
			if tt.hostSSHAuthSock != "" {
				t.Setenv("SSH_AUTH_SOCK", tt.hostSSHAuthSock)
			} else {
				t.Setenv("SSH_AUTH_SOCK", "")
			}

			conf, err := newConf(context.Background(), config.Config{
				ForwardAgent: tt.forwardAgent,
			})
			if err != nil {
				t.Fatalf("newConf() error = %v", err)
			}

			// Check SSH_AUTH_SOCK env var in config
			gotEnvVal, gotEnvSet := conf.Env["SSH_AUTH_SOCK"]
			if gotEnvSet != tt.wantEnvSet {
				t.Errorf("SSH_AUTH_SOCK env set = %v, want %v", gotEnvSet, tt.wantEnvSet)
			}
			if tt.wantEnvSet && gotEnvVal != tt.wantGuestSocket {
				t.Errorf("SSH_AUTH_SOCK env value = %v, want %v", gotEnvVal, tt.wantGuestSocket)
			}

			// Check port forward for SSH agent socket
			gotPortForward := false
			const expectedGuestSocket = "/run/host-services/ssh-auth.sock"
			for _, pf := range conf.PortForwards {
				if pf.GuestSocket == expectedGuestSocket {
					gotPortForward = true
					if tt.wantPortForward {
						if pf.HostSocket != tt.hostSSHAuthSock {
							t.Errorf("SSH agent port forward HostSocket = %v, want %v", pf.HostSocket, tt.hostSSHAuthSock)
						}
						if !pf.Reverse {
							t.Errorf("SSH agent port forward Reverse = %v, want true", pf.Reverse)
						}
					}
					break
				}
			}
			if gotPortForward != tt.wantPortForward {
				t.Errorf("SSH agent port forward present = %v, want %v", gotPortForward, tt.wantPortForward)
			}

			// Check provision script for creating /run/host-services directory
			gotProvision := false
			for _, p := range conf.Provision {
				if strings.Contains(p.Script, "mkdir -p /run/host-services") {
					gotProvision = true
					break
				}
			}
			if gotProvision != tt.wantProvision {
				t.Errorf("SSH agent provision script present = %v, want %v", gotProvision, tt.wantProvision)
			}
		})
	}
}
