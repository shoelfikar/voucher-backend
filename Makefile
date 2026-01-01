.PHONY: run build test clean migrate-up migrate-down install lint format help

# Variables
BINARY_NAME=voucher-api
MAIN_PATH=cmd/api/main.go
BIN_DIR=bin

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: Run the application
run:
	@echo "Running application..."
	go run $(MAIN_PATH)

## build: Build the application binary
build:
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary created at $(BIN_DIR)/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Clean complete"

## install: Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## format: Format code
format:
	@echo "Formatting code..."
	gofmt -w .
	@echo "Code formatted"

## dev: Run in development mode with auto-reload (requires air)
dev:
	@echo "Running in development mode..."
	air
