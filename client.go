// Package refyne provides the official Go SDK for the Refyne API.
//
// Refyne is an LLM-powered web extraction API that transforms unstructured
// websites into clean, typed JSON data.
//
// Basic usage:
//
//	client := refyne.NewClient("your-api-key")
//
//	result, err := client.Extract(ctx, refyne.ExtractInput{
//	    URL:    "https://example.com/product",
//	    Schema: map[string]any{"name": "string", "price": "number"},
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Data)
package refyne

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Version constants
const (
	SDKVersion        = "0.0.1"
	DefaultBaseURL    = "https://api.refyne.uk"
	DefaultTimeout    = 30 * time.Second
	DefaultMaxRetries = 3
)

// Client is the main Refyne SDK client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
	logger Logger

	// Sub-clients for organized API access
	Jobs    *JobsClient
	Schemas *SchemasClient
	Sites   *SitesClient
	Keys    *KeysClient
	LLM     *LLMClient
}

// ClientOption configures the client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets the maximum retry attempts.
func WithMaxRetries(retries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Refyne client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{},
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
		logger:     &noopLogger{},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize sub-clients
	c.Jobs = &JobsClient{client: c}
	c.Schemas = &SchemasClient{client: c}
	c.Sites = &SitesClient{client: c}
	c.Keys = &KeysClient{client: c}
	c.LLM = &LLMClient{client: c}

	return c
}

// ExtractInput contains parameters for single-page extraction.
type ExtractInput struct {
	URL       string          `json:"url"`
	Schema    any             `json:"schema"`
	FetchMode *string         `json:"fetch_mode,omitempty"`
	LLMConfig *LLMConfigInput `json:"llm_config,omitempty"`
}

// Extract extracts structured data from a single web page.
func (c *Client) Extract(ctx context.Context, input ExtractInput) (*ExtractOutputBody, error) {
	var result ExtractOutputBody
	err := c.request(ctx, http.MethodPost, "/api/v1/extract", input, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CrawlInput contains parameters for starting a crawl job.
type CrawlInput struct {
	URL        string          `json:"url"`
	Schema     any             `json:"schema"`
	Options    *CrawlOptions   `json:"options,omitempty"`
	WebhookURL *string         `json:"webhook_url,omitempty"`
	LLMConfig  *LLMConfigInput `json:"llm_config,omitempty"`
}

// Crawl starts an asynchronous crawl job.
func (c *Client) Crawl(ctx context.Context, input CrawlInput) (*CrawlJobResponseBody, error) {
	var result CrawlJobResponseBody
	err := c.request(ctx, http.MethodPost, "/api/v1/crawl", input, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AnalyzeInput contains parameters for website analysis.
type AnalyzeInput struct {
	URL   string `json:"url"`
	Depth *int   `json:"depth,omitempty"`
}

// Analyze analyzes a website to detect structure and suggest schemas.
func (c *Client) Analyze(ctx context.Context, input AnalyzeInput) (*AnalyzeResponseBody, error) {
	var result AnalyzeResponseBody
	err := c.request(ctx, http.MethodPost, "/api/v1/analyze", input, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUsage returns usage statistics for the current billing period.
func (c *Client) GetUsage(ctx context.Context) (*GetUsageOutputBody, error) {
	var result GetUsageOutputBody
	err := c.request(ctx, http.MethodGet, "/api/v1/usage", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// request performs an HTTP request with retry logic.
func (c *Client) request(ctx context.Context, method, path string, body any, result any) error {
	return c.requestWithRetry(ctx, method, path, body, result, 1)
}

func (c *Client) requestWithRetry(ctx context.Context, method, path string, body any, result any, attempt int) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("refyne-go/%s", SDKVersion))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Retry on network errors
		if attempt <= c.maxRetries {
			backoff := c.calculateBackoff(attempt)
			c.logger.Warn("Network error, retrying", map[string]any{
				"error":   err.Error(),
				"attempt": attempt,
				"backoff": backoff,
			})
			time.Sleep(backoff)
			return c.requestWithRetry(ctx, method, path, body, result, attempt+1)
		}
		return &NetworkError{Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Handle rate limiting
	if resp.StatusCode == http.StatusTooManyRequests && attempt <= c.maxRetries {
		retryAfter := c.parseRetryAfter(resp.Header.Get("Retry-After"))
		c.logger.Warn("Rate limited, retrying", map[string]any{
			"retry_after": retryAfter,
			"attempt":     attempt,
		})
		time.Sleep(retryAfter)
		return c.requestWithRetry(ctx, method, path, body, result, attempt+1)
	}

	// Handle server errors with retry
	if resp.StatusCode >= 500 && attempt <= c.maxRetries {
		backoff := c.calculateBackoff(attempt)
		c.logger.Warn("Server error, retrying", map[string]any{
			"status":  resp.StatusCode,
			"attempt": attempt,
			"backoff": backoff,
		})
		time.Sleep(backoff)
		return c.requestWithRetry(ctx, method, path, body, result, attempt+1)
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		return c.parseError(resp.StatusCode, respBody)
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

func (c *Client) calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(1<<(attempt-1)) * time.Second
	if backoff > 30*time.Second {
		backoff = 30 * time.Second
	}
	return backoff
}

func (c *Client) parseRetryAfter(header string) time.Duration {
	if header == "" {
		return time.Second
	}
	if seconds, err := strconv.Atoi(header); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return time.Second
}

func (c *Client) parseError(status int, body []byte) error {
	var errResp struct {
		Error  string            `json:"error"`
		Detail string            `json:"detail"`
		Errors map[string]string `json:"errors"`
	}
	_ = json.Unmarshal(body, &errResp)

	msg := errResp.Error
	if msg == "" {
		msg = http.StatusText(status)
	}

	switch status {
	case http.StatusBadRequest:
		return &ValidationError{APIError: APIError{Message: msg, Status: status}, Fields: errResp.Errors}
	case http.StatusUnauthorized:
		return &AuthError{APIError: APIError{Message: msg, Status: status}}
	case http.StatusForbidden:
		return &ForbiddenError{APIError: APIError{Message: msg, Status: status}}
	case http.StatusNotFound:
		return &NotFoundError{APIError: APIError{Message: msg, Status: status}}
	case http.StatusTooManyRequests:
		return &RateLimitError{APIError: APIError{Message: msg, Status: status}}
	default:
		return &APIError{Message: msg, Status: status, Detail: errResp.Detail}
	}
}
