package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// StringToURL converts a raw string to a *url.URL.
// It trims leading/trailing whitespace and returns an error if the string is empty or malformed.
func StringToURL(rawURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return nil, errors.New("url string cannot be empty")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %q: %w", rawURL, err)
	}
	return parsed, nil
}

// ParseURL is an alias for StringToURL for standard naming convention.
func ParseURL(rawURL string) (*url.URL, error) {
	return StringToURL(rawURL)
}

// MustStringToURL converts a raw string to a *url.URL, panicking if parsing fails.
// Ideal for test fixtures, constants, or setup where invalid URLs indicate a programming bug.
func MustStringToURL(rawURL string) *url.URL {
	u, err := StringToURL(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

// MustParseURL is an alias for MustStringToURL.
func MustParseURL(rawURL string) *url.URL {
	return MustStringToURL(rawURL)
}

// JoinURL combines a base URL string and one or more relative path segments into a single *url.URL.
// It normalizes path slashes cleanly.
func JoinURL(baseURL string, elem ...string) (*url.URL, error) {
	u, err := StringToURL(baseURL)
	if err != nil {
		return nil, err
	}
	return u.JoinPath(elem...), nil
}
