# Sultengutt Makefile
# Cross-platform build and test automation

# Variables
BINARY_NAME=sultengutt
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS=-ldflags="-w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build directories
BUILD_DIR=build
DIST_DIR=dist

# Platform-specific settings
ifeq ($(OS),Windows_NT)
    BINARY_EXT=.exe
    RM=del /Q
    MKDIR=mkdir
    RMDIR=rmdir /S /Q
else
    BINARY_EXT=
    RM=rm -f
    MKDIR=mkdir -p
    RMDIR=rm -rf
endif

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Sultengutt Build System"
	@echo "======================"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Testing targets
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./internal/...

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) -v -race ./...

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem ./...

## Building targets
.PHONY: build
build: clean ## Build for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	$(MKDIR) $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT) cmd/main.go

.PHONY: build-all
build-all: clean build-windows build-darwin ## Build for all platforms

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building for Windows..."
	$(MKDIR) $(DIST_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/main.go

.PHONY: build-darwin
build-darwin: ## Build for macOS (Intel and ARM)
	@echo "Building for macOS..."
	$(MKDIR) $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/main.go

.PHONY: install
install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(LDFLAGS) ./cmd

## Development targets
.PHONY: run
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run cmd/main.go

.PHONY: run-install
run-install: ## Run the install command
	@echo "Running $(BINARY_NAME) install..."
	$(GOCMD) run cmd/main.go install

.PHONY: run-status
run-status: ## Run the status command
	@echo "Running $(BINARY_NAME) status..."
	$(GOCMD) run cmd/main.go status

.PHONY: dev
dev: fmt vet test build ## Run development workflow (format, vet, test, build)

## Code quality targets
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: tidy
tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	$(GOMOD) tidy

.PHONY: download
download: ## Download Go modules
	@echo "Downloading Go modules..."
	$(GOMOD) download

## Release targets
.PHONY: release
release: ## Create a new release (prompts for version)
	@echo "Current version: $(VERSION)"
	@echo "Enter new version (e.g., v1.0.0):"
	@read NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then \
		echo "No version provided, aborting."; \
		exit 1; \
	fi; \
	echo "Creating release $$NEW_VERSION..."; \
	git tag $$NEW_VERSION; \
	git push origin $$NEW_VERSION; \
	echo "Release $$NEW_VERSION created and pushed!"

.PHONY: tag
tag: ## Create and push a git tag
	@if [ -z "$(TAG)" ]; then \
		echo "Usage: make tag TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating tag $(TAG)..."
	git tag $(TAG)
	git push origin $(TAG)

## Archive targets
.PHONY: archive
archive: build-all ## Create distribution archives
	@echo "Creating archives..."
	$(MKDIR) $(DIST_DIR)/archives
	@cd $(DIST_DIR) && \
	tar czf archives/$(BINARY_NAME)-windows-amd64.tar.gz $(BINARY_NAME)-windows-amd64.exe && \
	tar czf archives/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar czf archives/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@echo "Archives created in $(DIST_DIR)/archives/"

## Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@if [ -d $(BUILD_DIR) ]; then $(RMDIR) $(BUILD_DIR); fi
	@if [ -d $(DIST_DIR) ]; then $(RMDIR) $(DIST_DIR); fi
	@if [ -f coverage.out ]; then $(RM) coverage.out; fi
	@if [ -f coverage.html ]; then $(RM) coverage.html; fi

.PHONY: clean-test
clean-test: ## Clean test cache
	@echo "Cleaning test cache..."
	$(GOCMD) clean -testcache

.PHONY: clean-mod
clean-mod: ## Clean module cache
	@echo "Cleaning module cache..."
	$(GOCMD) clean -modcache

## Docker targets (if needed in the future)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@if [ -f Dockerfile ]; then \
		docker build -t $(BINARY_NAME):$(VERSION) .; \
	else \
		echo "No Dockerfile found"; \
	fi

## CI targets
.PHONY: ci
ci: fmt vet test-race build ## Run CI pipeline
	@echo "CI pipeline completed successfully!"

.PHONY: ci-coverage
ci-coverage: test-coverage ## Run CI with coverage
	@echo "CI with coverage completed successfully!"

## Information targets
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell $(GOCMD) version)"

.PHONY: deps
deps: ## Show dependency information
	@echo "Go modules:"
	$(GOMOD) list -m all

.PHONY: info
info: version deps ## Show build information
	@echo "Binary Name: $(BINARY_NAME)"
	@echo "Build Directory: $(BUILD_DIR)"
	@echo "Distribution Directory: $(DIST_DIR)"

## Quick targets
.PHONY: quick
quick: fmt test build ## Quick development cycle (format, test, build)

.PHONY: check
check: fmt vet lint test ## Run all checks

.PHONY: all
all: clean test-coverage build-all archive ## Build everything