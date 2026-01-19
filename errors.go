package refyne

import "fmt"

// Logger is the interface for custom logging.
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

// noopLogger is the default logger that does nothing.
type noopLogger struct{}

func (n *noopLogger) Debug(msg string, fields map[string]any) {}
func (n *noopLogger) Info(msg string, fields map[string]any)  {}
func (n *noopLogger) Warn(msg string, fields map[string]any)  {}
func (n *noopLogger) Error(msg string, fields map[string]any) {}

// APIError is the base error type for API errors.
type APIError struct {
	Message string
	Status  int
	Detail  string
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// ValidationError is returned when request validation fails.
type ValidationError struct {
	APIError
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// AuthError is returned when authentication fails.
type AuthError struct {
	APIError
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

// ForbiddenError is returned when access is denied.
type ForbiddenError struct {
	APIError
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Message)
}

// NotFoundError is returned when a resource is not found.
type NotFoundError struct {
	APIError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

// RateLimitError is returned when rate limit is exceeded.
type RateLimitError struct {
	APIError
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}

// NetworkError is returned when a network error occurs.
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}
