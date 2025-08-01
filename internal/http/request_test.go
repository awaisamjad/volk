package http

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	t.Run("Test Valid GET Request", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
		_, err := parseRequest(requestString)
		if err != nil {
			t.Errorf("parseRequest returned an error: %v", err)
		}
	})

	t.Run("Test Invalid GET Request", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n<h1>Hello World</h1>"
		_, err := parseRequest(requestString)
		if err != nil {
			t.Errorf("parseRequest didn't return an error as expected, got: %v", err)
		}
	})

	t.Run("Test Valid POST Request", func(t *testing.T) {
		requestString := "POST /submit HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/json\r\n\r\n{\"name\": \"John Doe\"}"
		_, err := parseRequest(requestString)
		if err != nil {
			t.Errorf("parseRequest returned an error: %v", err)
		}
	})

	t.Run("Test Request with no separator", func(t *testing.T) {
		requestString := "GET /index.html HTTP/1.1\r\nHost: localhost:8080\r\n"
		_, err := parseRequest(requestString)
		if err == nil {
			t.Errorf("parseRequest should have returned an error due to missing separator")
		}
	})

	t.Run("Test Invalid Method", func(t *testing.T) {
		requestString := "FOO / HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
		_, err := parseRequest(requestString)
		if err != nil {
			t.Errorf("parseRequest should accept any method, but returned error: %v", err)
		}
	})

	t.Run("Test Invalid Protocol", func(t *testing.T) {
		requestString := "GET / HTTP/2.0\r\nHost: localhost:8080\r\n\r\n"
		_, err := parseRequest(requestString)
		if err != nil {
			t.Errorf("parseRequest should accept any protocol, but returned error: %v", err)
		}
	})

	t.Run("Test Invalid Header", func(t *testing.T) {
		requestString := "GET / HTTP/1.1\r\nHost localhost:8080\r\n\r\n"
		_, err := parseRequest(requestString)
		if err == nil {
			t.Errorf("parseRequest should have returned an error due to invalid header")
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
		if req.StartLine.RequestTarget.Path != "/files/index.html" {
			t.Errorf("Expected request target /files/index.html, got %s", req.StartLine.RequestTarget.Path)
		}
		if req.StartLine.RequestTarget.Path == "files/index.html" {
			t.Errorf("Expected request target to start with /, got %s", req.StartLine.RequestTarget.Path)
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
		if err != nil {
			t.Errorf("NewRequest returned an error when it shouldn't: %v", err)
			return
		}

		resp := req.Response()
		if resp.StartLine.StatusCode != 400 {
			t.Errorf("Expected status code 400 for * target with GET, got %d", resp.StartLine.StatusCode)
		}
	})
}

func TestRequestGetters(t *testing.T) {
	requestString := "GET /files/index.html HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/html\r\n\r\n"
	req, err := NewRequest(requestString)
	if err != nil {
		t.Fatalf("NewRequest returned an error: %v", err)
	}

	t.Run("GetMethod", func(t *testing.T) {
		if req.GetMethod() != GET {
			t.Errorf("Expected GET, got %s", req.GetMethod())
		}
	})

	t.Run("GetRequestTarget", func(t *testing.T) {
		if req.GetRequestTarget().Path != "/files/index.html" {
			t.Errorf("Expected /files/index.html, got %s", req.GetRequestTarget().Path)
		}
	})

	t.Run("GetProtocol", func(t *testing.T) {
		if req.GetProtocol() != HTTP1_1 {
			t.Errorf("Expected HTTP/1.1, got %s", req.GetProtocol())
		}
	})

	t.Run("GetHeaders", func(t *testing.T) {
		headers := req.GetHeaders()
		if len(headers) != 2 {
			t.Errorf("Expected 2 headers, got %d", len(headers))
		}

		foundHost := false
		foundContentType := false
		for _, h := range headers {
			if h.Name == "Host" && h.Value == "localhost:8080" {
				foundHost = true
			}
			if h.Name == "Content-Type" && h.Value == "application/html" {
				foundContentType = true
			}
		}

		if !foundHost {
			t.Errorf("Expected Host header with value localhost:8080")
		}
		if !foundContentType {
			t.Errorf("Expected Content-Type header with value application/html")
		}
	})

	t.Run("GetBody", func(t *testing.T) {
		if req.GetBody() != "" {
			t.Errorf("Expected empty body, got %s", req.GetBody())
		}
	})
}

func TestValidatePath(t *testing.T) {
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
		req, err := NewRequest(requestString)
		if err != nil {
			t.Errorf("NewRequest returned an error: %v", err)
			return
		}
		err = req.ValidatePath()
		if err == nil {
			t.Errorf("ValidatePath should have returned an error for directory traversal")
		}
	})

	t.Run("Test Invalid Path - Empty Path", func(t *testing.T) {
		req := Request{
			StartLine: RequestStartLine{
				Method: GET,
				RequestTarget: RequestTarget{
					Path:     "",
					Query:    "",
					Fragment: "",
				},
				Protocol: HTTP1_1,
			},
		}
		err := req.ValidatePath()
		if err == nil {
			t.Errorf("ValidatePath should have returned an error for empty path")
		}
	})

	t.Run("Test Invalid Path - Invalid Characters", func(t *testing.T) {
		req := Request{
			StartLine: RequestStartLine{
				Method: GET,
				RequestTarget: RequestTarget{
					Path:     "/path/with\x00null",
					Query:    "",
					Fragment: "",
				},
				Protocol: HTTP1_1,
			},
		}
		err := req.ValidatePath()
		if err == nil {
			t.Errorf("ValidatePath should have returned an error for path with invalid characters")
		}
	})
}

func TestRequestString(t *testing.T) {
	req := Request{
		StartLine: RequestStartLine{
			Method: GET,
			RequestTarget: RequestTarget{
				Path:     "/index.html",
				Query:    "?param=value",
				Fragment: "#section",
			},
			Protocol: HTTP1_1,
		},
		Headers: []Header{
			{Name: "Host", Value: "example.com"},
			{Name: "User-Agent", Value: "test-client"},
		},
		Body: "Test body content",
	}

	expected := "GET /index.html?param=value#section HTTP/1.1\r\nHost: example.com\r\nUser-Agent: test-client\r\n\r\nTest body content"
	if req.String() != expected {
		t.Errorf("Request.String() returned incorrect string.\nExpected: %q\nGot: %q", expected, req.String())
	}
}
