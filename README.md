# Volk

Volk is a lightweight HTTP server written in Go, designed to serve static files with minimal configuration.

## Features

- Simple and fast file server
- Cross-platform support (Windows, macOS, Linux)
- Minimal dependencies
- Configuration via TOML files

## Current Limitations

- Only implements `GET` functionality (Will implement the rest in the future insha'Allah)
- Not the most performant

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/awais/volk/releases).

### Building from Source

#### Prerequisites

- Go 1.19 or newer
- [Just](https://just.systems/man/en/) command runner (optional but recommended)

#### Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/awais/volk.git
    cd /volk/cmd/server
    ```

2. Build:

    ```bash
    go build
    ```

## Usage

### Basic Usage

By default, Volk will serve files from the current directory and looks for `index.html`.

```bash
./volk
```

### Using Configuration File

Create a `server.toml` file:

```toml
port = 8080
root = "public"
index = "index.html"
```

Volk automatically looks for a `server.toml` file so no flags need to given:

```bash
./volk 
```

### Configuration Options

| Option  | Description                   | Default      |
| ------- | ----------------------------- | ------------ |
| `port`  | Port to listen on             | `8080`       |
| `root`  | Directory to serve files from | `.`          |
| `index` | Index file to serve           | `index.html` |

## Development

### Project Structure

```

.
├── build
│   ├── volk_darwin_amd64
│   ├── volk_darwin_arm64
│   ├── volk_linux_amd64
│   └── volk_windows_amd64.exe
├── cmd
│   ├── client
│   │   └── main.go
│   └── server
│       ├── about
│       │   └── index.html
│       ├── config-test.toml
│       ├── config.toml
│       ├── index.html
│       ├── server.toml
│       └── volk.go
├── config
│   └── config.go
├── go.mod
├── go.sum
├── internal
│   └── http
│       ├── fileserver.go
│       ├── http.go
│       └── http_test.go
├── justfile
├── main.go
├── README.md
└── test
    └── main.go
```

### Testing

Run tests with:

```bash
go test ./...
```

## Release Process

1. Update version in `justfile`
2. Build all binaries: `just build-all`
3. Create a release: `just release VERSION=vX.Y.Z`

## License

[MIT License](LICENSE)