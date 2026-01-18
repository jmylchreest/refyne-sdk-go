package refyne

import "fmt"

// RefyneError is the base error type for all SDK errors.
type RefyneError struct {
	// Message is the error message.
	Message string
	// Status is the HTTP status code.
	Status int
	// Detail is additional error detail.
	Detail string
}

// Error implements the error interface.
func (e *RefyneError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// RateLimitError is returned when rate limited.
type RateLimitError struct {
	RefyneError
	// RetryAfter is seconds to wait before retrying.
	RetryAfter int
}

// ValidationError is returned when validation fails.
type ValidationError struct {
	RefyneError
	// Errors contains field-level errors.
	Errors map[string]string
}

// AuthenticationError is returned when authentication fails.
type AuthenticationError struct {
	RefyneError
}

// ForbiddenError is returned when access is forbidden.
type ForbiddenError struct {
	RefyneError
}

// NotFoundError is returned when a resource is not found.
type NotFoundError struct {
	RefyneError
}

// UnsupportedAPIVersionError is returned when API version is incompatible.
type UnsupportedAPIVersionError struct {
	// APIVersion is the detected API version.
	APIVersion string
	// MinVersion is the minimum supported version.
	MinVersion string
	// MaxKnownVersion is the maximum known version.
	MaxKnownVersion string
}

// Error implements the error interface.
func (e *UnsupportedAPIVersionError) Error() string {
	return fmt.Sprintf(
		"API version %s is not supported. This SDK requires API version >= %s. "+
			"Please upgrade the API or use an older SDK version.",
		e.APIVersion, e.MinVersion,
	)
}
