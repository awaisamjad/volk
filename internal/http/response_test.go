package http

import (
	"testing"
)

func TestResponseParsing(t *testing.T) {
	t.Run("Test Valid HTTP Response", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 22\r\n\r\n<h1>Hello World</h1>"
		_, err := parseResponse(responseString)
		if err != nil {
			t.Errorf("parseResponse returned an error: %v", err)
		}
	})

	t.Run("Test Response without body", func(t *testing.T) {
		responseString := "HTTP/1.1 204 No Content\r\nServer: TestServer\r\nDate: Mon, 27 Jul 2022 12:28:53 GMT\r\n\r\n"
		_, err := parseResponse(responseString)
		if err != nil {
			t.Errorf("parseResponse returned an error: %v", err)
		}
	})

	t.Run("Test Response with missing separator", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 22"
		_, err := parseResponse(responseString)
		if err == nil {
			t.Errorf("parseResponse should have returned an error due to missing separator")
		}
	})

	t.Run("Test Invalid Protocol", func(t *testing.T) {
		responseString := "HTTP/2.0 200 OK\r\nContent-Type: text/html\r\n\r\n<h1>Hello World</h1>"
		resp, err := parseResponse(responseString)
		if err != nil {
			t.Errorf("parseResponse should accept any protocol, got error: %v", err)
		}
		if resp.GetProtocol() != "HTTP/2.0" {
			t.Errorf("Expected protocol HTTP/2.0, got %s", resp.GetProtocol())
		}
	})

	t.Run("Test Invalid Status Code", func(t *testing.T) {
		responseString := "HTTP/1.1 999 Invalid\r\nContent-Type: text/html\r\n\r\n<h1>Error</h1>"
		resp, err := parseResponse(responseString)
		if err != nil {
			t.Errorf("parseResponse should accept any status code, got error: %v", err)
		}
		if resp.GetStatusCode() != 999 {
			t.Errorf("Expected status code 999, got %d", resp.GetStatusCode())
		}
	})

	t.Run("Test Invalid Header Format", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type text/html\r\n\r\n<h1>Hello World</h1>"
		_, err := parseResponse(responseString)
		if err == nil {
			t.Errorf("parseResponse should have returned an error due to invalid header format")
		}
	})
}

func TestNewResponse(t *testing.T) {
	t.Run("Test NewResponse with valid response string", func(t *testing.T) {
		responseString := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<h1>Hello World</h1>"
		resp, err := NewResponse(responseString)
		if err != nil {
			t.Errorf("NewResponse returned an error: %v", err)
		}
		if resp.GetStatusCode() != 200 {
			t.Errorf("Expected status code 200, got %d", resp.GetStatusCode())
		}
		if resp.GetStatusText() != "OK" {
			t.Errorf("Expected status text OK, got %s", resp.GetStatusText())
		}
	})

	t.Run("Test NewResponse with invalid response string", func(t *testing.T) {
		responseString := "INVALID RESPONSE"
		_, err := NewResponse(responseString)
		if err == nil {
			t.Errorf("Expected error for invalid response string, got nil")
		}
	})
}

func TestResponseGetters(t *testing.T) {
	resp := Response{
		StartLine: ResponseStartLine{
			Protocol:   HTTP1_1,
			StatusCode: 200,
			StatusText: "OK",
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/html"},
			{Name: "Content-Length", Value: "22"},
		},
		Body: "<h1>Hello World</h1>",
	}

	t.Run("GetProtocol", func(t *testing.T) {
		if resp.GetProtocol() != HTTP1_1 {
			t.Errorf("Expected HTTP/1.1, got %s", resp.GetProtocol())
		}
	})

	t.Run("GetStatusCode", func(t *testing.T) {
		if resp.GetStatusCode() != 200 {
			t.Errorf("Expected 200, got %d", resp.GetStatusCode())
		}
	})

	t.Run("GetStatusText", func(t *testing.T) {
		if resp.GetStatusText() != "OK" {
			t.Errorf("Expected OK, got %s", resp.GetStatusText())
		}
	})

	t.Run("GetHeaders", func(t *testing.T) {
		headers := resp.GetHeaders()
		if len(headers) != 2 {
			t.Errorf("Expected 2 headers, got %d", len(headers))
		}

		foundContentType := false
		foundContentLength := false
		for _, h := range headers {
			if h.Name == "Content-Type" && h.Value == "text/html" {
				foundContentType = true
			}
			if h.Name == "Content-Length" && h.Value == "22" {
				foundContentLength = true
			}
		}

		if !foundContentType {
			t.Errorf("Expected Content-Type header with value text/html")
		}
		if !foundContentLength {
			t.Errorf("Expected Content-Length header with value 22")
		}
	})

	t.Run("GetBody", func(t *testing.T) {
		if resp.GetBody() != "<h1>Hello World</h1>" {
			t.Errorf("Expected <h1>Hello World</h1>, got %s", resp.GetBody())
		}
	})
}

func TestResponseString(t *testing.T) {
	resp := Response{
		StartLine: ResponseStartLine{
			Protocol:   HTTP1_1,
			StatusCode: 200,
			StatusText: "OK",
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/html"},
			{Name: "Content-Length", Value: "22"},
		},
		Body: "<h1>Hello World</h1>",
	}

	expected := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 22\r\n\r\n<h1>Hello World</h1>"
	if resp.String() != expected {
		t.Errorf("Response.String() returned incorrect string.\nExpected: %q\nGot: %q", expected, resp.String())
	}
}

func TestResponseStartLineString(t *testing.T) {
	startLine := ResponseStartLine{
		Protocol:   HTTP1_1,
		StatusCode: 200,
		StatusText: "OK",
	}

	expected := "HTTP/1.1 200 OK"
	if startLine.String() != expected {
		t.Errorf("ResponseStartLine.String() returned %q, expected %q", startLine.String(), expected)
	}
}
