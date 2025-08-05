// Package http implements a simple HTTP server and related utilities.
package http

import (
	"fmt"
	"strconv"
	"strings"
)

// ResponseStartLine represents the first line of an HTTP response
type ResponseStartLine struct {
	Protocol   Protocol
	StatusCode StatusCode
	StatusText StatusText
}

func (r ResponseStartLine) String() string {
	return fmt.Sprintf("%s %d %s", r.Protocol, r.StatusCode, r.StatusText)
}

// Response represents an HTTP response
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

// GetProtocol returns the response protocol
func (r Response) GetProtocol() Protocol {
	return r.StartLine.Protocol
}

// GetStatusCode returns the response status code
func (r Response) GetStatusCode() StatusCode {
	return r.StartLine.StatusCode
}

// GetStatusText returns the response status text
func (r Response) GetStatusText() StatusText {
	return r.StartLine.StatusText
}

// GetHeaders returns the response headers
func (r Response) GetHeaders() []Header {
	return r.Headers
}

// GetBody returns the response body
func (r Response) GetBody() string {
	return r.Body
}

// NewResponse creates a new Response from a response string
func NewResponse(response_string string) (Response, error) {
	response, err := parseResponse(response_string)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

// parseResponse parses a response string into a Response struct
func parseResponse(response string) (Response, error) {
	response = strings.Trim(response, " ")
	response_split := strings.Split(response, HeaderBodySeparator)
	if len(response_split) != 2 {
		return Response{}, fmt.Errorf("invalid response format: missing separator")
	}

	startline_headers := response_split[0]
	body := response_split[1]

	startline_headers_split := strings.Split(startline_headers, CRLF)
	if len(startline_headers_split) < 1 {
		return Response{}, fmt.Errorf("invalid response format: no startline")
	}

	startline := startline_headers_split[0]
	headers_strings := startline_headers_split[1:]

	startline_split := strings.Split(startline, " ")
	if len(startline_split) < 3 {
		return Response{}, fmt.Errorf("invalid startline format")
	}

	protocol := Protocol(startline_split[0])
	status_code, err := strconv.Atoi(startline_split[1])
	if err != nil {
		return Response{}, fmt.Errorf("invalid status code: %v", err)
	}

	status_text := strings.Join(startline_split[2:], " ")

	headers := []Header{}
	for _, header_str := range headers_strings {
		if header_str == "" {
			continue
		}

		header, err := parseHeader(header_str)
		if err != nil {
			return Response{}, fmt.Errorf("invalid header: %v", err)
		}

		headers = append(headers, header)
	}

	return Response{
		StartLine: ResponseStartLine{
			Protocol:   protocol,
			StatusCode: StatusCode(status_code),
			StatusText: StatusText(status_text),
		},
		Headers: headers,
		Body:    body,
	}, nil
}