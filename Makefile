# Makefile for Dummy HTTP Mock CORS Server

# Variables
BINARY_NAME=mock_cors_server
BINARY_PATH=bin/$(BINARY_NAME)
GO=go
PID_FILE=.server.pid

# Default target
.PHONY: all
all: build run

# Build the application
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p bin
	$(GO) build -o $(BINARY_PATH) ./cmd/server

# Run the application
.PHONY: run
run: build
	@echo "Running server..."
	./$(BINARY_PATH) & echo $$! > $(PID_FILE)
	@echo "Server started with PID: $$(cat $(PID_FILE))"

# Build and run in one command
.PHONY: build-run
build-run: build run

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	$(GO) test -v ./...

# Run end-to-end tests
.PHONY: test-e2e
test-e2e: build
	@echo "Running end-to-end tests..."
	./test_e2e.sh

# Run all tests
.PHONY: test-all
test-all: test test-e2e

# Stop the running server
.PHONY: stop
stop:
	@if [ -f $(PID_FILE) ]; then \
		echo "Stopping server with PID: $$(cat $(PID_FILE))"; \
		kill -9 $$(cat $(PID_FILE)) || true; \
		rm -f $(PID_FILE); \
		echo "Server stopped"; \
	else \
		echo "No PID file found. Server may not be running."; \
	fi

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf bin
	rm -f *.out
	rm -f $(PID_FILE)
	rm -f .test_server.pid
	rm -f .test_server.log

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Build and run the application (default)"
	@echo "  build      - Build the application"
	@echo "  run        - Run the built application"
	@echo "  stop       - Stop the running server"
	@echo "  build-run  - Build and run in one command"
	@echo "  test       - Run unit tests"
	@echo "  test-e2e   - Run end-to-end tests"
	@echo "  test-all   - Run all tests (unit + e2e)"
	@echo "  clean      - Remove compiled binaries and test files"
	@echo "  help       - Show this help message"
