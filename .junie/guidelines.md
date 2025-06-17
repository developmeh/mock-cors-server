# Dummy HTTP Passkeys Server Development Guidelines

This document provides guidelines for future development on the Dummy HTTP Passkeys Server project.

## Build/Configuration Instructions

### Makefile

The project uses a Makefile for common operations:

- `make build`: Compiles the server binary
- `make run`: Starts the server in the background and saves the PID for later termination
- `make stop`: Stops the running server using the saved PID
- `make test`: Runs the test script
- `make clean`: Removes compiled binaries and temporary files

When modifying the Makefile, ensure it:
- Starts the server in the background and pipes stdout to the shell
- Keeps track of the server PID for proper shutdown
- Provides clear error messages if operations fail

### GitHub Actions

The project should use GitHub Actions for CI/CD:
- Create a `.github/workflows` directory with workflow files for:
  - Running tests on pull requests
  - Building and testing on multiple platforms
  - Publishing releases when tags are pushed

### GoReleaser

The project should use GoReleaser for creating and publishing releases:
- Create a `.goreleaser.yml` configuration file
- Configure it to build binaries for multiple platforms
- Set it up to publish artifacts through GitHub Actions when tags are pushed
- Include checksums and signature files for verification

### Configuration

The server should be configurable with the following options:

- **Port**: The HTTP port should be configurable (currently hardcoded to 8081)
- **Routes**: The server should support multiple routes, each serving a static file
- **CORS/Preflight**: All preflight options and headers should be configurable:
  - Allow Origins
  - Allow Methods
  - Allow Headers
  - Allow Credentials
  - Max Age

Configuration should be implemented using Viper and Cobra for a consistent CLI experience:
- Command-line flags
- Environment variables
- Configuration files (JSON, YAML, etc.)

## Testing Information

### Test Coverage

All modules should have comprehensive tests:
- Unit tests for individual functions
- Integration tests for API endpoints
- End-to-end tests for the complete flow

### Running Tests

Always run tests before considering a request complete:
- Use `go test ./...` for running all tests
- Use `make test` for running the end-to-end test script

### E2E Testing

The project should include a shell script for end-to-end testing that:
- Starts the server in the background
- Runs tests against the running server
- Cleans up processes after testing
- Returns a non-zero exit code if tests fail
- Does not block the shell

### Process Cleanup

Always ensure that test processes are cleaned up:
- Use defer statements to ensure cleanup happens
- Track PIDs of spawned processes
- Use signal handlers to catch interrupts and clean up

## Additional Development Information

### Project structure
This is a golang project and should be organized like a standard
golang project
- use standard folder structure

### CLI Tool Design

This is a CLI tool, so always review changes with these considerations:
- Update Viper/Cobra configuration when adding new features
- Maintain backward compatibility when possible
- Document all CLI commands and options
- Provide sensible defaults for all options

### Command Deprecation

Try not to deprecate CLI commands:
- Mark commands as deprecated before removing them
- Provide migration paths for users
- Keep deprecated commands working for at least one major version

### Configuration Versioning

Configuration should be versioned:
- Include a version field in configuration files
- Provide migration tools for upgrading configurations
- Document breaking changes between versions

### Route Configuration

Routes should be enumerable, allowing the config to specify:
- Multiple routes that each serve a static file
- Global or individual preflight options for each route
- Content types and other response headers

### Security Considerations

When developing, keep security in mind:
- Validate all input parameters
- Use proper CORS settings to prevent unauthorized access
- Log security-relevant events
- Don't expose sensitive information in logs or error messages

## Implementation Roadmap

1. Update the server to use Viper/Cobra for configuration
2. Make the port configurable
3. Implement enumerable routes with individual preflight options
4. Add GitHub Actions workflows
5. Configure GoReleaser
6. Improve test coverage and create E2E test scripts
7. Document all changes in README.md