.PHONY: help build install clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the tracehub binary
	@echo "Building tracehub..."
	@go build -o tracehub ./cmd/tracehub

install: build ## Install to /usr/local/bin
	@echo "Installing to /usr/local/bin..."
	@sudo mv tracehub /usr/local/bin/
	@echo "✓ Installed! Run 'tracehub --help' to get started"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f tracehub
	@rm -rf bin/ dist/

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

.DEFAULT_GOAL := help
