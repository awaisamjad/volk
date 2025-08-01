// Package http implements a simple HTTP server and related utilities.
package http

import (
	"errors"
	"strings"
)

// Query errors
var (
	ErrQueryNotFound       = errors.New("query not found")
	ErrFragmentBeforeQuery = errors.New("fragment comes before query")
	ErrQueryEmpty          = errors.New("query is empty")
)

// Query represents HTTP query parameters
type Query struct {
	Params map[string][]string
}

// findQuery extracts the query string from a request target.
//
// Returns:
//   - query: The extracted query string, including the '?' prefix.
//   - queryIndex: The index of the '?' character in the requestTarget string.
//   - error: An error if the query is not found.  Will be ErrQueryNotFound if no query exists.
func findQuery(requestTarget string) (string, int, error) {
	query := ""

	queryIdx := strings.Index(requestTarget, "?")

	if queryIdx == -1 {
		return "", -1, ErrQueryNotFound
	}

	isThereAFragment := false
	fragmentIdx := strings.Index(requestTarget, "#")

	if fragmentIdx != -1 {
		isThereAFragment = true
	}

	if isThereAFragment {
		if fragmentIdx < queryIdx {
			return "", -1, ErrFragmentBeforeQuery
		}

		query = requestTarget[queryIdx:fragmentIdx]
	} else {
		query = requestTarget[queryIdx:]
	}

	if query == "" {
		return "", -1, ErrQueryEmpty
	}

	return query, queryIdx, nil
}

// parseQuery parses a query string into a Query struct.
//
// The query string should be in the format "key1=value1&key2=value2...".
// It handles cases where values are missing (e.g., "key1&key2=value2") by assigning an empty string.
//
// Returns:
//   - Query: A Query struct containing the parsed parameters.
//   - error: An error if parsing fails.
func parseQuery(query string) (Query, error) {
	if len(query) > 0 && query[0] == '?' {
		query = query[1:]
	}

	if query == "" {
		return Query{Params: map[string][]string{}}, ErrQueryEmpty
	}

	params := make(map[string][]string)
	// querySplit := strings.Split(query, "&")
	queryIter := strings.SplitSeq(query, "&")

	for param := range queryIter {
		// Skip empty parameters
		if param == "" {
			continue
		}

		keyValue := strings.SplitN(param, "=", 2)
		key := keyValue[0]

		// Handle case where there's no value
		if len(keyValue) == 1 {
			params[key] = append(params[key], "")
		} else {
			value := keyValue[1]
			params[key] = append(params[key], value)
		}
	}

	return Query{Params: params}, nil
}

// FindAndParseQuery extracts and parses the query string from a request target.
//
// It first uses findQuery to locate the query string. If a query string is found,
// it then uses parseQuery to parse it into a Query struct.
//
// Returns:
//   - Query: A Query struct containing the parsed parameters. If no query is found, an empty Query is returned.
//   - int: The index of the '?' character in the requestTarget string. Returns -1 if no query is found.
//   - error: An error if one occurred during the process. Returns nil if no error occurred or if no query was found.
func FindAndParseQuery(requestTarget string) (Query, int, error) {
	query, queryIndex, err := findQuery(requestTarget)
	if err != nil {
		// If no query is found, return an empty Query struct and no error.
		if errors.Is(err, ErrQueryNotFound) {
			return Query{Params: map[string][]string{}}, -1, nil
		}
		return Query{}, -1, err
	}
	parsedQuery, err := parseQuery(query)
	return parsedQuery, queryIndex, nil
}
