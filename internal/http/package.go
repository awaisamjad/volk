// Package http implements a simple HTTP server and related utilities.
//
// This package provides functionality for parsing HTTP requests, generating HTTP responses,
// and serving static files. It supports all standard HTTP methods (GET, POST, PUT, etc.)
// and includes utilities for handling headers, query parameters, request paths, and more.
//
// The package is organized into multiple files:
// - constants.go: HTTP constants, methods, protocols, status codes
// - header.go: Header type and operations
// - query.go: Query parameter handling
// - fragment.go: Fragment handling
// - request_target.go: RequestTarget type and operations
// - request.go: Request type, parsing, and validation
// - response.go: Response type and creation
// - methods.go: HTTP method implementations (GET, POST, etc.)
// - server.go: FileServer and server-related code
// - util.go: Utility functions for HTTP operations
// - package.go: Package documentation and initialization
package http

// HTTPMessage is an interface that represents an HTTP message
// Both Request and Response implement this interface
type HTTPMessage interface {
	// GetHeaders returns the headers of the HTTP message
	GetHeaders() []Header

	// GetBody returns the body of the HTTP message
	GetBody() string

	// String returns a string representation of the HTTP message
	String() string
}

// DefaultFileServer is the default file server used for serving static files
var DefaultFileServer *FileServer

// init initializes the HTTP package
func init() {
	// Initialize with nil, will be set when the server starts
	DefaultFileServer = nil
}

// SetDefaultFileServer sets the default file server
func SetDefaultFileServer(fs *FileServer) {
	DefaultFileServer = fs
}
