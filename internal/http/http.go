// Package http provides functionality for HTTP protocol handling.
//
// It includes structures and methods for parsing, validating, and processing HTTP
// requests and responses. The package supports the various HTTP methods : GET,
// POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, and CONNECT, along with HTTP
// protocol versions 0.9, 1.0, and 1.1.
//
// Key features include:
// - Request and response parsing from raw strings
// - HTTP message formatting to standard format
// - Path validation for security concerns
// - Header parsing and validation
// - Status code handling with appropriate text responses
// - Query parameter and URL fragment parsing
// - Method-specific request handlers
//
// The package forms the foundation for building both HTTP clients and servers,
// with particular focus on correctness and security.

package http

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"volk/config"
)

var (
	DefaultFileServer *FileServer
)

const (
	// CRLF is the HTTP line ending (Carriage Return + Line Feed)
	CRLF = "\r\n"

	// HeaderBodySeparator is the separator between headers and body (double CRLF)
	HeaderBodySeparator = CRLF + CRLF

	// HeaderSeparator is the separator between individual headers
	HeaderSeparator = CRLF
)

type Method string

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

type Protocol string

const (
	HTTP1_1 Protocol = "HTTP/1.1"
	HTTP1_0 Protocol = "HTTP/1.0"
	HTTP0_9 Protocol = "HTTP/0.9"
)

type StatusCode uint64
type StatusText string

var StatusCodeMap = map[StatusCode]StatusText{
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	409: "Conflict",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

type Header struct {
	Name  string
	Value string
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s", h.Name, h.Value)
}

type RequestStartLine struct {
	Method        Method
	RequestTarget string
	Protocol      Protocol
}

func (r RequestStartLine) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.RequestTarget, r.Protocol)
}

// HTTPMessage defines common behavior for HTTP messages (requests and responses)
type HTTPMessage interface {
	// GetHeaders returns all headers in the message
	GetHeaders() []Header

	// GetBody returns the message body as a string
	GetBody() string

	// String returns the complete formatted HTTP message
	String() string
}

type Request struct {
	StartLine RequestStartLine
	Headers   []Header
	Body      string
}

func (r Request) String() string {
	var builder strings.Builder

	builder.WriteString(r.StartLine.String())
	builder.WriteString("\r\n")

	for _, header := range r.Headers {
		builder.WriteString(header.String())
		builder.WriteString("\r\n")
	}

	builder.WriteString("\r\n")
	builder.WriteString(r.Body)

	return builder.String()
}

// GetHeaders returns all headers in the request
func (r Request) GetHeaders() []Header {
	return r.Headers
}

// GetBody returns the request body as a string
func (r Request) GetBody() string {
	return r.Body
}

// GetMethod returns the HTTP method of the request
func (r Request) GetMethod() Method {
	return r.StartLine.Method
}

// GetRequestTarget returns the request target (path)
func (r Request) GetRequestTarget() string {
	return r.StartLine.RequestTarget
}

// GetProtocol returns the HTTP protocol version
func (r Request) GetProtocol() Protocol {
	return r.StartLine.Protocol
}

var (
	ErrInvalidPathChars     = errors.New("path contains invalid characters")
	ErrDirectoryTraversal   = errors.New("path attempts directory traversal")
	ErrEmptyPath            = errors.New("path cannot be empty")
	ErrForbiddenPathSegment = errors.New("path contains forbidden segment")
)

func (r Request) ValidatePath() error {
	requestTarget := r.GetRequestTarget()
	if requestTarget == "" {
		return ErrEmptyPath
	}

	if requestTarget == "*" && r.StartLine.Method != OPTIONS {
		return fmt.Errorf("%s cannot use * as path", r.StartLine.Method)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return errors.New("path doesn't contain slash")
	}

	if strings.Contains(requestTarget, "../") {
		return ErrDirectoryTraversal
	}

	//Parse the URL to leverage net/url's safety features and isolate the path.
	// url.Parse itself can catch many invalid URI formats
	u, err := url.Parse("http://" + "dummy.com" + requestTarget)
	if err != nil {
		return fmt.Errorf("malformed URL: %w", err)
	}
	parsedPath := u.Path
	cleanedPath := path.Clean(parsedPath)
	if cleanedPath != "/" && !strings.HasPrefix(cleanedPath, "/") {
		return ErrDirectoryTraversal
	}

	return nil
}

type ResponseStartLine struct {
	Protocol   Protocol
	StatusCode StatusCode
	StatusText StatusText
}

func (r ResponseStartLine) String() string {
	return fmt.Sprintf("%s %d %s", r.Protocol, r.StatusCode, r.StatusText)
}

type Response struct {
	StartLine ResponseStartLine
	Headers   []Header
	Body      string
}

func (r Response) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %d %s\r\n",
		r.GetProtocol(), r.GetStatusCode(), r.GetStatusText()))

	for _, header := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", header.Name, header.Value))
	}
	builder.WriteString("\r\n")
	builder.WriteString(r.Body)

	return builder.String()
}

func (r Response) GetProtocol() Protocol {
	return r.StartLine.Protocol
}

func (r Response) GetStatusCode() StatusCode {
	return r.StartLine.StatusCode
}

func (r Response) GetStatusText() StatusText {
	return r.StartLine.StatusText
}
func (r Response) GetHeaders() []Header {
	return r.Headers
}

func (r Response) GetBody() string {
	return r.Body
}

func NewRequest(request_string string) (Request, error) {
	request, err := parseRequest(request_string)
	if err != nil {
		return Request{}, err
	}
	return request, nil
}

func NewResponse(response_string string) (Response, error) {
	response, err := parseResponse(response_string)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

// parseRequest parses a raw HTTP request string into a structured Request object.
// It extracts the method, request target, protocol, headers, and body from the string.
//
// This function is intended for internal use within the http package.
// Use NewRequest() to create a new Request struct from a raw request string.
//
// Returns an error if the request string is malformed or contains invalid data.
func parseRequest(request string) (Request, error) {
	// Seperate first by empty line to get startline and headers together and optional body by itself
	request = strings.Trim(request, " ")
	request_split := strings.Split(request, HeaderBodySeparator)
	if len(request_split) != 2 {
		return Request{}, fmt.Errorf("invalid request format: missing separator")
	}

	startline_headers := request_split[0]
	body := request_split[1]

	startline_headers_split := strings.Split(startline_headers, HeaderSeparator)
	startline := startline_headers_split[0]
	startline_split := strings.Split(startline, " ")

	if len(startline_split) != 3 {
		return Request{}, fmt.Errorf("invalid startline format: should have 3 parts")
	}
	method := startline_split[0]
	request_target := startline_split[1]
	protocol := startline_split[2]

	//? Only these methods should have a body
	if (Method(method) != POST && Method(method) != PUT && Method(method) != PATCH) && (len(body) != 0) {
		return Request{}, fmt.Errorf("%s can not have a body", method)
	}

	// Validate method
	switch Method(method) {
	case GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT:
		// Valid method, no action needed
	default:
		return Request{}, fmt.Errorf("HTTP Method %s is not supported", method)
	}

	// Validate Request Target
	err := Request{
		StartLine: RequestStartLine{
			Method:        Method(method),
			RequestTarget: request_target,
			Protocol:      Protocol(protocol),
		},
	}.ValidatePath()

	if err != nil {
		return Request{}, err
	}

	// Validate protocol
	switch Protocol(protocol) {
	case HTTP1_1, HTTP1_0, HTTP0_9:
		// Valid protocol, no action needed
	default:
		return Request{}, fmt.Errorf("invalid protocol: %s", protocol)
	}

	var headers []Header
	for _, header := range startline_headers_split[1:] {
		parsed_header, err := parseHeader(header)
		if err != nil {
			return Request{}, fmt.Errorf("invalid header: %s %s", header, err)
		}

		headers = append(headers, parsed_header)
	}

	return Request{
		RequestStartLine{
			Method:        Method(method),
			RequestTarget: request_target,
			Protocol:      Protocol(protocol),
		},
		headers,
		body,
	}, nil
}

// parseResponse parses a raw HTTP response string into a structured Response object.
// It extracts the protocol, status code, headers, and body from the string.
//
// This function is intended for internal use within the http package.
// Use NewResponse() to create a new Response struct from a raw response string.
//
// Returns an error if the response string is malformed or contains invalid data.
func parseResponse(response string) (Response, error) {
	// Seperate first by empty line to get startline and headers together and optional body by itself
	response = strings.Trim(response, " ")
	response_split := strings.Split(response, HeaderBodySeparator)
	if len(response_split) != 2 {
		return Response{}, fmt.Errorf("invalid response format: missing separator")
	}

	startline_headers := response_split[0]
	body := response_split[1]

	startline_headers_split := strings.Split(startline_headers, HeaderSeparator)
	startline := startline_headers_split[0]
	startline_split := strings.Split(startline, " ")

	//? strategy is to split on whitespace take first two as protocol and status code and join rest as status text
	if len(startline_split) < 2 {
		return Response{}, fmt.Errorf("invalid startline format: missing either protocol or status code")
	}

	var protocol string
	var statusCode uint64
	// if len(startline_split) == 2 {
	protocol = startline_split[0]
	switch Protocol(protocol) {
	case HTTP1_1, HTTP1_0, HTTP0_9:
		// Valid protocol, no action needed
	default:
		return Response{}, fmt.Errorf("invalid protocol: %s", protocol)
	}

	var err error
	statusCode, err = strconv.ParseUint(startline_split[1], 10, 64)
	if err != nil {
		return Response{}, fmt.Errorf("invalid statuscode: %d", statusCode)
	}
	if _, ok := StatusCodeMap[StatusCode(statusCode)]; !ok {
		return Response{}, fmt.Errorf("invalid statuscode: %d", statusCode)
	}
	// }

	statusText := strings.Join(startline_split[2:], " ")
	found := false
	for _, text := range StatusCodeMap {
		if text == StatusText(statusText) {
			found = true
			break
		}
	}
	if !found {
		return Response{}, fmt.Errorf("invalid status text: %s", statusText)
	}

	var headers []Header
	for _, header := range startline_headers_split[1:] {
		parsed_header, err := parseHeader(header)
		if err != nil {
			return Response{}, fmt.Errorf("invalid header: %s %s", header, err)
		}

		headers = append(headers, parsed_header)
	}

	return Response{
		ResponseStartLine{
			Protocol:   Protocol(protocol),
			StatusCode: StatusCode(statusCode),
			StatusText: StatusText(statusText),
		},
		headers,
		body,
	}, nil
}

func parseHeader(header string) (Header, error) {
	var tokenPattern = regexp.MustCompile(`^[!#$%&'*+\.^_` + "`" + `|~0-9a-zA-Z-]+$`)
	firstColonIdx := strings.Index(header, ":")
	if firstColonIdx == -1 {
		return Header{}, fmt.Errorf("invalid header format: missing colon")
	}

	name := strings.TrimSpace(header[:firstColonIdx])
	value := strings.TrimSpace(header[firstColonIdx+1:])

	if name == "" {
		return Header{}, fmt.Errorf("invalid header: name is empty")
	}

	if !tokenPattern.MatchString(name) {
		return Header{}, fmt.Errorf("invalid header name: '%s'", name)
	}

	// ? we will allow empty values
	// if value == "" {
	// 	return Header{}, fmt.Errorf("invalid header: value is empty for '%s'", name)
	// }

	return Header{
		Name:  name,
		Value: value,
	}, nil
}

// func parseHeader(header string) (Header, error) {
// 	first_colon_idx := strings.Index(header, ":")
// 	if first_colon_idx == -1 {
// 		return Header{}, fmt.Errorf("invalid header format: missing colon")
// 	}

// 	name := header[:first_colon_idx]
// 	value := header[first_colon_idx+1:]
// 	name = strings.TrimSpace(name)
// 	value = strings.TrimSpace(value)
// 	var tokenPattern = regexp.MustCompile(`^[!#$%&'*+\.^_` + "`" + `|~0-9a-zA-Z-]+$`)

// 	if !tokenPattern.MatchString(name) {
// 		return Header{}, fmt.Errorf("invalid header name: '%s'", name)
// 	}

// 	if len(strings.Split(name, " ")) != 1 {
// 		return Header{}, fmt.Errorf("invalid header format: most likely missing colon after name")
// 	}

// 	if name == "" {
// 		return Header{}, fmt.Errorf("invalid header format: header name cannot be empty in '%s'", header)
// 	}
// 	if value == "" {
// 		return Header{}, fmt.Errorf("invalid header format: header value cannot be empty in '%s'", header)
// 	}

// 	return Header{
// 		Name:  name,
// 		Value: value,
// 	}, nil
// }

// ?category=electronics&brand=sony&price_max=500
type Query struct {
	Params map[string][]string
}

func findQuery(requestTarget string) (string, error) {

	query := ""

	queryIdx := strings.Index(requestTarget, "?")

	if queryIdx == -1 {
		log.Println("")
		return "", errors.New("query not found")
	}

	isThereAFragment := false
	fragmentIdx := strings.Index(requestTarget, "#")
	if fragmentIdx != -1 {
		isThereAFragment = true
	}

	if isThereAFragment && fragmentIdx > queryIdx {
		for i := queryIdx; i < fragmentIdx; i++ {
			query += string(requestTarget[i])
		}
	} else {
		for i := queryIdx; i < len(requestTarget); i++ {
			query += string(requestTarget[i])
		}
	}

	query = strings.Trim(query, " ")
	if len(query) == 0 {
		return "", errors.New("query is empty")
	}

	if query[0] == '?' {
		query = query[1:]
	} else {
		return "", errors.New("query doesnt start with ?")
	}

	return query, nil
}

// ? parseQuery takes the string given by findQuery() and then parses it. it is not responsible for locating the query in the request-target
// func parseQuery(query Query) (Query, error) {

// 	queryIdx := strings.Index(requestTarget, "?")
// 	isQueryPresent := false
// 	if queryIdx != -1 {
// 		isQueryPresent = true
// 	}
// queries := strings.Split(query, "&")
// params := make(map[string][]string, 0)
// // Process each query parameter
// for _, q := range queries {
// 	if q == "" || !strings.Contains(q, "=") {
// 		continue // Skip empty or malformed parameters
// 	}

// 	parts := strings.SplitN(q, "=", 2)
// 	key := parts[0]
// 	value := ""

// 	if len(parts) > 1 {
// 		value = parts[1]
// 	}

// 	// URL decode the key and value
// 	decodedKey, err := url.QueryUnescape(key)
// 	if err != nil {
// 		return Query{}, fmt.Errorf("invalid query parameter key: %w", err)
// 	}

// 	decodedValue, err := url.QueryUnescape(value)
// 	if err != nil {
// 		return Query{}, fmt.Errorf("invalid query parameter value: %w", err)
// 	}

// 	// Add to params map (support for multiple values with same key)
// 	params[decodedKey] = append(params[decodedKey], decodedValue)
// }
// 	return Query{}, nil
// }

type Fragment string

// ? Given a request-target findFragment will find the fragment (if any) and return it
func findFragment(requestTarget string) (Fragment, error) {
	fragment := ""

	fragmentIdx := strings.Index(requestTarget, "#")

	if fragmentIdx == -1 {
		return "", errors.New("Fragment not found")
	}

	for i := fragmentIdx; i < len(requestTarget); i++ {
		fragment += string(requestTarget[i])
	}

	fragment = strings.Trim(string(fragment), " ")

	if strings.Contains(fragment, " ") {
		return "", errors.New("Invalid Fragment : cannot have spaces")
	}

	return Fragment(fragment), nil
}

// parseFragment processes a URI fragment identifier.
//
// It takes a Fragment type (string wrapper) and validates that:
//   - It's not an empty string
//   - It starts with '#' character
//   - It contains no whitespace
//
// The function strips the leading '#' character and returns the fragment content.
// For example, "#section1" becomes "section1".
//
// Returns:
//   - Fragment: The processed fragment without the '#' prefix
//   - error: If validation fails, an appropriate error is returned
func parseFragment(fragment Fragment) (Fragment, error) {
	if fragment == "" {
		return "", errors.New("fragment cannot be empty")
	}

	if !strings.HasPrefix(string(fragment), "#") {
		return "", fmt.Errorf("invalid fragment: %s must start with '#'", fragment)
	}

	if strings.Contains(string(fragment), "\n") {
		return "", errors.New("fragment cannot contain whitespace")
	}

	fragment = fragment[1:]
	fragment = Fragment(strings.ReplaceAll(string(fragment), " ", "%20"))
	fragment = Fragment(strings.ReplaceAll(string(fragment), "#", "%23"))

	return fragment, nil
}

func parseRequestTarget(requestTarget string) (string, error) {
	// fragment, err := findFragment(requestTarget)
	return requestTarget, nil
}

// Get the response from a Request object
func (rq *Request) Response() Response {
	switch rq.GetMethod() {
	case GET:
		return rq.GET()
	case POST:
		return rq.POST()
	case PUT:
		return rq.PUT()
	case DELETE:
		return rq.DELETE()
	case PATCH:
		return rq.PATCH()
	case HEAD:
		return rq.HEAD()
	case OPTIONS:
		return rq.OPTIONS()
	case TRACE:
		return rq.TRACE()
	case CONNECT:
		return rq.CONNECT()
	default:
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   rq.StartLine.Protocol,
				StatusCode: 501,
				StatusText: StatusCodeMap[501],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "501 Not Implemented: Only GET is currently implemented",
		}

	}
}

func (rq *Request) GET() Response {
	path := rq.GetRequestTarget()
	switch path {
	case "*":
		//? only OPTIONS method allowed to use *
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   rq.StartLine.Protocol,
				StatusCode: 400,
				StatusText: StatusCodeMap[400],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "400 Bad Request: '*' only allowed with OPTIONS method",
		}

	default:
		if DefaultFileServer != nil {
			return DefaultFileServer.ServeFile(rq)
		} else {
			// Fallback to default config if no file server is configured
			fileserver := NewFileServer(config.FileServerConfig{
				DocumentRoot: ".",
				DefaultFile:  "index.html",
			})
			return fileserver.ServeFile(rq)
		}
	}
}

func (rq *Request) POST() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: POST is not implemented",
	}
}

func (rq *Request) PUT() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: PUT is not implemented",
	}
}

func (rq *Request) DELETE() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: DELETE is not implemented",
	}
}

func (rq *Request) PATCH() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: PATCH is not implemented",
	}
}

func (rq *Request) HEAD() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: HEAD is not implemented",
	}
}

func (rq *Request) OPTIONS() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: OPTIONS is not implemented",
	}
}

func (rq *Request) TRACE() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: TRACE is not implemented",
	}
}

func (rq *Request) CONNECT() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: CONNECT is not implemented",
	}
}
