// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound indicates the requested resource was not found (HTTP 404).
	ErrNotFound = errors.New("not found")

	// ErrAuthenticationRequired indicates authentication credentials are missing or invalid (HTTP 401).
	ErrAuthenticationRequired = errors.New("authentication required")

	// ErrRateLimitExceeded indicates the API rate limit has been exceeded (HTTP 429).
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrGatewayTimeout indicates the API gateway timed out (HTTP 504).
	ErrGatewayTimeout = errors.New("gateway timeout")
)

// APIError wraps HTTP error responses with status code and response body.
type APIError struct {
	StatusCode int
	Body       string
	Err        error
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("%s (HTTP %d): %s", e.Err, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("%s (HTTP %d)", e.Err, e.StatusCode)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// IsNotFound returns true if the error represents an HTTP 404 response.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
