// Package http implements a simple HTTP server and related utilities.
package http

// CRLF is the standard HTTP line ending
const CRLF = "\r\n"

// HeaderBodySeparator is the separator between HTTP headers and body
const HeaderBodySeparator = CRLF + CRLF

// HeaderSeparator is the separator between HTTP header name and value
const HeaderSeparator = ": "

// Method represents an HTTP method
type Method string

// HTTP methods as defined in RFC 7231
const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	HEAD    Method = "HEAD"
	OPTIONS Method = "OPTIONS"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
)

// Protocol represents an HTTP protocol version
type Protocol string

// HTTP protocol versions
const (
	HTTP1_1 Protocol = "HTTP/1.1"
	HTTP1_0 Protocol = "HTTP/1.0"
	HTTP0_9 Protocol = "HTTP/0.9"
)

// StatusCode represents an HTTP status code
type StatusCode int

// StatusText represents an HTTP status text
type StatusText string

// StatusCodeMap maps status codes to their text representations
var StatusCodeMap = map[StatusCode]StatusText{
	200: "OK",
	201: "Created",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	304: "Not Modified",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
}
