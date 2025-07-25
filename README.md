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
3. Run:
    ```bash
    ./volk
    ```
## Usage

Volk requires a `server.toml` for configuration and will look for a `server.toml` file in the following directories
- ./ (current directory)
- /config
- /etc/volk

To create a default configuration file in the current directory run either 

```bash 
volk serve -C
```
or
```bash 
volk serve --createConfig
```


### Default Configuration File


```toml
# Default Config File
[server]
port = 8000             # Port the server listens on
host = "0.0.0.0"        # Host address to bind to
read_timeout = 30       # Read timeout in seconds
write_timeout = 30      # Write timeout in seconds
max_connections = 100   # Maximum number of concurrent connections

[file_server]
document_root = "."                  # Root directory for serving files
default_file = "index.html"          # Default file to serve if a directory is requested
allow_directory_listing = false      # Whether to allow directory listing
[file_server.mime_type_overrides]
".dat" = "application/octet-stream"  # Override MIME type for .dat files
".custom" = "text/plain"             # Override MIME type for .custom files

[security]
allow_directory_traversal = false # Whether to allow directory traversal (should be false in production)
max_request_size = 1048576        # Maximum request size in bytes (1MB)
rate_limit = 60                   # Number of allowed requests per minute
allowed_origins = ["*"]           # Allowed CORS origins

[logging]
format = "plain"   # Logging format (plain, verbose)
file_path = ""     # Path to the log file (empty for stdout)
access_logs = true # Enable/disable access logs
```
## Development

### Project Structure (changes a lot)

```
.
├── build
│   ├── server.toml
│   ├── volk_darwin_amd64
│   ├── volk_darwin_arm64
│   ├── volk_linux_amd64
│   └── volk_windows_amd64.exe
├── config
│   ├── config.go
│   ├── config-test.toml
│   └── default-config.toml
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
├── test
│   └── main.go
└── volk
    ├── cmd
    │   ├── dumpconfig.go
    │   ├── root.go
    │   ├── serve.go
    │   └── version.go
    ├── index.html
    ├── main.go
    └── volk
```

### Testing

Run tests with:

```bash
go test ./...
```

## Release Process

I use `just` as my command runner. Check my commands in the `justfile`

## License

[MIT License](LICENSE)