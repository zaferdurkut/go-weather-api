package support

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
