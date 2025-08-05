// Package http implements a simple HTTP server and related utilities.
package http

import (
	"fmt"
)

// RequestTarget represents an HTTP request target (path, query, fragment)
type RequestTarget struct {
	Path     string
	Query    string
	Fragment string
}

func (r RequestTarget) String() string {
	return fmt.Sprintf("%s%s%s", r.Path, r.Query, r.Fragment)
}

// parseRequestTarget extracts the path from a request target
func parseRequestTarget(requestTarget string) (string, error) {
	fragment, fragmentIdx, err := FindAndParseFragment(requestTarget)
	if err != nil && err != ErrFragmentNotFound {
		return "", err
	}

	query, queryIdx, err := FindAndParseQuery(requestTarget)
	if err != nil && err != ErrQueryNotFound && err != ErrQueryEmpty {
		return "", err
	}

	path := requestTarget

	if fragment != "" {
		path = requestTarget[:fragmentIdx]
	}

	if len(query.Params) > 0 {
		path = path[:queryIdx]
	}

	return path, nil
}
