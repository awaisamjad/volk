// Package http implements a simple HTTP server and related utilities.
package http

import (
	"errors"
	"fmt"
	"strings"
)

// Request errors
var (
	ErrInvalidPathChars     = errors.New("path contains invalid characters")
	ErrDirectoryTraversal   = errors.New("path attempts directory traversal")
	ErrEmptyPath            = errors.New("path cannot be empty")
	ErrForbiddenPathSegment = errors.New("path contains forbidden segment")
)

// RequestStartLine represents the first line of an HTTP request
type RequestStartLine struct {
	Method        Method
	RequestTarget RequestTarget
	Protocol      Protocol
}

func (r RequestStartLine) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.RequestTarget.String(), r.Protocol)
}

// Request represents an HTTP request
type Request struct {
	StartLine RequestStartLine
	Headers   []Header
	Body      string
}

func (r Request) String() string {
	var sb strings.Builder
	sb.WriteString(r.StartLine.String())
	sb.WriteString(CRLF)

	for _, header := range r.Headers {
		sb.WriteString(header.String())
		sb.WriteString(CRLF)
	}

	sb.WriteString(CRLF)
	sb.WriteString(r.Body)

	return sb.String()
}

// GetHeaders returns the request headers
func (r Request) GetHeaders() []Header {
	return r.Headers
}

// GetBody returns the request body
func (r Request) GetBody() string {
	return r.Body
}

// GetMethod returns the request method
func (r Request) GetMethod() Method {
	return r.StartLine.Method
}

// GetRequestTarget returns the request target
func (r Request) GetRequestTarget() RequestTarget {
	return r.StartLine.RequestTarget
}

// GetProtocol returns the request protocol
func (r Request) GetProtocol() Protocol {
	return r.StartLine.Protocol
}

// NewRequest creates a new Request from a request string
func NewRequest(request_string string) (Request, error) {
	request, err := parseRequest(request_string)
	if err != nil {
		return Request{}, err
	}
	return request, nil
}

// ValidatePath validates the path in the request
func (r Request) ValidatePath() error {
	requestTarget := r.GetRequestTarget().String()
	if requestTarget == "" {
		return ErrEmptyPath
	}

	if requestTarget == "*" && r.StartLine.Method != OPTIONS {
		return fmt.Errorf("%s cannot use * as path", r.StartLine.Method)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return fmt.Errorf("path must start with /: %s", requestTarget)
	}

	if strings.Contains(requestTarget, "..") {
		return ErrDirectoryTraversal
	}

	segments := strings.SplitSeq(requestTarget, "/")
	for segment := range segments {
		if segment == "." || segment == ".." {
			return ErrForbiddenPathSegment
		}
	}

	// Check for invalid characters
	for _, c := range requestTarget {
		if c < 32 || c > 126 {
			return ErrInvalidPathChars
		}
	}

	return nil
}

// parseRequest parses a request string into a Request struct
func parseRequest(request string) (Request, error) {
	request = strings.Trim(request, " ")
	request_split := strings.Split(request, HeaderBodySeparator)
	if len(request_split) != 2 {
		return Request{}, fmt.Errorf("invalid request format: missing separator")
	}

	startline_headers := request_split[0]
	body := request_split[1]

	startline_headers_split := strings.Split(startline_headers, CRLF)
	if len(startline_headers_split) < 1 {
		return Request{}, fmt.Errorf("invalid request format: no startline")
	}

	startline := startline_headers_split[0]
	headers_strings := startline_headers_split[1:]

	startline_split := strings.Split(startline, " ")
	if len(startline_split) != 3 {
		return Request{}, fmt.Errorf("invalid startline format")
	}

	method := Method(startline_split[0])
	request_target_str := startline_split[1]
	protocol := Protocol(startline_split[2])

	path, err := parseRequestTarget(request_target_str)
	if err != nil {
		return Request{}, fmt.Errorf("invalid request target: %v", err)
	}

	request_target := RequestTarget{
		Path:     path,
		Query:    "",
		Fragment: "",
	}

	query, _, err := FindAndParseQuery(request_target_str)
	if err == nil {
		request_target.Query = "?" + strings.Join(func() []string {
			queryParts := []string{}
			for k, vs := range query.Params {
				for _, v := range vs {
					if v == "" {
						queryParts = append(queryParts, k)
					} else {
						queryParts = append(queryParts, k+"="+v)
					}
				}
			}
			return queryParts
		}(), "&")
	}

	fragment, _, err := FindAndParseFragment(request_target_str)
	if err == nil {
		request_target.Fragment = string(fragment)
	}

	headers := []Header{}
	for _, header_str := range headers_strings {
		if header_str == "" {
			continue
		}

		header, err := parseHeader(header_str)
		if err != nil {
			return Request{}, fmt.Errorf("invalid header: %v", err)
		}

		headers = append(headers, header)
	}

	return Request{
		StartLine: RequestStartLine{
			Method:        method,
			RequestTarget: request_target,
			Protocol:      protocol,
		},
		Headers: headers,
		Body:    body,
	}, nil
}
