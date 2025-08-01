// Package http implements a simple HTTP server and related utilities.
package http

import (
	"errors"
	"strings"
)

// Fragment represents an HTTP URL fragment
type Fragment string

// Fragment errors
var (
	ErrFragmentNotFound     = errors.New("fragment not found")
	ErrFragmentEmpty        = errors.New("fragment cannot be empty")
	ErrFragmentNoHashPrefix = errors.New("fragment must start with '#'")
	ErrFragmentWhitespace   = errors.New("fragment cannot contain whitespace")
)

// findFragment extracts the fragment from a request target.
// It returns the fragment, its starting index in the request target, and an error if any.
//
// Returns:
//   - fragment: The extracted fragment string, including the '#' prefix.
//   - fragmentIndex: The index of the '#' character in the requestTarget string.
//   - error: An error if the fragment is not found.  Will be ErrFragmentNotFound if no fragment exists.
func findFragment(requestTarget string) (Fragment, int, error) {
	fragment := ""

	fragmentIdx := strings.Index(requestTarget, "#")

	if fragmentIdx == -1 {
		return "", -1, ErrFragmentNotFound
	}

	for i := fragmentIdx; i < len(requestTarget); i++ {
		fragment += string(requestTarget[i])
	}

	return Fragment(fragment), fragmentIdx, nil
}

// parseFragment validates and parses a fragment.
//
// It performs the following checks:
//   - The fragment is not empty.
//   - The fragment starts with a '#'.
//   - The fragment does not contain whitespace characters (tabs, carriage returns, or newlines).
//
// Returns:
//   - fragment: The validated fragment.
//   - error: An error if the fragment is invalid.s a fragment
func parseFragment(fragment Fragment) (Fragment, error) {
	if fragment == "" {
		return "", ErrFragmentEmpty
	}

	if !strings.HasPrefix(string(fragment), "#") {
		return "", ErrFragmentNoHashPrefix
	}

	if strings.Contains(string(fragment), "\t") || strings.Contains(string(fragment), "\r") || strings.Contains(string(fragment), "\n") {
		return "", ErrFragmentWhitespace
	}

	return fragment, nil
}

// FindAndParseFragment finds and parses a fragment from a request target.
//
// It combines the functionality of findFragment and parseFragment.
// It first attempts to find the fragment within the request target.
// If a fragment is found, it then validates and parses the fragment.
//
// Returns:
//   - fragment: The validated fragment.
//   - fragmentIndex: The index of the '#' character in the requestTarget string. Will return -1 if not found.
//   - error: An error if any of the find or parse operations fail.
//     Possible errors include ErrFragmentNotFound, ErrFragmentEmpty, ErrFragmentNoHashPrefix, and ErrFragmentWhitespace.
func FindAndParseFragment(requestTarget string) (Fragment, int, error) {
	fragment, fragmentIdx, err := findFragment(requestTarget)
	if err != nil {
		return "", -1, err
	}

	parsedFragment, err := parseFragment(fragment)
	return parsedFragment, fragmentIdx, err
}
