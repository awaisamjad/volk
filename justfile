APP_NAME := "volk"

VERSION := "v0.1.0"

BUILD_DIR := "build"

# -s: Omit symbol table and debug info.
# -w: Omit DWARF symbol table.
GO_BUILD_FLAGS := "-ldflags '-s -w'"


# Create the build directory if it doesn't exist.
# @mkdir {{BUILD_DIR}}

build-linux:
    @echo "Building for Linux (AMD64)..."
    GOOS=linux GOARCH=amd64 go build {{GO_BUILD_FLAGS}} -o {{BUILD_DIR}}/{{APP_NAME}}_linux_amd64 ./volk/main.go

build-windows:
    @echo "Building for Windows (AMD64)..."
    GOOS=windows GOARCH=amd64 go build {{GO_BUILD_FLAGS}} -o {{BUILD_DIR}}/{{APP_NAME}}_windows_amd64.exe ./volk/main.go
build-macos-amd64:
    @echo "Building for macOS (AMD64)..."
    GOOS=darwin GOARCH=amd64 go build {{GO_BUILD_FLAGS}} -o {{BUILD_DIR}}/{{APP_NAME}}_darwin_amd64 ./volk/main.go

build-macos-arm64:
    @echo "Building for macOS (ARM64)..."
    GOOS=darwin GOARCH=arm64 go build {{GO_BUILD_FLAGS}} -o {{BUILD_DIR}}/{{APP_NAME}}_darwin_arm64 ./volk/main.go

build-all: build-linux build-windows build-macos-amd64 build-macos-arm64
    @echo "All binaries built successfully in '{{BUILD_DIR}}/'"

# Create a GitHub release with the built binaries.
# This recipe requires the 'gh' CLI to be installed and authenticated.
release: build-all
    @echo "Attempting to create GitHub release {{VERSION}}..."
    # Check if gh CLI is installed
    @command -v gh >/dev/null 2>&1 || { echo >&2 "Error: 'gh' CLI is not installed. Please install it from https://cli.github.com/"; exit 1; }
    # Check if gh CLI is authenticated
    @gh auth status >/dev/null 2>&1 || { echo >&2 "Error: 'gh' CLI is not authenticated. Please run 'gh auth login' to authenticate."; exit 1; }

    # Create the release.
    # {{VERSION}}: This will be the Git tag name for the release (e.g., v1.0.0).
    # --notes: Provides the release description. You can point this to a file
    #          like `--notes-file RELEASE_NOTES.md` for more detailed notes.
    # --title: Sets the visible title of the release on GitHub.
    # {{BUILD_DIR}}/*: Attaches all files within the 'build' directory as release assets.
    gh release create {{VERSION}} \
        --notes "Automated release for {{APP_NAME}} {{VERSION}}" \
        --title "Release {{VERSION}}" \
        {{BUILD_DIR}}/*

    @echo "GitHub release '{{VERSION}}' created successfully!"
    @echo "Check your repository's releases page on GitHub."

go-wrk:
    go-wrk -c 2048 -d 10 http://localhost:6543/
# --- Utility Recipes ---

# Clean up all generated build artifacts
clean:
    @echo "Cleaning up build directory '{{BUILD_DIR}}/'..."
    @rm -rf {{BUILD_DIR}}
    @echo "Clean complete."

# Default recipe: If you just type 'just' without arguments, it will run 'build-all'.
default: build-all
