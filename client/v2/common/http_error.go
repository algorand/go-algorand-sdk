package common

import "fmt"

// HTTPError contains the raw HTTP error body plus parsed REST error fields when available.
type HTTPError struct {
	StatusCode int
	Message    string
	Data       map[string]any
	RawBody    []byte
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.RawBody)
}

// BadRequest represents a 400 HTTP response.
type BadRequest struct {
	*HTTPError
}

func (e BadRequest) Unwrap() error {
	return e.HTTPError
}

// InvalidToken represents a 401 HTTP response.
type InvalidToken struct {
	*HTTPError
}

func (e InvalidToken) Unwrap() error {
	return e.HTTPError
}

// NotFound represents a 404 HTTP response.
type NotFound struct {
	*HTTPError
}

func (e NotFound) Unwrap() error {
	return e.HTTPError
}

// InternalError represents a 500 HTTP response.
type InternalError struct {
	*HTTPError
}

func (e InternalError) Unwrap() error {
	return e.HTTPError
}
