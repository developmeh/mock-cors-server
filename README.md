# Dummy HTTP Mock CORS Server

A configurable HTTP server written in Go that handles mock CORS requests. This server is designed for testing and development purposes, providing configurable routes, CORS settings, and dummy responses.

## 📖 User Guide

**👉 [Complete User Guide with Examples](USER_GUIDE.md)** - Comprehensive guide with examples for various scenarios, configuration methods, and troubleshooting tips.

## Features

- **Configurable Routes**: Support for multiple routes with individual settings
- **Flexible CORS**: Global and per-route CORS configuration
- **CLI Interface**: Built with Cobra for easy command-line usage
- **Configuration Management**: Support for config files, environment variables, and CLI flags
- **Multiple Platforms**: Cross-platform builds for Linux, macOS, and Windows
- **Comprehensive Testing**: Unit tests and end-to-end testing scripts
- **CI/CD Ready**: GitHub Actions workflows for testing and releases

## Requirements

- Go 1.19 or higher

## Installation

### From Source

1. Clone this repository:
   ```bash
   git clone https://github.com/pscarrone/dummy_http_passkeys.git
   cd dummy_http_passkeys
   ```

2. Build the application:
   ```bash
   make build
   ```

### From Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/pscarrone/dummy_http_passkeys/releases).

### Using Homebrew (macOS/Linux)

```bash
brew install pscarrone/tap/dummy-http-passkeys
```

## Configuration

The server can be configured in three ways (in order of precedence):

1. **Command-line flags**
2. **Environment variables** (prefixed with `MOCK_CORS_`)
3. **Configuration file** (`config.yaml`)

### Configuration File

Create a `config.yaml` file in the current directory, `$HOME/.dummy_http_passkeys/`, or `/etc/dummy_http_passkeys/`:

```yaml
version: "1.0.0"
port: 8081

# Global CORS settings
cors:
  allow_origins:
    - "*"
  allow_methods:
    - "GET"
    - "POST"
    - "OPTIONS"
  allow_headers:
    - "Content-Type"
    - "Authorization"
    - "site-token"
    - "client-id"
    - "placement-id"
    - "integrator-id"
    - "oauth-type"
  allow_credentials: true
  max_age: 86400

# Routes configuration
routes:
  - path: "/v1/json/begin"
    content_type: "application/json"
    # Uses global CORS settings

  - path: "/api/v2/test"
    content_type: "application/json"
    cors:
      allow_origins:
        - "https://example.com"
      allow_methods:
        - "POST"
      allow_headers:
        - "Content-Type"
      allow_credentials: false
      max_age: 3600
```

### Environment Variables

```bash
export MOCK_CORS_PORT=8081
export MOCK_CORS_CORS_ALLOW_ORIGINS="*"
```

### Command-line Flags

```bash
mock-cors-server --port 8081 --config /path/to/config.yaml
```

## Usage

### Using the Makefile

```bash
# Build the executable
make build

# Run the server
make run

# Build and run in one command (default)
make all

# Run unit tests
make test

# Run end-to-end tests
make test-e2e

# Run all tests
make test-all

# Clean up compiled binaries and test files
make clean

# Show available commands
make help
```

### Using the Binary Directly

```bash
# Start server with default settings
./mock_cors_server

# Start server on custom port
./mock_cors_server --port 9000

# Start server with custom config file
./mock_cors_server --config /path/to/config.yaml

# Show help
./mock_cors_server --help
```

## API Endpoints

### Default Route: POST /v1/json/begin

Returns a JSON object with dummy data for mock CORS testing.

#### Example Request:

```bash
curl -X POST http://localhost:8081/v1/json/begin \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{}'
```

#### Example Response:

```json
{
  "status": "success",
  "challenge": "abc123xyz789",
  "timestamp": "2025-06-17T13:31:43.306302Z",
  "expiresIn": 300,
  "sessionId": "12345-67890-abcde-fghij"
}
```

### CORS Preflight

The server supports CORS preflight requests:

```bash
curl -X OPTIONS http://localhost:8081/v1/json/begin \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type"
```

## Testing

### Unit Tests

```bash
make test
```

### End-to-End Tests

```bash
make test-e2e
```

The E2E test script will:
- Start the server in the background
- Run comprehensive tests (OPTIONS, POST, invalid methods)
- Validate JSON responses
- Clean up processes automatically

### Manual Testing

Test the server manually using curl:

```bash
# Test CORS preflight
curl -X OPTIONS http://localhost:8081/v1/json/begin \
  -H "Origin: http://localhost:3000" \
  -v

# Test POST request
curl -X POST http://localhost:8081/v1/json/begin \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{}' \
  -v

# Test invalid method (should return 405)
curl -X GET http://localhost:8081/v1/json/begin -v
```

## Development

### Project Structure

```
.
├── cmd/server/          # Main application entry point
├── pkg/server/          # Public server package
├── internal/config/     # Private configuration package
├── .github/workflows/   # GitHub Actions CI/CD
├── config.yaml         # Sample configuration file
├── test_e2e.sh         # End-to-end test script
├── Makefile            # Build and test commands
└── .goreleaser.yml     # Release configuration
```

### Building for Multiple Platforms

```bash
# Build for current platform
make build

# Build for all platforms (requires GoReleaser)
goreleaser build --snapshot --clean
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test-all`
5. Submit a pull request

## CI/CD

The project uses GitHub Actions for:

- **Continuous Integration**: Run tests on multiple Go versions
- **Cross-platform builds**: Build binaries for Linux, macOS, and Windows
- **Automated releases**: Publish releases when tags are pushed
- **Homebrew formula**: Automatically update Homebrew tap

## License

This project is licensed under the MIT License - see the LICENSE file for details.
