package support

import "fmt"

// ErrNotFound is a custom error type used when a resource is not found.
// This allows handlers to distinguish between a generic error and a "not found" condition,
// enabling them to return a 404 HTTP status code.
type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

// NewErrNotFound creates a new ErrNotFound error.
func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{Message: message}
}

// ErrBadRequest represents validation or client input errors (HTTP 400).
type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string               { return e.Message }
func NewErrBadRequest(message string) *ErrBadRequest { return &ErrBadRequest{Message: message} }

// ErrUnauthorized represents authentication failures (HTTP 401).
type ErrUnauthorized struct{ Message string }

func (e *ErrUnauthorized) Error() string                 { return e.Message }
func NewErrUnauthorized(message string) *ErrUnauthorized { return &ErrUnauthorized{Message: message} }

// ErrForbidden represents authorization failures (HTTP 403).
type ErrForbidden struct{ Message string }

func (e *ErrForbidden) Error() string              { return e.Message }
func NewErrForbidden(message string) *ErrForbidden { return &ErrForbidden{Message: message} }

// ErrTimeout represents request timeout to upstream or internal operations (HTTP 504 suggested).
type ErrTimeout struct{ Message string }

func (e *ErrTimeout) Error() string            { return e.Message }
func NewErrTimeout(message string) *ErrTimeout { return &ErrTimeout{Message: message} }

// ErrUpstream represents upstream dependency failures (HTTP 502/503 suggested).
type ErrUpstream struct {
	StatusCode int
	Body       string
}

func (e *ErrUpstream) Error() string {
	return fmt.Sprintf("upstream error status=%d body=%s", e.StatusCode, e.Body)
}
func NewErrUpstream(status int, body string) *ErrUpstream {
	return &ErrUpstream{StatusCode: status, Body: body}
}
