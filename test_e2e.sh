#!/bin/bash

# E2E Test Script for Dummy HTTP Mock CORS Server
# This script starts the server, runs tests, and cleans up processes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVER_BINARY="./bin/mock_cors_server"
TEST_PORT=8081
PID_FILE=".test_server.pid"
LOG_FILE=".test_server.log"

# Cleanup function
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"

    # Kill server if PID file exists
    if [ -f "$PID_FILE" ]; then
        SERVER_PID=$(cat "$PID_FILE")
        if kill -0 "$SERVER_PID" 2>/dev/null; then
            echo "Stopping server (PID: $SERVER_PID)"
            kill "$SERVER_PID" || true
            sleep 1
            # Force kill if still running
            if kill -0 "$SERVER_PID" 2>/dev/null; then
                kill -9 "$SERVER_PID" || true
            fi
        fi
        rm -f "$PID_FILE"
    fi

    # Clean up log file
    rm -f "$LOG_FILE"

    echo -e "${GREEN}Cleanup completed${NC}"
}

# Set up trap to ensure cleanup on exit
trap cleanup EXIT INT TERM

# Function to check if server is running
check_server() {
    local max_attempts=10
    local attempt=1

    echo "Waiting for server to start..."

    while [ $attempt -le $max_attempts ]; do
        if curl -s -o /dev/null -w "%{http_code}" "http://localhost:$TEST_PORT/v1/json/begin" -X OPTIONS > /dev/null 2>&1; then
            echo -e "${GREEN}Server is running!${NC}"
            return 0
        fi

        echo "Attempt $attempt/$max_attempts: Server not ready yet..."
        sleep 1
        attempt=$((attempt + 1))
    done

    echo -e "${RED}Server failed to start within expected time${NC}"
    return 1
}

# Function to run tests
run_tests() {
    echo -e "${YELLOW}Running E2E tests...${NC}"

    # Test 1: OPTIONS request (CORS preflight)
    echo "Test 1: OPTIONS request"
    response=$(curl -s -o /dev/null -w "%{http_code}" \
        -X OPTIONS \
        -H "Origin: http://localhost:3000" \
        -H "Access-Control-Request-Method: POST" \
        -H "Access-Control-Request-Headers: Content-Type, Authorization" \
        "http://localhost:$TEST_PORT/v1/json/begin")

    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✓ OPTIONS test passed${NC}"
    else
        echo -e "${RED}✗ OPTIONS test failed (HTTP $response)${NC}"
        return 1
    fi

    # Test 2: POST request
    echo "Test 2: POST request"
    response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -H "Origin: http://localhost:3000" \
        -d '{}' \
        "http://localhost:$TEST_PORT/v1/json/begin")

    # Extract HTTP status code (last 3 characters)
    http_code="${response: -3}"
    response_body="${response%???}"

    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ POST test passed${NC}"
        echo "Response body: $response_body"

        # Validate JSON response
        if echo "$response_body" | jq . > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Response is valid JSON${NC}"
        else
            echo -e "${RED}✗ Response is not valid JSON${NC}"
            return 1
        fi
    else
        echo -e "${RED}✗ POST test failed (HTTP $http_code)${NC}"
        echo "Response: $response_body"
        return 1
    fi

    # Test 3: Invalid method
    echo "Test 3: Invalid method (GET)"
    response=$(curl -s -o /dev/null -w "%{http_code}" \
        -X GET \
        "http://localhost:$TEST_PORT/v1/json/begin")

    if [ "$response" = "405" ]; then
        echo -e "${GREEN}✓ Invalid method test passed${NC}"
    else
        echo -e "${RED}✗ Invalid method test failed (HTTP $response, expected 405)${NC}"
        return 1
    fi

    echo -e "${GREEN}All tests passed!${NC}"
    return 0
}

# Main execution
main() {
    echo -e "${YELLOW}Starting E2E tests for Dummy HTTP Mock CORS Server${NC}"

    # Check if binary exists
    if [ ! -f "$SERVER_BINARY" ]; then
        echo -e "${RED}Server binary not found: $SERVER_BINARY${NC}"
        echo "Please run 'make build' first"
        exit 1
    fi

    # Check if required tools are available
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}curl is required but not installed${NC}"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}jq not found, JSON validation will be skipped${NC}"
    fi

    # Start server in background
    echo "Starting server..."
    "$SERVER_BINARY" --port "$TEST_PORT" > "$LOG_FILE" 2>&1 &
    SERVER_PID=$!
    echo $SERVER_PID > "$PID_FILE"

    echo "Server started with PID: $SERVER_PID"

    # Check if server started successfully
    if ! check_server; then
        echo -e "${RED}Failed to start server${NC}"
        echo "Server log:"
        cat "$LOG_FILE" || true
        exit 1
    fi

    # Run tests
    if run_tests; then
        echo -e "${GREEN}E2E tests completed successfully!${NC}"
        exit 0
    else
        echo -e "${RED}E2E tests failed!${NC}"
        echo "Server log:"
        cat "$LOG_FILE" || true
        exit 1
    fi
}

# Run main function
main "$@"
