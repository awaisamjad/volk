package http

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	t.Run("Test Valid GET Request", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		_, err := ParseRequest(requestString)
		if err != nil {
			t.Errorf("ParseRequest returned an error: %v", err)
		}
	})

	t.Run("Test Invalid GET Request", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n<h1>Hello World</h1>"
		_, err := ParseRequest(requestString)
		if err == nil {
			t.Errorf("ParseRequest returned an error: %v", err)
		}
	})
	t.Run("Test Valid POST Request", func(t *testing.T) {
		requestString := "POST /submit HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/json\r\n\r\n{\"name\": \"John Doe\"}"
		_, err := ParseRequest(requestString)
		if err != nil {
			t.Errorf("ParseRequest returned an error: %v", err)
		}
	})

	t.Run("Test Request with no separator", func(t *testing.T) {
		requestString := "GET /index.html HTTP/1.1\r\nHost: localhost:8080\r\n"
		_, err := ParseRequest(requestString)
		if err == nil {
			t.Errorf("ParseRequest should have returned an error due to missing separator. Error : %v", err)
		}
	})

	t.Run("Test Invalid Method", func(t *testing.T) {
		requestString := "FOO / HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
		_, err := ParseRequest(requestString)
		if err == nil {
			t.Errorf("ParseRequest should have returned an error due to invalid method")
		}
	})

	t.Run("Test Invalid Protocol", func(t *testing.T) {
		requestString := "GET / HTTP/2.0\r\nHost: localhost:8080\r\n\r\n"
		_, err := ParseRequest(requestString)
		if err == nil {
			t.Errorf("ParseRequest should have returned an error due to invalid protocol")
		}
	})

	t.Run("Test Invalid Header", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost localhost:8080\r\n\r\n"
		_, err := ParseRequest(requestString)
		if err == nil {
			t.Errorf("ParseRequest should have returned an error due to invalid header")
		}
	})

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
func TestResponseParsing(t *testing.T) {
	t.Run("Test Valid HTTP Response", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 22\r\n\r\n<h1>Hello World</h1>"
		_, err := ParseResponse(responseString)
		if err != nil {
			t.Errorf("ParseResponse returned an error: %v", err)
		}
	})

	t.Run("Test Response without body", func(t *testing.T) {
		responseString := "HTTP/1.1 204 No Content\r\nServer: TestServer\r\nDate: Mon, 27 Jul 2022 12:28:53 GMT\r\n\r\n"
		_, err := ParseResponse(responseString)
		if err != nil {
			t.Errorf("ParseResponse returned an error: %v", err)
		}
	})

	t.Run("Test Response with missing separator", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 22"
		_, err := ParseResponse(responseString)
		if err == nil {
			t.Errorf("ParseResponse should have returned an error due to missing separator")
		}
	})

	t.Run("Test Invalid Protocol", func(t *testing.T) {
		responseString := "HTTP/2.0 200 OK\r\nContent-Type: text/html\r\n\r\n<h1>Hello World</h1>"
		_, err := ParseResponse(responseString)
		if err == nil {
			t.Errorf("ParseResponse should have returned an error due to invalid protocol")
		}
	})

	t.Run("Test Invalid Status Code", func(t *testing.T) {
		responseString := "HTTP/1.1 999 Invalid\r\nContent-Type: text/html\r\n\r\n<h1>Error</h1>"
		_, err := ParseResponse(responseString)
		if err == nil {
			t.Errorf("ParseResponse should have returned an error due to invalid status code")
		}
	})

	t.Run("Test Invalid Header Format", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type text/html\r\n\r\n<h1>Hello World</h1>"
		_, err := ParseResponse(responseString)
		if err == nil {
			t.Errorf("ParseResponse should have returned an error due to invalid header format")
		}
	})
}

func TestNewRequestAndResponse(t *testing.T) {
	t.Run("Test NewRequest with valid request string", func(t *testing.T) {
		requestString := "GET /files/index.html HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}
		if req.StartLine.Method != GET {
			t.Errorf("Expected method GET, got %s", req.StartLine.Method)
		}
		if req.StartLine.RequestTarget != "/files/index.html" {
			t.Errorf("Expected request target /files/index.html, got %s", req.StartLine.RequestTarget)
		}
		if req.StartLine.RequestTarget == "files/index.html" {
			t.Errorf("Expected request target %s got files/index.html", req.StartLine.RequestTarget)
		}
		if req.StartLine.Protocol != HTTP1_1 {
			t.Errorf("Expected protocol HTTP/1.1, got %s", req.StartLine.Protocol)
		}
	})

	t.Run("Test Request.Response for GET request", func(t *testing.T) {
		requestString := "GET /files/index.html HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}

		resp := req.Response()
		if resp.StartLine.Protocol != HTTP1_1 {
			t.Errorf("Expected protocol HTTP/1.1, got %s", resp.StartLine.Protocol)
		}

		// The actual status code might vary depending on if the file exists
		// so we're not testing the exact value here
	})
	t.Run("Test directory for request-target for GET request", func(t *testing.T) {
		requestString := "GET /files/ HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}

		resp := req.Response()
		if resp.StartLine.Protocol != HTTP1_1 {
			t.Errorf("Expected protocol HTTP/1.1, got %s", resp.StartLine.Protocol)
		}

		// The actual status code might vary depending on if the file exists
		// so we're not testing the exact value here
	})
	t.Run("Test Request.Response for non-GET method", func(t *testing.T) {
		requestString := "POST /submit HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/json\r\n\r\n{\"data\":\"test\"}"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}

		resp := req.Response()
		if resp.StartLine.StatusCode != 501 {
			t.Errorf("Expected status code 501, got %d", resp.StartLine.StatusCode)
		}
		if resp.StartLine.StatusText != "Not Implemented" {
			t.Errorf("Expected status text Not Implemented, got %s", resp.StartLine.StatusText)
		}
	})

	t.Run("Test NewRequest with invalid request string", func(t *testing.T) {
		requestString := "INVALID REQUEST"
		_, err := NewRequest(requestString)
		if err == nil {
			t.Errorf("Expected error for invalid request string, got nil")
		}
	})

	t.Run("Test GET request with * target", func(t *testing.T) {
		requestString := "GET * HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
		req, err := NewRequest(requestString)
		if err == nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}

		resp := req.Response()
		if resp.StartLine.StatusCode != 501 {
			t.Errorf("Expected status code 501 for * target, got %d", resp.StartLine.StatusCode)
		}
	})
}

func TestFileServerPathValidation(t *testing.T) {
	t.Run("Test Valid Path", func(t *testing.T) {
		requestString := "GET /files/index.html HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}
		err = req.ValidatePath()
		if err != nil {
			t.Errorf("ValidatePath returned an error for valid path: %v", err)
		}
	})

	t.Run("Test Directory Path", func(t *testing.T) {
		requestString := "GET /files/ HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}
		err = req.ValidatePath()
		if err != nil {
			t.Errorf("ValidatePath returned an error for valid directory path: %v", err)
		}
	})

	t.Run("Test Invalid Path - Directory Traversal", func(t *testing.T) {
		requestString := "GET /files/../../../etc/passwd HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		_, err := NewRequest(requestString)
		if err == nil {
			t.Errorf("NewRequest returned an error: %v", err)
		}
	})
}

func TestFindQuery(t *testing.T) {
	tests := []struct {
		url   string
		query string
		valid bool
	}{
		{"/page?name=value", "name=value", true},
		{"/page?name=value&age=30", "name=value&age=30", true},
		{"/page", "", false},
		{"/page?", "", true},
		{"/page?name=", "name=", true},
		{"/page?=value", "=value", true},
		{"/page?name=value#section", "name=value", true},
		{"/page?name with space=value", "name with space=value", true},
		{"/page?name=value with space?", "name=value with space?", true},
		{"/page??doubleQuestion", "?doubleQuestion", true},
		{"/page?special=!@$%^&*()", "special=!@$%^&*()", true},
	}

	for _, test := range tests {
		t.Run("URL: "+test.url, func(t *testing.T) {
			query, err := findQuery(test.url)
			if test.valid {
				if err != nil {
					t.Errorf("findQuery(%q) failed: %v, but should have succeeded", test.url, err)
				}
				if query != test.query {
					t.Errorf("findQuery(%q) returned %q, expected %q", test.url, query, test.query)
				}
			} else {
				if err == nil {
					t.Errorf("findQuery(%q) succeeded, but should have failed", test.url)
				}
			}
		})
	}
}

func TestFindFragment(t *testing.T) {
	tests := []struct {
		url      string
		fragment string
		valid    bool
	}{
		{"/page#section1", "#section1", true},
		{"/page?query#section2", "#section2", true},
		{"/page", "", false},
		{"/page?query", "", false},
		{"/page#", "#", true},
		{"/page# ", "#", true},
		{"/page#section with space", "", false},
		{"/page?query#section with space", "", false},
		{"/page#", "#", true},
		{"/page#", "#", true},
	}

	for _, test := range tests {
		t.Run("URL: "+test.url, func(t *testing.T) {
			fragment, err := findFragment(test.url)
			// if err == nil {
			// 	t.Errorf("URL : %s, Expected Fragement : %s, Fragment found : %v", test.url, test.fragment, fragment)
			// }

			if test.fragment != "" {
				if err != nil {
					t.Errorf("findFragment(%q) failed: %v, but should have succeeded", test.url, err)
				}
				if string(fragment) != test.fragment {
					t.Errorf("findFragment(%q) returned %q, expected %q", test.url, fragment, test.fragment)
				}
			} else {
				if err == nil {
					t.Errorf("findFragment(%q) succeeded, but should have failed", test.url)
				}
			}
		})
	}
}

func TestParseFragment(t *testing.T) {
	tests := []struct {
		input    Fragment
		expected Fragment
		valid    bool
		errMsg   string
	}{
		{"#section1", "section1", true, ""},
		{"#header", "header", true, ""},
		{"", "", false, "fragment cannot be empty"},
		{"section1", "", false, "must start with '#'"},
		{"# section1", "%20section1", true, ""},
		{"#multiple#hash", "multiple%23hash", true, ""},
		{"#123", "123", true, ""},
		{"#special-chars_@$%^&*()", "special-chars_@$%^&*()", true, ""},
		{"#", "", true, ""},
		{"##doubleHash", "%23doubleHash", true, ""},
		{"# ", "%20", true, ""},
		{"#line\nbreak", "", false, "fragment cannot contain whitespace"},
		{"#line\nbreak", "", false, "fragment cannot contain whitespace"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ParseFragment(%q)", test.input), func(t *testing.T) {
			result, err := parseFragment(test.input)

			if test.valid {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				if result != test.expected {
					t.Errorf("expected %q, got %q", test.expected, result)
				}
			} else {
				if err == nil {
					t.Errorf("expected error but got success")
				}
				if test.errMsg != "" && err != nil && !strings.Contains(err.Error(), test.errMsg) {
					t.Errorf("expected error message containing %q, got %q", test.errMsg, err.Error())
				}
			}
		})
	}
}
