# colima — Agent Guide

## Project overview

This repository is a Go implementation of **Colima** ("Containers on Lima"), a CLI that provides container runtimes on macOS (and Linux) with minimal setup.

From the source tree:

- The CLI is built with `spf13/cobra` (`cmd/root/root.go`).
- The main entrypoint is `cmd/colima/main.go`.
- The build pipeline is driven by the repo `Makefile`, producing binaries under `_output/binaries/`.
- CI is implemented with GitHub Actions (`.github/workflows/go.yml` and `.github/workflows/integration.yml`).

## Repository layout

- `cmd/`
  - CLI commands (Cobra) and the binary entrypoint.
  - `cmd/colima/main.go` imports the command tree and calls `root.Execute()`.
  - `cmd/root/root.go` defines the root command, global flags, and profile selection behavior.
- `config/`
  - Version metadata is injected at build time via ldflags in the `Makefile` (`github.com/abiosoft/colima/config`).
- `scripts/`
  - Helper scripts for integration testing and optional components (e.g. `scripts/integration.sh`, `scripts/build_vmnet.sh`).
- `embedded/`
  - Embedded assets and helper scripts such as `embedded/images/images_sha.sh` (referenced by `make images-sha`).
- `.github/workflows/`
  - Build/test workflows and macOS integration workflows.

## Key configuration files

- `go.mod`
  - Go module: `github.com/abiosoft/colima`.
- `Makefile`
  - Authoritative developer entrypoint for building, testing, installing, formatting, and integration tests.
- `.github/workflows/go.yml`
  - CI for build/test and release artifact generation.
- `.github/workflows/integration.yml`
  - CI integration runs on macOS and exercises `colima start` with Docker/containerd and optional Kubernetes.

## Build and test commands

### Build

The default target builds the `colima` binary into `_output/binaries/`:

```bash
make
```

Under the hood this runs a platform-specific build (see `Makefile`):

- Output binary name: `colima-$(OS)-$(ARCH)`
- Output directory: `_output/binaries/`
- `ldflags` inject version metadata into `github.com/abiosoft/colima/config`.

### Install

Installs the built binary into `/usr/local/bin/colima`:

```bash
sudo make install
```

### Unit tests

```bash
make test
```

Or directly:

```bash
go test -v ./...
```

### Formatting and linting

- Format:

```bash
make fmt
```

- Lint (requires `golangci-lint` in PATH):

```bash
make lint
```

### Integration tests

The `integration` target builds first and then runs `scripts/integration.sh` using the built binary path:

```bash
make integration
```

## Runtime / behavior notes

### Profiles

The root command supports a `--profile/-p` flag (default `default`). Some commands also accept the profile as a positional argument when `--profile` is not set (see `cmd/root/root.go`).

### Supported runtimes

Per `README.md`, Colima supports multiple runtimes (Docker, containerd, Incus) and optional Kubernetes.

## CI / release

- `go.yml` runs build/test and also produces binaries for Linux and macOS via `make` and uploads them as workflow artifacts.
- `integration.yml` runs integration scenarios on macOS runners. It installs CLI dependencies using Homebrew and validates Kubernetes/Docker startup via `colima start ...`.

## Security / operational considerations

- Colima manages local VMs/container runtimes and may require elevated privileges for install steps (e.g. writing to `/usr/local/bin`).
- The integration workflow installs and uses system tools (`kubectl`, `docker`, `lima`), so running integration locally can affect local Docker/VM state.

## Working conventions for agents

- Prefer `make` targets over ad-hoc commands.
- When changing CLI behavior, locate the relevant Cobra command under `cmd/` and confirm how profile selection and logging initialization behave in `cmd/root/root.go`.
- Keep changes compatible with both macOS and Linux (CI builds both).
