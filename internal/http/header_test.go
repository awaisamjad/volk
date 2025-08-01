package http

import (
	"testing"
)

func TestHeaderString(t *testing.T) {
	testCases := []struct {
		name     string
		header   Header
		expected string
	}{
		{
			name:     "simple header",
			header:   Header{Name: "Content-Type", Value: "text/html"},
			expected: "Content-Type: text/html",
		},
		{
			name:     "empty value",
			header:   Header{Name: "Connection", Value: ""},
			expected: "Connection: ",
		},
		{
			name:     "empty name",
			header:   Header{Name: "", Value: "keep-alive"},
			expected: ": keep-alive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.header.String()
			if result != tc.expected {
				t.Errorf("Header.String() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected Header
		valid    bool
	}{
		{"Valid header", "Host: example.com", Header{Name: "Host", Value: "example.com"}, true},
		{"Header with spaces", "Content-Type: application/json", Header{Name: "Content-Type", Value: "application/json"}, true},
		{"Header with special chars", "X-Custom-Header: value!#$%&'*+.^_`|~", Header{Name: "X-Custom-Header", Value: "value!#$%&'*+.^_`|~"}, true},
		{"Header with empty value", "Empty: ", Header{Name: "Empty", Value: ""}, true},
		{"Header without colon", "InvalidHeader", Header{}, false},
		{"Header with invalid name chars", "Invalid Header: value", Header{}, false},
		{"Empty header", "", Header{}, false},
		{"Header with multiple colons", "Set-Cookie: name=value; expires=date", Header{Name: "Set-Cookie", Value: "name=value; expires=date"}, true},
		{"Header with whitespace", "   Server:   Apache   ", Header{Name: "Server", Value: "Apache"}, true},
		{"Header with empty name", ": value", Header{}, false},
		{"Header with numeric values", "Max-Age: 3600", Header{Name: "Max-Age", Value: "3600"}, true},
		{"Header with hyphens", "X-Forwarded-For: 192.168.1.1", Header{Name: "X-Forwarded-For", Value: "192.168.1.1"}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := parseHeader(test.header)

			if test.valid {
				if err != nil {
					t.Errorf("Expected valid header, got error: %v", err)
				}
				if result.Name != test.expected.Name || result.Value != test.expected.Value {
					t.Errorf("Expected %v, got %v", test.expected, result)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid header '%s', got none", test.header)
				}
			}
		})
	}
}