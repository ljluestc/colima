# Fix SSH_AUTH_SOCK mapping inside VM

## Summary

Fixes #1330

When `--ssh-agent` is enabled, the `SSH_AUTH_SOCK` environment variable is set inside the VM but the socket path doesn't actually exist. This PR fixes the issue by using Lima's socket forwarding to map the host's SSH agent socket to a fixed path inside the VM.

## Problem

Previously, when starting colima with `--ssh-agent`:
```bash
colima start --ssh-agent
```

The VM would have `SSH_AUTH_SOCK` set to a path like `/tmp/ssh-yhBbSVSpFi/agent.1121`, but this socket didn't exist:
```bash
$ colima ssh env | grep SSH_AUTH_SOCK
SSH_AUTH_SOCK=/tmp/ssh-yhBbSVSpFi/agent.1121

$ colima ssh -- ls -la /tmp/ssh-yhBbSVSpFi/agent.1121
ls: /tmp/ssh-yhBbSVSpFi/agent.1121: No such file or directory
```

This made it impossible to use SSH agent forwarding inside Docker containers.

## Solution

When `ForwardAgent` is enabled and `SSH_AUTH_SOCK` is set on the host:

1. Add a Lima `portForward` entry to map the host's `SSH_AUTH_SOCK` to `/run/host-services/ssh-auth.sock` inside the VM
2. Set the `SSH_AUTH_SOCK` environment variable in the VM to point to this path
3. Add a provision script to ensure the `/run/host-services` directory exists

## After this fix

```bash
$ colima ssh env | grep SSH_AUTH_SOCK
SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock

$ colima ssh -- ls -la /run/host-services/ssh-auth.sock
srwxr-xr-x 1 root root 0 Dec 27 23:30 /run/host-services/ssh-auth.sock
```

Docker containers can now mount the SSH agent socket:
```bash
docker run --rm -v /run/host-services/ssh-auth.sock:/run/host-services/ssh-auth.sock \
  -e SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock \
  alpine ssh-add -l
```

## Changes

- `environment/vm/lima/limaconfig/config.go`: Add `Reverse` field to `PortForward` struct to enable host-to-guest socket forwarding
- `environment/vm/lima/yaml.go`: Add SSH agent socket forwarding via Lima's `portForwards` with `Reverse: true` when `ForwardAgent` is enabled
- `environment/vm/lima/yaml_test.go`: Add comprehensive tests for SSH agent socket forwarding

## Testing

- All existing tests pass
- New tests cover:
  - ForwardAgent enabled with SSH_AUTH_SOCK set (socket forwarding configured)
  - ForwardAgent enabled without SSH_AUTH_SOCK (no socket forwarding)
  - ForwardAgent disabled with SSH_AUTH_SOCK set (no socket forwarding)
  - ForwardAgent disabled without SSH_AUTH_SOCK (no socket forwarding)
