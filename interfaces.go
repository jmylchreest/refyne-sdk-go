package refyne

import "net/http"

// Logger defines the interface for SDK logging.
// Implement this interface to provide custom logging behavior.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, meta map[string]any)
	// Info logs an info message.
	Info(msg string, meta map[string]any)
	// Warn logs a warning message.
	Warn(msg string, meta map[string]any)
	// Error logs an error message.
	Error(msg string, meta map[string]any)
}

// noopLogger is the default logger that does nothing.
type noopLogger struct{}

func (l *noopLogger) Debug(msg string, meta map[string]any) {}
func (l *noopLogger) Info(msg string, meta map[string]any)  {}
func (l *noopLogger) Warn(msg string, meta map[string]any)  {}
func (l *noopLogger) Error(msg string, meta map[string]any) {}

// HTTPClient defines the interface for HTTP operations.
// Implement this interface to provide custom HTTP behavior.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// defaultHTTPClient wraps the standard http.Client.
type defaultHTTPClient struct {
	client *http.Client
}

func (c *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Cache defines the interface for caching API responses.
// Implement this interface to provide custom caching behavior.
type Cache interface {
	// Get retrieves a cached entry by key.
	Get(key string) (*CacheEntry, bool)
	// Set stores an entry in the cache.
	Set(key string, entry *CacheEntry)
	// Delete removes an entry from the cache.
	Delete(key string)
}

// CacheEntry represents a cached response.
type CacheEntry struct {
	// Value is the cached response data.
	Value any
	// ExpiresAt is the Unix timestamp when entry expires.
	ExpiresAt int64
	// CacheControl contains the parsed directives.
	CacheControl CacheControlDirectives
}

// CacheControlDirectives contains parsed Cache-Control header values.
type CacheControlDirectives struct {
	NoStore              bool
	NoCache              bool
	Private              bool
	MaxAge               *int
	StaleWhileRevalidate *int
}
