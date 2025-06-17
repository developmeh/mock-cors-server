# Dummy HTTP Mock CORS Server - User Guide

This comprehensive guide provides examples and scenarios for using the Dummy HTTP Mock CORS Server CLI tool.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Configuration Methods](#configuration-methods)
3. [Route Types and Examples](#route-types-and-examples)
4. [CORS Configuration Scenarios](#cors-configuration-scenarios)
5. [Common Use Cases](#common-use-cases)
6. [Advanced Examples](#advanced-examples)
7. [Troubleshooting](#troubleshooting)

## Quick Start

### Basic Usage

Start the server with default settings:
```bash
mock-cors-server
```

Start with a specific port:
```bash
mock-cors-server --port 8080
```

Start with a custom config file:
```bash
mock-cors-server --config /path/to/my-config.yaml
```

## Configuration Methods

The server supports three configuration methods (in order of precedence):

### 1. Command-Line Flags
```bash
# Override port via command line
mock-cors-server --port 9000

# Use custom config file
mock-cors-server --config ./my-custom-config.yaml

# Combine flags
mock-cors-server --config ./config.yaml --port 8080
```

### 2. Environment Variables
```bash
# Set port via environment variable
export MOCK_CORS_PORT=8080
mock-cors-server

# Multiple environment variables
export MOCK_CORS_PORT=8080
export MOCK_CORS_CORS_ALLOW_ORIGINS="https://example.com,https://test.com"
mock-cors-server
```

### 3. Configuration File
The server looks for `config.yaml` in these locations (in order):
- Current directory (`./config.yaml`)
- Home directory (`$HOME/.dummy_http_passkeys/config.yaml`)
- System directory (`/etc/dummy_http_passkeys/config.yaml`)

## Route Types and Examples

### 1. Dummy Routes (Hardcoded JSON Responses)

Perfect for mocking API endpoints with predefined responses.

```yaml
routes:
  - path: "/v1/json/begin"
    type: "dummy"
    content_type: "application/json"
    # Returns hardcoded JSON response for passkeys authentication
```

**Example Usage:**
```bash
curl http://localhost:8081/v1/json/begin
```

**Use Cases:**
- Mock authentication endpoints
- Simulate API responses during development
- Testing frontend applications

### 2. Static File Routes

Serve static files from the filesystem.

```yaml
routes:
  - path: "/static/example.html"
    type: "static"
    file_path: "./static/example.html"
    # Content type auto-detected from file extension
  
  - path: "/assets/logo.png"
    type: "static"
    file_path: "./assets/logo.png"
  
  - path: "/docs/api.txt"
    type: "static"
    file_path: "./documentation/api.txt"
    content_type: "text/plain"  # Override auto-detection
```

**Example Usage:**
```bash
curl http://localhost:8081/static/example.html
curl http://localhost:8081/assets/logo.png
```

**Use Cases:**
- Serve HTML pages for testing
- Provide static assets (images, CSS, JS)
- Serve documentation files

### 3. JSON Blob Routes

Return custom JSON responses defined in the configuration.

```yaml
routes:
  - path: "/api/custom/response"
    type: "json"
    json_content: '{"message": "Hello from JSON blob", "status": "ok", "data": {"key": "value"}}'
    content_type: "application/json"
  
  - path: "/api/user/profile"
    type: "json"
    json_content: '{"id": 123, "name": "John Doe", "email": "john@example.com", "roles": ["user", "admin"]}'
    content_type: "application/json"
```

**Example Usage:**
```bash
curl http://localhost:8081/api/custom/response
curl http://localhost:8081/api/user/profile
```

**Use Cases:**
- Mock API responses with custom data
- Test different response formats
- Simulate various API states

## CORS Configuration Scenarios

### Global CORS Settings

Apply the same CORS settings to all routes:

```yaml
cors:
  allow_origins:
    - "*"  # Allow all origins (development only)
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  allow_credentials: true
  max_age: 86400  # 24 hours
```

### Restrictive CORS for Production-like Testing

```yaml
cors:
  allow_origins:
    - "https://myapp.com"
    - "https://staging.myapp.com"
  allow_methods:
    - "GET"
    - "POST"
    - "OPTIONS"
  allow_headers:
    - "Content-Type"
    - "Authorization"
  allow_credentials: true
  max_age: 3600  # 1 hour
```

### Per-Route CORS Override

Different CORS settings for specific routes:

```yaml
# Global CORS settings
cors:
  allow_origins: ["*"]
  allow_methods: ["GET", "POST", "OPTIONS"]
  allow_headers: ["Content-Type"]
  allow_credentials: false
  max_age: 3600

routes:
  # This route uses global CORS settings
  - path: "/public/api"
    type: "json"
    json_content: '{"public": true}'
  
  # This route overrides CORS settings
  - path: "/secure/api"
    type: "json"
    json_content: '{"secure": true}'
    cors:
      allow_origins:
        - "https://secure.example.com"
      allow_methods:
        - "POST"
      allow_headers:
        - "Content-Type"
        - "Authorization"
        - "X-API-Key"
      allow_credentials: true
      max_age: 1800  # 30 minutes
```

## Common Use Cases

### 1. Passkeys/WebAuthn Development

Configuration for testing passkeys authentication:

```yaml
version: "1.0.0"
port: 8081

cors:
  allow_origins:
    - "https://localhost:3000"
    - "https://127.0.0.1:3000"
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

routes:
  - path: "/v1/json/begin"
    type: "dummy"
    content_type: "application/json"
  
  - path: "/v1/json/finish"
    type: "dummy"
    content_type: "application/json"
  
  - path: "/v1/json/register"
    type: "json"
    json_content: '{"challenge": "mock-challenge", "user": {"id": "user123", "name": "testuser"}}'
    content_type: "application/json"
```

**Usage:**
```bash
mock-cors-server --config passkeys-config.yaml
```

### 2. Frontend Development Mock API

```yaml
version: "1.0.0"
port: 3001

cors:
  allow_origins:
    - "http://localhost:3000"
    - "http://127.0.0.1:3000"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  allow_credentials: true
  max_age: 86400

routes:
  # User endpoints
  - path: "/api/users"
    type: "json"
    json_content: '[{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]'
  
  - path: "/api/users/1"
    type: "json"
    json_content: '{"id": 1, "name": "John Doe", "email": "john@example.com"}'
  
  # Product endpoints
  - path: "/api/products"
    type: "json"
    json_content: '[{"id": 1, "name": "Product A", "price": 99.99}, {"id": 2, "name": "Product B", "price": 149.99}]'
  
  # Static assets
  - path: "/images/logo.png"
    type: "static"
    file_path: "./assets/logo.png"
```

### 3. Testing CORS Policies

```yaml
version: "1.0.0"
port: 8082

# Strict CORS for testing
cors:
  allow_origins:
    - "https://trusted-domain.com"
  allow_methods:
    - "GET"
    - "POST"
  allow_headers:
    - "Content-Type"
  allow_credentials: false
  max_age: 300

routes:
  - path: "/api/strict"
    type: "json"
    json_content: '{"message": "This endpoint has strict CORS"}'
  
  # Override with more permissive CORS for specific endpoint
  - path: "/api/permissive"
    type: "json"
    json_content: '{"message": "This endpoint allows more origins"}'
    cors:
      allow_origins:
        - "https://trusted-domain.com"
        - "https://dev-domain.com"
        - "http://localhost:3000"
      allow_methods:
        - "GET"
        - "POST"
        - "PUT"
        - "DELETE"
        - "OPTIONS"
      allow_headers:
        - "Content-Type"
        - "Authorization"
        - "X-Custom-Header"
      allow_credentials: true
      max_age: 3600
```

## Advanced Examples

### Multi-Environment Configuration

Create different configs for different environments:

**development.yaml:**
```yaml
version: "1.0.0"
port: 8081

cors:
  allow_origins: ["*"]
  allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allow_headers: ["*"]
  allow_credentials: true
  max_age: 86400

routes:
  - path: "/api/debug"
    type: "json"
    json_content: '{"environment": "development", "debug": true}'
```

**staging.yaml:**
```yaml
version: "1.0.0"
port: 8081

cors:
  allow_origins:
    - "https://staging.myapp.com"
    - "https://staging-admin.myapp.com"
  allow_methods: ["GET", "POST", "OPTIONS"]
  allow_headers: ["Content-Type", "Authorization"]
  allow_credentials: true
  max_age: 3600

routes:
  - path: "/api/status"
    type: "json"
    json_content: '{"environment": "staging", "debug": false}'
```

**Usage:**
```bash
# Development
mock-cors-server --config development.yaml

# Staging
mock-cors-server --config staging.yaml
```

### Complex Route Configuration

```yaml
version: "1.0.0"
port: 8081

cors:
  allow_origins: ["*"]
  allow_methods: ["GET", "POST", "OPTIONS"]
  allow_headers: ["Content-Type", "Authorization"]
  allow_credentials: true
  max_age: 86400

routes:
  # Health check endpoint
  - path: "/health"
    type: "json"
    json_content: '{"status": "healthy", "timestamp": "2024-01-01T00:00:00Z"}'
  
  # API versioning
  - path: "/api/v1/users"
    type: "json"
    json_content: '{"version": "v1", "users": [{"id": 1, "name": "User1"}]}'
  
  - path: "/api/v2/users"
    type: "json"
    json_content: '{"version": "v2", "users": [{"id": 1, "name": "User1", "email": "user1@example.com"}]}'
  
  # File downloads
  - path: "/downloads/sample.pdf"
    type: "static"
    file_path: "./files/sample.pdf"
    content_type: "application/pdf"
  
  # Custom headers for specific route
  - path: "/api/special"
    type: "json"
    json_content: '{"special": true}'
    cors:
      allow_origins: ["https://special.example.com"]
      allow_methods: ["POST"]
      allow_headers: ["Content-Type", "X-Special-Token"]
      allow_credentials: true
      max_age: 1800
```

## Troubleshooting

### Common Issues and Solutions

#### 1. CORS Errors in Browser
**Problem:** Browser shows CORS errors when making requests.

**Solution:** Check your `allow_origins` configuration:
```yaml
cors:
  allow_origins:
    - "http://localhost:3000"  # Add your frontend URL
    - "https://yourdomain.com"
```

#### 2. Port Already in Use
**Problem:** Error "port already in use" when starting server.

**Solutions:**
```bash
# Use a different port
mock-cors-server --port 8082

# Or kill the process using the port
lsof -ti:8081 | xargs kill -9
```

#### 3. Config File Not Found
**Problem:** Server can't find configuration file.

**Solutions:**
```bash
# Specify config file explicitly
mock-cors-server --config ./config.yaml

# Or place config.yaml in current directory
cp config.yaml.sample config.yaml
```

#### 4. Static Files Not Served
**Problem:** Static files return 404 errors.

**Solution:** Check file paths in configuration:
```yaml
routes:
  - path: "/static/file.html"
    type: "static"
    file_path: "./static/file.html"  # Ensure this path exists
```

### Debug Mode

To see which config file is being used:
```bash
mock-cors-server 2>&1 | grep "Using config file"
```

### Testing CORS Settings

Use curl to test CORS headers:
```bash
# Test preflight request
curl -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  http://localhost:8081/api/test

# Test actual request
curl -X POST \
  -H "Origin: https://example.com" \
  -H "Content-Type: application/json" \
  http://localhost:8081/api/test
```

### Environment Variable Reference

All configuration options can be set via environment variables with the `MOCK_CORS_` prefix:

```bash
export MOCK_CORS_PORT=8081
export MOCK_CORS_CORS_ALLOW_ORIGINS="https://example.com,https://test.com"
export MOCK_CORS_CORS_ALLOW_METHODS="GET,POST,OPTIONS"
export MOCK_CORS_CORS_ALLOW_HEADERS="Content-Type,Authorization"
export MOCK_CORS_CORS_ALLOW_CREDENTIALS=true
export MOCK_CORS_CORS_MAX_AGE=3600
```

## Getting Help

For more information:
- Check the main [README.md](README.md) for installation and basic setup
- View the sample configuration: [config.yaml.sample](config.yaml.sample)
- Report issues on the project's GitHub repository

---

*This user guide covers the most common scenarios. For advanced use cases or custom requirements, refer to the source code or create an issue on the project repository.*