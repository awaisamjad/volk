package http

import (
	"reflect"
	"testing"
)

func TestFindQuery(t *testing.T) {
	tests := []struct {
		url   string
		query string
		valid bool
	}{
		{"/page?name=value ", "?name=value", true},
		{"/page?name=value&age=30", "?name=value&age=30", true},
		{"/page", "", false},
		{"/page?", "?", true},
		{"/page?name=", "?name=", true},
		{"/page?=value", "?=value", true},
		{"/page?name=value#section", "?name=value", true},
		{"/page?name with space=value # fdsfds", "?name with space=value", true},
		{"/page?name=value with space?", "?name=value with space?", true},
		{"/page??doubleQuestion", "??doubleQuestion", true},
		{"/page?special=!@$%^&*()", "?special=!@$%^&*()", true},
	}

	for _, test := range tests {
		t.Run("URL: "+test.url, func(t *testing.T) {
			query, _, err := findQuery(test.url)
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

func TestParseQuery(t *testing.T) {
	compareQueries := func(q1, q2 Query) bool {
		if len(q1.Params) != len(q2.Params) {
			return false
		}
		for key, values1 := range q1.Params {
			values2, ok := q2.Params[key]
			if !ok || len(values1) != len(values2) {
				return false
			}
			for i, v1 := range values1 {
				if v1 != values2[i] {
					return false
				}
			}
		}
		return true
	}

	tests := []struct {
		url     string
		queries Query
		valid   bool
	}{
		{
			"/page?name=value",
			Query{Params: map[string][]string{"name": {"value"}}},
			true,
		},
		{
			"/page?name=value&age=30",
			Query{Params: map[string][]string{"name": {"value"}, "age": {"30"}}},
			true,
		},
		{
			"/page?",
			Query{Params: map[string][]string{}},
			true,
		},
		{
			"/page?name=",
			Query{Params: map[string][]string{"name": {""}}},
			true,
		},
		{
			"/page?=value",
			Query{Params: map[string][]string{"": {"value"}}},
			true,
		},
		{
			"/page?name=value#section",
			Query{Params: map[string][]string{"name": {"value"}}},
			true,
		},
		{
			"/page?name%20with%20space=value",
			Query{Params: map[string][]string{"name with space": {"value"}}},
			true,
		},
		{
			"/page?name=value%20with%20space?",
			Query{Params: map[string][]string{"name": {"value with space?"}}},
			true,
		},
		{
			"/page??doubleQuestion",
			Query{Params: map[string][]string{"?doubleQuestion": {""}}},
			true,
		},
		{
			"/page?name=value&name=another",
			Query{Params: map[string][]string{"name": {"value", "another"}}},
			true,
		},
	}

	for _, test := range tests {
		t.Run("URL: "+test.url, func(t *testing.T) {
			foundQuery, _, err := findQuery(test.url)
			if err != nil {
				t.Errorf("findQuery(%q) failed: %v, but should have succeeded", test.url, err)
			}

			parsedQuery, err := parseQuery(foundQuery)
			if test.valid {
				if err != nil {
					t.Errorf("parseQuery(%q) failed: %v, but should have succeeded", test.url, err)
				}
				if !compareQueries(parsedQuery, test.queries) {
					t.Errorf("parseQuery(%q) returned %v, expected %v", test.url, parsedQuery, test.queries)
				}
			} else {
				if err == nil {
					t.Errorf("parseQuery(%q) succeeded, but should have failed", test.url)
				}
			}
		})
	}
}

func TestFindAndParseQuery(t *testing.T) {
	tests := []struct {
		name           string
		requestTarget  string
		expectedQuery  Query
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Simple query",
			requestTarget: "/page?name=value",
			expectedQuery: Query{Params: map[string][]string{"name": {"value"}}},
			expectedError: false,
		},
		{
			name:          "No query",
			requestTarget: "/page",
			expectedQuery: Query{Params: map[string][]string{}},
			expectedError: false,
		},
		{
			name:          "Multiple parameters",
			requestTarget: "/page?name=value&age=30",
			expectedQuery: Query{Params: map[string][]string{"name": {"value"}, "age": {"30"}}},
			expectedError: false,
		},
		{
			name:          "Empty query",
			requestTarget: "/page?",
			expectedQuery: Query{Params: map[string][]string{}},
			expectedError: false,
		},
		{
			name:          "Query with fragment",
			requestTarget: "/page?name=value#section",
			expectedQuery: Query{Params: map[string][]string{"name": {"value"}}},
			expectedError: false,
		},
		{
			name:          "Duplicate parameter",
			requestTarget: "/page?name=value&name=another",
			expectedQuery: Query{Params: map[string][]string{"name": {"value", "another"}}},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := FindAndParseQuery(tt.requestTarget)

			if tt.expectedError {
				if err == nil {
					t.Errorf("FindAndParseQuery(%q) succeeded, but should have failed", tt.requestTarget)
				}
			} else {
				if err != nil {
					t.Errorf("FindAndParseQuery(%q) failed: %v, but should have succeeded", tt.requestTarget, err)
				}

				if !reflect.DeepEqual(query, tt.expectedQuery) {
					t.Errorf("FindAndParseQuery(%q) returned %v, expected %v", tt.requestTarget, query, tt.expectedQuery)
				}
			}
		})
	}
}
