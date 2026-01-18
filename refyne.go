// Package refyne provides the official Go SDK for the Refyne API.
//
// Refyne is an LLM-powered web extraction API that transforms unstructured
// websites into clean, typed JSON data.
//
// Basic usage:
//
//	client := refyne.NewClient("your-api-key")
//
//	result, err := client.Extract(ctx, refyne.ExtractRequest{
//	    URL: "https://example.com/product",
//	    Schema: map[string]any{
//	        "name":  "string",
//	        "price": "number",
//	    },
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
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Version information
const (
	// SDKVersion is the current SDK version
	SDKVersion = "0.0.0"
	// MinAPIVersion is the minimum API version this SDK supports
	MinAPIVersion = "0.0.0"
	// MaxKnownAPIVersion is the maximum API version this SDK was built against
	MaxKnownAPIVersion = "0.0.0"
)

// DefaultBaseURL is the default API base URL
const DefaultBaseURL = "https://api.refyne.uk"

// Client is the main Refyne SDK client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient HTTPClient
	logger     Logger
	cache      Cache
	cacheOn    bool
	userAgent  string
	timeout    time.Duration
	maxRetries int

	apiVersionChecked bool
	authHash          string

	mu sync.RWMutex

	// Sub-clients
	Jobs    *JobsService
	Schemas *SchemasService
	Sites   *SitesService
	Keys    *KeysService
	LLM     *LLMService
}

// Option configures a Client
type Option func(*Client)

// WithBaseURL sets a custom base URL
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithLogger sets a custom logger
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithCache sets a custom cache
func WithCache(cache Cache) Option {
	return func(c *Client) {
		c.cache = cache
	}
}

// WithCacheEnabled enables or disables caching
func WithCacheEnabled(enabled bool) Option {
	return func(c *Client) {
		c.cacheOn = enabled
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets the maximum retry attempts
func WithMaxRetries(retries int) Option {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

// WithUserAgentSuffix appends to the User-Agent string
func WithUserAgentSuffix(suffix string) Option {
	return func(c *Client) {
		c.userAgent = buildUserAgent(suffix)
	}
}

// NewClient creates a new Refyne client with the given API key and options.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: &defaultHTTPClient{client: &http.Client{}},
		logger:     &noopLogger{},
		cache:      NewMemoryCache(100),
		cacheOn:    true,
		userAgent:  buildUserAgent(""),
		timeout:    30 * time.Second,
		maxRetries: 3,
		authHash:   hashString(apiKey),
	}

	for _, opt := range opts {
		opt(c)
	}

	// Warn about insecure connections
	if !strings.HasPrefix(c.baseURL, "https://") {
		c.logger.Warn("API base URL is not using HTTPS. This is insecure.", map[string]any{
			"baseURL": c.baseURL,
		})
	}

	// Initialize sub-services
	c.Jobs = &JobsService{client: c}
	c.Schemas = &SchemasService{client: c}
	c.Sites = &SitesService{client: c}
	c.Keys = &KeysService{client: c}
	c.LLM = &LLMService{client: c}

	return c
}

// Extract extracts structured data from a single web page.
func (c *Client) Extract(ctx context.Context, req ExtractRequest) (*ExtractResponse, error) {
	body := map[string]any{
		"url":    req.URL,
		"schema": req.Schema,
	}
	if req.FetchMode != "" {
		body["fetchMode"] = req.FetchMode
	}
	if req.LLMConfig != nil {
		body["llmConfig"] = req.LLMConfig
	}

	var resp ExtractResponse
	if err := c.request(ctx, http.MethodPost, "/api/v1/extract", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Crawl starts an asynchronous crawl job.
func (c *Client) Crawl(ctx context.Context, req CrawlRequest) (*CrawlJobCreated, error) {
	body := map[string]any{
		"url":    req.URL,
		"schema": req.Schema,
	}
	if req.Options != nil {
		body["options"] = req.Options
	}
	if req.WebhookURL != "" {
		body["webhookUrl"] = req.WebhookURL
	}
	if req.LLMConfig != nil {
		body["llmConfig"] = req.LLMConfig
	}

	var resp CrawlJobCreated
	if err := c.request(ctx, http.MethodPost, "/api/v1/crawl", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Analyze analyzes a website to detect structure and suggest schemas.
func (c *Client) Analyze(ctx context.Context, req AnalyzeRequest) (*AnalyzeResponse, error) {
	body := map[string]any{
		"url": req.URL,
	}
	if req.Depth > 0 {
		body["depth"] = req.Depth
	}

	var resp AnalyzeResponse
	if err := c.request(ctx, http.MethodPost, "/api/v1/analyze", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsage returns usage statistics for the current billing period.
func (c *Client) GetUsage(ctx context.Context) (*UsageResponse, error) {
	var resp UsageResponse
	if err := c.request(ctx, http.MethodGet, "/api/v1/usage", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) request(ctx context.Context, method, path string, body any, result any) error {
	return c.requestWithOptions(ctx, method, path, body, result, false)
}

func (c *Client) requestWithOptions(ctx context.Context, method, path string, body any, result any, skipCache bool) error {
	urlStr := c.baseURL + path
	cacheKey := GenerateCacheKey(method, urlStr, c.authHash)

	// Check cache for GET requests
	if method == http.MethodGet && c.cacheOn && !skipCache {
		if entry, ok := c.cache.Get(cacheKey); ok {
			if data, err := json.Marshal(entry.Value); err == nil {
				return json.Unmarshal(data, result)
			}
		}
	}

	resp, err := c.executeWithRetry(ctx, method, urlStr, body, 1)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check API version on first request
	c.mu.Lock()
	if !c.apiVersionChecked {
		if apiVersion := resp.Header.Get("X-API-Version"); apiVersion != "" {
			if err := CheckAPIVersionCompatibility(apiVersion, c.logger); err != nil {
				c.mu.Unlock()
				return err
			}
		} else {
			c.logger.Warn("API did not return X-API-Version header", nil)
		}
		c.apiVersionChecked = true
	}
	c.mu.Unlock()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp, respBody)
	}

	// Parse response
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache GET responses
	if method == http.MethodGet && c.cacheOn {
		cacheControl := resp.Header.Get("Cache-Control")
		if entry := CreateCacheEntry(result, cacheControl); entry != nil {
			c.cache.Set(cacheKey, entry)
		}
	}

	return nil
}

func (c *Client) executeWithRetry(ctx context.Context, method, urlStr string, body any, attempt int) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Retry on network errors
		if attempt <= c.maxRetries {
			backoff := min(time.Duration(1<<(attempt-1))*time.Second, 30*time.Second)
			c.logger.Warn("Network error. Retrying", map[string]any{
				"error":      err.Error(),
				"attempt":    attempt,
				"maxRetries": c.maxRetries,
				"backoff":    backoff,
			})
			time.Sleep(backoff)
			return c.executeWithRetry(ctx, method, urlStr, body, attempt+1)
		}
		return nil, fmt.Errorf("network error: %w", err)
	}

	// Handle rate limiting with retry
	if resp.StatusCode == http.StatusTooManyRequests && attempt <= c.maxRetries {
		retryAfter := 1
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if v, err := strconv.Atoi(ra); err == nil {
				retryAfter = v
			}
		}
		c.logger.Warn("Rate limited. Retrying", map[string]any{
			"retryAfter": retryAfter,
			"attempt":    attempt,
			"maxRetries": c.maxRetries,
		})
		resp.Body.Close()
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return c.executeWithRetry(ctx, method, urlStr, body, attempt+1)
	}

	// Handle server errors with retry
	if resp.StatusCode >= 500 && attempt <= c.maxRetries {
		backoff := min(time.Duration(1<<(attempt-1))*time.Second, 30*time.Second)
		c.logger.Warn("Server error. Retrying", map[string]any{
			"status":     resp.StatusCode,
			"attempt":    attempt,
			"maxRetries": c.maxRetries,
			"backoff":    backoff,
		})
		resp.Body.Close()
		time.Sleep(backoff)
		return c.executeWithRetry(ctx, method, urlStr, body, attempt+1)
	}

	return resp, nil
}

func parseErrorResponse(resp *http.Response, body []byte) error {
	var errResp struct {
		Error  string            `json:"error"`
		Detail string            `json:"detail"`
		Errors map[string]string `json:"errors"`
	}
	_ = json.Unmarshal(body, &errResp)

	msg := errResp.Error
	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return &ValidationError{RefyneError: RefyneError{Message: msg, Status: 400}, Errors: errResp.Errors}
	case http.StatusUnauthorized:
		return &AuthenticationError{RefyneError: RefyneError{Message: msg, Status: 401}}
	case http.StatusForbidden:
		return &ForbiddenError{RefyneError: RefyneError{Message: msg, Status: 403}}
	case http.StatusNotFound:
		return &NotFoundError{RefyneError: RefyneError{Message: msg, Status: 404}}
	case http.StatusTooManyRequests:
		retryAfter := 60
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if v, err := strconv.Atoi(ra); err == nil {
				retryAfter = v
			}
		}
		return &RateLimitError{RefyneError: RefyneError{Message: msg, Status: 429}, RetryAfter: retryAfter}
	default:
		return &RefyneError{Message: msg, Status: resp.StatusCode, Detail: errResp.Detail}
	}
}

func buildUserAgent(suffix string) string {
	ua := fmt.Sprintf("Refyne-SDK-Go/%s (Go/%s; %s/%s)", SDKVersion, runtime.Version()[2:], runtime.GOOS, runtime.GOARCH)
	if suffix != "" {
		ua += " " + suffix
	}
	return ua
}

func hashString(s string) string {
	var h uint32
	for _, c := range s {
		h = ((h << 5) - h) + uint32(c)
	}
	return strconv.FormatUint(uint64(h), 36)
}

// GenerateCacheKey generates a cache key from request details.
func GenerateCacheKey(method, urlStr, authHash string) string {
	parts := []string{strings.ToUpper(method), urlStr}
	if authHash != "" {
		parts = append(parts, authHash)
	}
	return strings.Join(parts, ":")
}
