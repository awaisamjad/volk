// Package http implements a simple HTTP server and related utilities.
package http

import (
	"fmt"
	"regexp"
	"strings"
)

// Header represents an HTTP header
type Header struct {
	Name  string
	Value string
}

// String returns a string representation of the header
func (h Header) String() string {
	return h.Name + HeaderSeparator + h.Value
}

// parseHeader parses a header string into a Header struct
func parseHeader(header string) (Header, error) {
	var tokenPattern = regexp.MustCompile(`^[!#$%&'*+\.^_` + "`" + `|~0-9a-zA-Z-]+$`)
	firstColonIdx := strings.Index(header, ":")
	if firstColonIdx == -1 {
		return Header{}, fmt.Errorf("invalid header format: missing colon")
	}

	name := strings.TrimSpace(header[:firstColonIdx])
	value := strings.TrimSpace(header[firstColonIdx+1:])

	if name == "" {
		return Header{}, fmt.Errorf("invalid header format: empty name")
	}

	// RFC 7230 section 3.2.6 states that field names are tokens
	if !tokenPattern.MatchString(name) {
		return Header{}, fmt.Errorf("invalid header format: invalid characters in name")
	}

	return Header{
		Name:  name,
		Value: value,
	}, nil
}
