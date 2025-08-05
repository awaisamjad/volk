# Volk

Volk is a lightweight HTTP server/File Server written in Go, designed to serve static files with minimal configuration.

## Features

- Simple and fast file server
- Customizable configuration via TOML files
- Cross-platform support (Windows, macOS, Linux)
- Minimal external dependencies
- Command-line interface with simple commands

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/awaisamjad/volk/releases).

### Building from Source

#### Prerequisites

- Go 1.24 or newer

#### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/awaisamjad/volk.git
   cd volk
   ```

2. Build:

   ```bash
   go build -o volk ./volk/main.go
   ```

3. Or use the provided Just commands:

   ```bash
   just build-all
   ```

## Usage

### Quick Start

1. Create an `index.html` file in your project directory
2. Run `./volk serve` in the same directory

By default, Volk will:

- Look for an `index.html` file in the current directory
- Serve the file on `http://0.0.0.0:6543`

### Configuration

Volk can be configured using a `volk_config.toml` file. You can generate a default configuration file with:

```bash
./volk dump-config
```

### Default Configuration

```toml
[server]
port = 8000           # Port the server listens on
read_timeout = 30     # Read timeout in seconds

[file_server]
document_root = "."             # Root directory for serving files
default_file = "index.html"     # Default file to serve if a directory is requested

[logging]
format = "plain"   # Logging format (plain, verbose)
file_path = ""     # Path to the log file (empty for stdout)
access_logs = true # Enable/disable access logs
```

## Project Structure

```
.
├── build                 # Build outputs
├── config                # Configuration related code
│   ├── config.go
│   └── default-config.toml
├── internal              # Internal packages
│   └── http              # HTTP implementation
│       ├── fileserver.go
│       ├── request.go
│       ├── response.go
│       └── ...
├── volk                  # Main application code
│   ├── cmd               # CLI commands
│   │   ├── dump_default_config.go
│   │   ├── root.go
│   │   └── serve.go
│   └── main.go           # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── justfile              # Just commands for building and releasing
└── README.md             # This file
```

## Development

### Building for Different Platforms

The project includes a `justfile` with recipes for building binaries for various platforms:

```bash
# Build for all supported platforms
just build-all

# Build for specific platforms
just build-linux
just build-windows
just build-macos-amd64
just build-macos-arm64

# Clean build artifacts
just clean
```

### Testing

Run tests with:

```bash
go test ./...
```

## Release Process

The project uses `just` as a command runner for managing builds and releases.

To create a new release:

```bash
just release
```

This will build binaries for all platforms and create a GitHub release using the `gh` CLI.

## License

[MIT License](LICENSE)
