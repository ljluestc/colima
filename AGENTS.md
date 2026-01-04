# AGENTS.md - Colima Project

This file provides guidance to AI agents working on the Colima project.

## Project Overview

Colima is a container runtime for macOS with minimal setup, supporting Docker and Kubernetes. It provides:

- **Language**: Go (Golang) 1.23+
- **Type**: Container runtime / CLI application
- **Purpose**: Lightweight container runtime for macOS
- **Features**: Docker and Kubernetes support, minimal setup, good performance

## Key Configuration Files

- `go.mod` - Go module definition and dependencies
- `go.sum` - Dependency checksums
- `Makefile` - Build automation scripts
- `.golangci.yml` - GolangCI-Lint configuration
- `README.md` - Project documentation

## Build and Test Commands

### Installation
```bash
# Install dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Development
```bash
# Build the project
make build

# Build for specific platform
make build-darwin
make build-linux

# Install locally
make install
```

### Testing
```bash
# Run tests
make test

# Run specific test
go test ./path/to/package

# Linting
golangci-lint run

# Format code
gofmt -w .
```

### Makefile Commands
```bash
# Common make commands
make build      # Build the project
make test       # Run test suite
make install    # Install locally
make clean      # Clean build artifacts
make lint       # Run linter
make format     # Format code
```

## Project Structure

```
cmd/                    # Command-line applications
  colima/               # Main Colima CLI
app/                    # Application core
cli/                    # CLI utilities
config/                 # Configuration management
core/                   # Core functionality
daemon/                 # Daemon processes
environment/            # Environment management
store/                  # Storage management
util/                   # Utilities
docs/                   # Documentation
scripts/                # Build and utility scripts
embedded/               # Embedded resources
integration/            # Integration components
```

## Code Style Guidelines

- **Go Standards**: Follow official Go code review comments
- **Formatting**: Use `gofmt` for consistent formatting
- **Linting**: GolangCI-Lint configuration in `.golangci.yml`
- **Naming**: Use camelCase for variables, PascalCase for exported types
- **Error Handling**: Explicit error handling (no panic for expected errors)
- **Documentation**: Add godoc comments for exported functions/types
- **CLI Design**: Follow good CLI design principles

## Testing Instructions

- **Unit Tests**: Located alongside source files
- **Integration Tests**: In `integration/` directory
- **Mocking**: Use interfaces for mocking dependencies
- **Test Coverage**: Aim for high coverage of core functionality
- **Platform Testing**: Test on macOS (primary) and Linux

## Security Considerations

- **Container Security**: Secure container runtime isolation
- **Privilege Escalation**: Minimize required privileges
- **Input Validation**: Validate all user inputs
- **Error Messages**: Don't expose sensitive information
- **Dependency Management**: Regularly update dependencies
- **Network Security**: Secure network communications

## Performance Considerations

- **Resource Usage**: Minimize memory and CPU usage
- **Startup Time**: Optimize application startup
- **Container Operations**: Efficient container management
- **I/O Operations**: Optimize filesystem operations
- **Concurrency**: Use goroutines effectively

## Container Runtime Features

- **Docker Support**: Full Docker compatibility
- **Kubernetes Support**: Kubernetes cluster management
- **Resource Management**: CPU and memory limits
- **Networking**: Container networking support
- **Volume Management**: Persistent storage support
- **Cross-Platform**: macOS and Linux support

## CLI Design

- **User-Friendly**: Intuitive command interface
- **Help System**: Comprehensive help and documentation
- **Error Handling**: Clear error messages
- **Configuration**: Flexible configuration options
- **Output Formatting**: Support for different output formats

## Git Conventions

- **Commit Messages**: Clear, descriptive commit messages
- **Branching**: Use feature branches for new development
- **Pull Requests**: Required for merging to main branch
- **Tags**: Use semantic versioning for releases
- **Changelog**: Maintain changelog for significant changes

## CI/CD

- **GitHub Actions**: Configured in `.github/workflows/`
- **Automated Testing**: Runs on every push/PR
- **Build Verification**: Ensures project builds successfully
- **Linting Checks**: Code quality gates
- **Release Process**: Automated release workflows
- **Cross-Platform**: Build and test on multiple platforms

## Documentation

- **README.md**: Main project documentation
- **docs/`: Additional documentation
- **Godoc**: Use godoc comments for inline documentation
- **Examples**: Include practical usage examples
- **Troubleshooting**: Common issues and solutions

## Dependency Management

- **Go Modules**: Uses Go modules for dependency management
- **Minimal Dependencies**: Keep dependencies minimal
- **Updates**: Regularly update dependencies
- **Compatibility**: Ensure backward compatibility

## macOS-Specific Considerations

- **System Integration**: macOS-specific integrations
- **Permissions**: Handle macOS permission model
- **Performance**: Optimize for macOS performance
- **UI Integration**: Consider macOS UI guidelines
- **Installation**: macOS package management

## Container Ecosystem Integration

- **Docker CLI**: Full Docker CLI compatibility
- **Kubernetes CLI**: kubectl integration
- **Containerd**: Containerd support
- **OCI Standards**: OCI compliance
- **CNCF Ecosystem**: Integration with CNCF tools

## Future Enhancements

- **Additional Features**: Support for more container features
- **Performance**: Optimize container operations
- **Monitoring**: Add metrics and monitoring support
- **Documentation**: Expand usage examples and tutorials
- **Platform Support**: Additional platform support

## Community and Support

- **Issue Tracking**: GitHub issues for bug reports
- **Contributions**: Welcome community contributions
- **Documentation**: Improve and expand documentation
- **Support**: Provide user support channels
- **Roadmap**: Clear project roadmap and vision

## Task Implementation
1. **Analyze Requirements**: Refer to `README.md` for detailed feature specifications and system design.
2. **Implementation**: Modify source code in the respective directories (e.g., `src/`, `internal/`).
3. **Verification**: Run provided build and test commands (see above) to ensure correctness.
4. **Push Changes**:
   - Commit changes: `git commit -m "feat: implement <feature>"`
   - Push to remote: `git push origin <branch-name>`
