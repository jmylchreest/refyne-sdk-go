package refyne

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey 'test-api-key', got '%s'", client.apiKey)
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
	if client.timeout != DefaultTimeout {
		t.Errorf("expected timeout %v, got %v", DefaultTimeout, client.timeout)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("expected maxRetries %d, got %d", DefaultMaxRetries, client.maxRetries)
	}

	// Check sub-clients are initialized
	if client.Jobs == nil {
		t.Error("Jobs sub-client not initialized")
	}
	if client.Schemas == nil {
		t.Error("Schemas sub-client not initialized")
	}
	if client.Sites == nil {
		t.Error("Sites sub-client not initialized")
	}
	if client.Keys == nil {
		t.Error("Keys sub-client not initialized")
	}
	if client.LLM == nil {
		t.Error("LLM sub-client not initialized")
	}
}

func TestClientOptions(t *testing.T) {
	customURL := "https://custom.api.test"
	customTimeout := 60 * time.Second
	customRetries := 5

	client := NewClient("test-api-key",
		WithBaseURL(customURL+"/"),
		WithTimeout(customTimeout),
		WithMaxRetries(customRetries),
	)

	if client.baseURL != customURL {
		t.Errorf("expected baseURL '%s', got '%s'", customURL, client.baseURL)
	}
	if client.timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, client.timeout)
	}
	if client.maxRetries != customRetries {
		t.Errorf("expected maxRetries %d, got %d", customRetries, client.maxRetries)
	}
}

func TestAuthenticationHeader(t *testing.T) {
	apiKey := "test-bearer-token"
	var capturedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_jobs":        0,
			"total_charged_usd": 0,
			"byok_jobs":         0,
		})
	}))
	defer server.Close()

	client := NewClient(apiKey, WithBaseURL(server.URL))
	_, err := client.GetUsage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Bearer " + apiKey
	if capturedAuth != expected {
		t.Errorf("expected Authorization '%s', got '%s'", expected, capturedAuth)
	}
}

func TestUserAgentHeader(t *testing.T) {
	var capturedUA string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_jobs":        0,
			"total_charged_usd": 0,
			"byok_jobs":         0,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	_, err := client.GetUsage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "refyne-go/" + SDKVersion
	if capturedUA != expected {
		t.Errorf("expected User-Agent '%s', got '%s'", expected, capturedUA)
	}
}

func TestExtract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/extract" {
			t.Errorf("expected path '/api/v1/extract', got '%s'", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["url"] != "https://example.com" {
			t.Errorf("expected url 'https://example.com', got '%v'", body["url"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"title": "Test"},
			"url":        "https://example.com",
			"fetched_at": "2024-01-01T00:00:00Z",
			"usage": map[string]any{
				"input_tokens":  100,
				"output_tokens": 50,
				"cost_usd":      0.001,
			},
			"metadata": map[string]any{
				"provider":            "test",
				"model":               "test-model",
				"fetch_duration_ms":   100,
				"extract_duration_ms": 200,
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	result, err := client.Extract(context.Background(), ExtractInput{
		URL:    "https://example.com",
		Schema: map[string]any{"title": "string"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Url != "https://example.com" {
		t.Errorf("expected url 'https://example.com', got '%s'", result.Url)
	}
}

func TestCrawl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/crawl" {
			t.Errorf("expected path '/api/v1/crawl', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"job_id": "job-123",
			"status": "pending",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	result, err := client.Crawl(context.Background(), CrawlInput{
		URL:    "https://example.com",
		Schema: map[string]any{"title": "string"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.JobId != "job-123" {
		t.Errorf("expected job_id 'job-123', got '%s'", result.JobId)
	}
}

func TestJobsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/jobs" {
			t.Errorf("expected path '/api/v1/jobs', got '%s'", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got '%s'", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"jobs":  []any{},
			"total": 0,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	_, err := client.Jobs.List(context.Background(), &ListOptions{Limit: 10})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJobsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/jobs/job-123" {
			t.Errorf("expected path '/api/v1/jobs/job-123', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":                 "job-123",
			"type":               "crawl",
			"status":             "completed",
			"url":                "https://example.com",
			"page_count":         5,
			"token_usage_input":  1000,
			"token_usage_output": 500,
			"cost_usd":           0.01,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	result, err := client.Jobs.Get(context.Background(), "job-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "job-123" {
		t.Errorf("expected id 'job-123', got '%s'", result.Id)
	}
}

func TestError400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":  "validation failed",
			"errors": map[string]string{"url": "required"},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	_, err := client.Extract(context.Background(), ExtractInput{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if valErr.Fields["url"] != "required" {
		t.Errorf("expected field error 'required', got '%s'", valErr.Fields["url"])
	}
}

func TestError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid token"})
	}))
	defer server.Close()

	client := NewClient("bad-key", WithBaseURL(server.URL))
	_, err := client.GetUsage(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	authErr, ok := err.(*AuthError)
	if !ok {
		t.Fatalf("expected AuthError, got %T", err)
	}
	if !strings.Contains(authErr.Error(), "invalid token") {
		t.Errorf("expected error message to contain 'invalid token', got '%s'", authErr.Error())
	}
}

func TestError403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{"error": "forbidden"})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	_, err := client.GetUsage(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, ok := err.(*ForbiddenError)
	if !ok {
		t.Fatalf("expected ForbiddenError, got %T", err)
	}
}

func TestError404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "not found"})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	_, err := client.Jobs.Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, ok := err.(*NotFoundError)
	if !ok {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
}

func TestError429RateLimitWithRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.Header().Set("Retry-After", "0")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"error": "rate limited"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_jobs":        0,
			"total_charged_usd": 0,
			"byok_jobs":         0,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL), WithTimeout(5*time.Second))
	_, err := client.GetUsage(context.Background())

	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestError500WithRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": "internal error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_jobs":        0,
			"total_charged_usd": 0,
			"byok_jobs":         0,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL), WithTimeout(5*time.Second))
	_, err := client.GetUsage(context.Background())

	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestSchemasCRUD(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/schemas":
			json.NewEncoder(w).Encode(map[string]any{"schemas": []any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/schemas":
			json.NewEncoder(w).Encode(map[string]any{
				"id":          "schema-1",
				"name":        "Test",
				"schema_yaml": "type: object",
			})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/schemas/"):
			json.NewEncoder(w).Encode(map[string]any{
				"id":          "schema-1",
				"name":        "Test",
				"schema_yaml": "type: object",
			})
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/v1/schemas/"):
			json.NewEncoder(w).Encode(map[string]any{
				"id":          "schema-1",
				"name":        "Updated",
				"schema_yaml": "type: object",
			})
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v1/schemas/"):
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	// List
	_, err := client.Schemas.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	// Create
	created, err := client.Schemas.Create(ctx, CreateSchemaInput{
		Name:       "Test",
		SchemaYAML: "type: object",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.Id != "schema-1" {
		t.Errorf("expected id 'schema-1', got '%s'", created.Id)
	}

	// Get
	_, err = client.Schemas.Get(ctx, "schema-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Update
	updated, err := client.Schemas.Update(ctx, "schema-1", CreateSchemaInput{
		Name:       "Updated",
		SchemaYAML: "type: object",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name 'Updated', got '%s'", updated.Name)
	}

	// Delete
	err = client.Schemas.Delete(ctx, "schema-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestSitesCRUD(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sites":
			json.NewEncoder(w).Encode(map[string]any{"sites": []any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sites":
			json.NewEncoder(w).Encode(map[string]any{
				"id":   "site-1",
				"name": "Test Site",
				"url":  "https://example.com",
			})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/sites/"):
			json.NewEncoder(w).Encode(map[string]any{
				"id":   "site-1",
				"name": "Test Site",
				"url":  "https://example.com",
			})
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/v1/sites/"):
			json.NewEncoder(w).Encode(map[string]any{
				"id":   "site-1",
				"name": "Updated Site",
				"url":  "https://example.com",
			})
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v1/sites/"):
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	// List
	_, err := client.Sites.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	// Create
	created, err := client.Sites.Create(ctx, CreateSiteInput{
		Name: "Test Site",
		URL:  "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.Id != "site-1" {
		t.Errorf("expected id 'site-1', got '%s'", created.Id)
	}

	// Get
	_, err = client.Sites.Get(ctx, "site-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Update
	updated, err := client.Sites.Update(ctx, "site-1", CreateSiteInput{
		Name: "Updated Site",
		URL:  "https://example.com",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name == nil || *updated.Name != "Updated Site" {
		t.Errorf("expected name 'Updated Site', got '%v'", updated.Name)
	}

	// Delete
	err = client.Sites.Delete(ctx, "site-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestLLMOperations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.URL.Path == "/api/v1/llm/providers":
			json.NewEncoder(w).Encode(map[string]any{
				"providers": []any{
					map[string]any{"id": "anthropic", "name": "Anthropic"},
					map[string]any{"id": "openai", "name": "OpenAI"},
				},
			})
		case strings.HasPrefix(r.URL.Path, "/api/v1/llm/models/"):
			json.NewEncoder(w).Encode(map[string]any{
				"models": []any{
					map[string]any{"id": "model-1", "name": "Model 1"},
				},
			})
		case r.URL.Path == "/api/v1/llm/keys" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{"keys": []any{}})
		case r.URL.Path == "/api/v1/llm/keys" && r.Method == http.MethodPut:
			json.NewEncoder(w).Encode(map[string]any{
				"id":       "key-1",
				"provider": "anthropic",
			})
		case r.URL.Path == "/api/v1/llm/chain" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{"chain": []any{}})
		case r.URL.Path == "/api/v1/llm/chain" && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	// List providers
	providers, err := client.LLM.ListProviders(ctx)
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}
	if providers.Providers == nil || len(*providers.Providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(*providers.Providers))
	}

	// List models
	_, err = client.LLM.ListModels(ctx, "anthropic")
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	// List keys
	_, err = client.LLM.ListKeys(ctx)
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}

	// Upsert key
	_, err = client.LLM.UpsertKey(ctx, UpsertKeyInput{
		Provider:     "anthropic",
		APIKey:       "sk-test",
		DefaultModel: "claude-3-5-sonnet",
	})
	if err != nil {
		t.Fatalf("UpsertKey failed: %v", err)
	}

	// Get chain
	_, err = client.LLM.GetChain(ctx)
	if err != nil {
		t.Fatalf("GetChain failed: %v", err)
	}

	// Set chain
	err = client.LLM.SetChain(ctx, []ChainEntry{
		{Provider: "anthropic", Model: "claude-3-5-sonnet", IsEnabled: true},
	})
	if err != nil {
		t.Fatalf("SetChain failed: %v", err)
	}
}

func TestCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client := NewClient("test-key", WithHTTPClient(customClient))

	if client.httpClient != customClient {
		t.Error("custom HTTP client not set")
	}
}

func TestBackoffCalculation(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
		{5, 16 * time.Second},
		{6, 30 * time.Second}, // Max capped at 30s
		{7, 30 * time.Second},
	}

	for _, tt := range tests {
		got := client.calculateBackoff(tt.attempt)
		if got != tt.expected {
			t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.expected)
		}
	}
}

func TestRetryAfterParsing(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		header   string
		expected time.Duration
	}{
		{"", 1 * time.Second},
		{"5", 5 * time.Second},
		{"0", 0},
		{"invalid", 1 * time.Second},
	}

	for _, tt := range tests {
		got := client.parseRetryAfter(tt.header)
		if got != tt.expected {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", tt.header, got, tt.expected)
		}
	}
}
