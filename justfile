# justfile for gh-prs GitHub CLI extension

# Configuration
binary_name := "gh-prs"
build_dir := "dist"
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`

# Go build flags
ldflags := "-ldflags \"-X main.version=" + version + " -s -w\""

# Default recipe
default: clean build

# Build for current platform
build:
    @echo "Building {{binary_name}} version {{version}} for current platform..."
    go build {{ldflags}} -o {{binary_name}} .

# Build for specific platform
build-platform goos goarch:
    #!/usr/bin/env bash
    output={{build_dir}}/{{binary_name}}-{{goos}}-{{goarch}}
    if [ "{{goos}}" = "windows" ]; then
        output="${output}.exe"
    fi
    echo "Building for {{goos}}/{{goarch}}..."
    GOOS={{goos}} GOARCH={{goarch}} go build {{ldflags}} -o "$output" .

# Build for all platforms
build-all: clean
    @echo "Building {{binary_name}} version {{version}} for all platforms..."
    mkdir -p {{build_dir}}
    just build-platform linux amd64
    just build-platform linux arm64
    just build-platform darwin amd64
    just build-platform darwin arm64
    just build-platform windows amd64
    just build-platform windows arm64
    @echo "Built binaries for all platforms in {{build_dir}}/"

# Install extension locally for testing
local-install: build
    @echo "Installing {{binary_name}} locally..."
    gh extension install .

# Remove locally installed extension
local-uninstall:
    @echo "Removing local {{binary_name}} extension..."
    -gh extension remove prs 2>/dev/null

# Run tests
test:
    @echo "Running tests..."
    go test ./...

# Run linter
lint:
    @echo "Running linter..."
    golangci-lint run

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf {{build_dir}} {{binary_name}}

# Show version
version:
    @echo {{version}}

# List available recipes
list:
    @just --list