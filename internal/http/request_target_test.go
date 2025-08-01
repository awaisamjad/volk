package http

import (
	"testing"
)

func TestRequestTargetString(t *testing.T) {
	tests := []struct {
		name     string
		target   RequestTarget
		expected string
	}{
		{
			name: "Simple path",
			target: RequestTarget{
				Path:     "/index.html",
				Query:    "",
				Fragment: "",
			},
			expected: "/index.html",
		},
		{
			name: "Path with query",
			target: RequestTarget{
				Path:     "/search",
				Query:    "?q=keyword",
				Fragment: "",
			},
			expected: "/search?q=keyword",
		},
		{
			name: "Path with fragment",
			target: RequestTarget{
				Path:     "/about",
				Query:    "",
				Fragment: "#team",
			},
			expected: "/about#team",
		},
		{
			name: "Path with query and fragment",
			target: RequestTarget{
				Path:     "/search",
				Query:    "?q=keyword",
				Fragment: "#results",
			},
			expected: "/search?q=keyword#results",
		},
		{
			name: "Root path",
			target: RequestTarget{
				Path:     "/",
				Query:    "",
				Fragment: "",
			},
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.target.String()
			if result != tt.expected {
				t.Errorf("RequestTarget.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParseRequestTarget(t *testing.T) {
	tests := []struct {
		name          string
		requestTarget string
		expectedPath  string
		expectedError bool
	}{
		{
			name:          "Simple path",
			requestTarget: "/index.html",
			expectedPath:  "/index.html",
			expectedError: false,
		},
		{
			name:          "Path with query",
			requestTarget: "/search?q=keyword",
			expectedPath:  "/search",
			expectedError: false,
		},
		{
			name:          "Path with fragment",
			requestTarget: "/about#team",
			expectedPath:  "/about",
			expectedError: false,
		},
		{
			name:          "Path with query and fragment",
			requestTarget: "/search?q=keyword#results",
			expectedPath:  "/search",
			expectedError: false,
		},
		{
			name:          "Root path",
			requestTarget: "/",
			expectedPath:  "/",
			expectedError: false,
		},
		{
			name:          "Path with special characters",
			requestTarget: "/path/!@#$%^&*()/resource",
			expectedPath:  "/path/!@#$%^&*()/resource",
			expectedError: false,
		},
		{
			name:          "Path with space",
			requestTarget: "/path with space/resource",
			expectedPath:  "/path with space/resource",
			expectedError: false,
		},
		{
			name:          "Path with encoded characters",
			requestTarget: "/path/%20/resource",
			expectedPath:  "/path/%20/resource",
			expectedError: false,
		},
		{
			name:          "Empty path",
			requestTarget: "",
			expectedPath:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := parseRequestTarget(tt.requestTarget)

			if tt.expectedError {
				if err == nil {
					t.Errorf("parseRequestTarget(%q) succeeded, but should have failed", tt.requestTarget)
				}
			} else {
				if err != nil {
					t.Errorf("parseRequestTarget(%q) failed: %v, but should have succeeded", tt.requestTarget, err)
				}

				if path != tt.expectedPath {
					t.Errorf("parseRequestTarget(%q) returned path %q, expected %q", tt.requestTarget, path, tt.expectedPath)
				}
			}
		})
	}
}
