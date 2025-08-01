package http

import (
	"fmt"
	"strings"
	"testing"
)

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
		{"/page#section with space", "#section with space", true},
		{"/page?query#section with space", "#section with space", true},
		{"/page#", "#", true},
		{"/page#", "#", true},
	}

	for _, test := range tests {
		t.Run("URL: "+test.url, func(t *testing.T) {
			fragment, _, err := findFragment(test.url)
			if test.valid {
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
		{"#section1", "#section1", true, ""},
		{"#header", "#header", true, ""},
		{"", "", false, "fragment cannot be empty"},
		{"section1", "", false, "must start with '#'"},
		{"# section1", "# section1", true, ""},
		{"#multiple#hash", "#multiple#hash", true, ""},
		{"#123", "#123", true, ""},
		{"#special-chars_@$%^&*()", "#special-chars_@$%^&*()", true, ""},
		{"#", "#", true, ""},
		{"##doubleHash", "##doubleHash", true, ""},
		{"# ", "# ", true, ""},
		{"#line\nbreak", "", false, "fragment cannot contain whitespace"},
		{"#line\rbreak", "", false, "fragment cannot contain whitespace"},
		{"#line\tbreak", "", false, "fragment cannot contain whitespace"},
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

func TestFindAndParseFragment(t *testing.T) {
	tests := []struct {
		name             string
		requestTarget    string
		expectedFragment Fragment
		expectedError    bool
		expectedErrMsg   string
	}{
		{
			name:             "Simple fragment",
			requestTarget:    "/page#section1",
			expectedFragment: "#section1",
			expectedError:    false,
		},
		{
			name:             "No fragment",
			requestTarget:    "/page",
			expectedFragment: "",
			expectedError:    true,
		},
		{
			name:             "Fragment with query",
			requestTarget:    "/page?name=value#section1",
			expectedFragment: "#section1",
			expectedError:    false,
		},
		{
			name:             "Empty fragment",
			requestTarget:    "/page#",
			expectedFragment: "#",
			expectedError:    false,
		},
		{
			name:             "Fragment with spaces",
			requestTarget:    "/page#section with spaces",
			expectedFragment: "#section with spaces",
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fragment, _, err := FindAndParseFragment(tt.requestTarget)

			if tt.expectedError {
				if err == nil {
					t.Errorf("FindAndParseFragment(%q) succeeded, but should have failed", tt.requestTarget)
				}
			} else {
				if err != nil {
					t.Errorf("FindAndParseFragment(%q) failed: %v, but should have succeeded", tt.requestTarget, err)
				}

				if fragment != tt.expectedFragment {
					t.Errorf("FindAndParseFragment(%q) returned %q, expected %q", tt.requestTarget, fragment, tt.expectedFragment)
				}
			}
		})
	}
}
