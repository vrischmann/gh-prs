# Makefile for gh-prs GitHub CLI extension

# Version can be set via environment variable or git tag
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY_NAME = gh-prs
BUILD_DIR = dist

# Go build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -s -w"

# Platforms to build for
PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all build clean test install local-install help

# Default target
all: clean build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) version $(VERSION) for current platform..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for all platforms
build-all: clean
	@echo "Building $(BINARY_NAME) version $(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		GOOS=$(word 1,$(subst /, ,$(platform))) \
		GOARCH=$(word 2,$(subst /, ,$(platform))) \
		go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(subst /,-,$(platform))$(if $(findstring windows,$(platform)),.exe) . && \
	) echo "Built binaries for all platforms in $(BUILD_DIR)/"

# Install extension locally for testing
local-install: build
	@echo "Installing $(BINARY_NAME) locally..."
	@gh extension install .

# Remove locally installed extension
local-uninstall:
	@echo "Removing local $(BINARY_NAME) extension..."
	@gh extension remove prs 2>/dev/null || true

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(BINARY_NAME)

# Show version
version:
	@echo $(VERSION)

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-all     - Build for all supported platforms"
	@echo "  local-install - Install extension locally for testing"
	@echo "  local-uninstall - Remove locally installed extension"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  version       - Show version"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION       - Override version (default: git describe or 'dev')"