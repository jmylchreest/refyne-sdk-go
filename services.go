package refyne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// JobsClient handles job-related operations.
type JobsClient struct {
	client *Client
}

// ListOptions contains options for listing jobs.
type ListOptions struct {
	Limit  int
	Offset int
}

// List returns all jobs.
func (j *JobsClient) List(ctx context.Context, opts *ListOptions) (*ListJobsOutputBody, error) {
	path := "/api/v1/jobs"
	if opts != nil {
		params := ""
		if opts.Limit > 0 {
			params += fmt.Sprintf("limit=%d", opts.Limit)
		}
		if opts.Offset > 0 {
			if params != "" {
				params += "&"
			}
			params += fmt.Sprintf("offset=%d", opts.Offset)
		}
		if params != "" {
			path += "?" + params
		}
	}

	var result ListJobsOutputBody
	if err := j.client.request(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a job by ID.
func (j *JobsClient) Get(ctx context.Context, id string) (*JobResponse, error) {
	var result JobResponse
	if err := j.client.request(ctx, http.MethodGet, "/api/v1/jobs/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ResultsOptions contains options for getting job results.
type ResultsOptions struct {
	Merge bool
}

// GetResults returns job results.
func (j *JobsClient) GetResults(ctx context.Context, id string, opts *ResultsOptions) (json.RawMessage, error) {
	path := "/api/v1/jobs/" + id + "/results"
	if opts != nil && opts.Merge {
		path += "?merge=true"
	}

	var result json.RawMessage
	if err := j.client.request(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SchemasClient handles schema operations.
type SchemasClient struct {
	client *Client
}

// List returns all schemas.
func (s *SchemasClient) List(ctx context.Context) (*ListSchemasOutputBody, error) {
	var result ListSchemasOutputBody
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/schemas", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a schema by ID.
func (s *SchemasClient) Get(ctx context.Context, id string) (*SchemaOutput, error) {
	var result SchemaOutput
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/schemas/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateInput contains parameters for creating a schema.
type CreateSchemaInput struct {
	Name       string `json:"name"`
	SchemaYAML string `json:"schema_yaml"`
	Visibility string `json:"visibility,omitempty"`
}

// Create creates a new schema.
func (s *SchemasClient) Create(ctx context.Context, input CreateSchemaInput) (*SchemaOutput, error) {
	var result SchemaOutput
	if err := s.client.request(ctx, http.MethodPost, "/api/v1/schemas", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a schema.
func (s *SchemasClient) Update(ctx context.Context, id string, input CreateSchemaInput) (*SchemaOutput, error) {
	var result SchemaOutput
	if err := s.client.request(ctx, http.MethodPut, "/api/v1/schemas/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a schema.
func (s *SchemasClient) Delete(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/schemas/"+id, nil, nil)
}

// SitesClient handles site operations.
type SitesClient struct {
	client *Client
}

// List returns all sites.
func (s *SitesClient) List(ctx context.Context) (*ListSavedSitesOutputBody, error) {
	var result ListSavedSitesOutputBody
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/sites", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a site by ID.
func (s *SitesClient) Get(ctx context.Context, id string) (*SavedSiteOutput, error) {
	var result SavedSiteOutput
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/sites/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSiteInput contains parameters for creating a site.
type CreateSiteInput struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	DefaultSchemaID string `json:"default_schema_id,omitempty"`
	FetchMode       string `json:"fetch_mode,omitempty"`
}

// Create creates a new site.
func (s *SitesClient) Create(ctx context.Context, input CreateSiteInput) (*SavedSiteOutput, error) {
	var result SavedSiteOutput
	if err := s.client.request(ctx, http.MethodPost, "/api/v1/sites", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a site.
func (s *SitesClient) Update(ctx context.Context, id string, input CreateSiteInput) (*SavedSiteOutput, error) {
	var result SavedSiteOutput
	if err := s.client.request(ctx, http.MethodPut, "/api/v1/sites/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a site.
func (s *SitesClient) Delete(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/sites/"+id, nil, nil)
}

// KeysClient handles API key operations.
type KeysClient struct {
	client *Client
}

// List returns all API keys.
func (k *KeysClient) List(ctx context.Context) (*ListKeysOutputBody, error) {
	var result ListKeysOutputBody
	if err := k.client.request(ctx, http.MethodGet, "/api/v1/keys", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new API key.
func (k *KeysClient) Create(ctx context.Context, name string) (*CreateKeyOutputBody, error) {
	var result CreateKeyOutputBody
	if err := k.client.request(ctx, http.MethodPost, "/api/v1/keys", map[string]string{"name": name}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Revoke revokes an API key.
func (k *KeysClient) Revoke(ctx context.Context, id string) error {
	return k.client.request(ctx, http.MethodDelete, "/api/v1/keys/"+id, nil, nil)
}

// LLMClient handles LLM configuration.
type LLMClient struct {
	client *Client
}

// ListProviders returns available LLM providers.
func (l *LLMClient) ListProviders(ctx context.Context) (*ListProvidersOutputBody, error) {
	var result ListProvidersOutputBody
	if err := l.client.request(ctx, http.MethodGet, "/api/v1/llm/providers", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListModels returns available models for a provider.
func (l *LLMClient) ListModels(ctx context.Context, provider string) (*UserListModelsOutputBody, error) {
	var result UserListModelsOutputBody
	if err := l.client.request(ctx, http.MethodGet, "/api/v1/llm/models/"+provider, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListKeys returns configured LLM provider keys.
func (l *LLMClient) ListKeys(ctx context.Context) (*ListUserServiceKeysOutputBody, error) {
	var result ListUserServiceKeysOutputBody
	if err := l.client.request(ctx, http.MethodGet, "/api/v1/llm/keys", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpsertKeyInput contains parameters for upserting an LLM key.
type UpsertKeyInput struct {
	Provider     string `json:"provider"`
	APIKey       string `json:"api_key"`
	DefaultModel string `json:"default_model"`
	BaseURL      string `json:"base_url,omitempty"`
}

// UpsertKey adds or updates an LLM provider key.
func (l *LLMClient) UpsertKey(ctx context.Context, input UpsertKeyInput) (*UserServiceKeyResponse, error) {
	var result UserServiceKeyResponse
	if err := l.client.request(ctx, http.MethodPut, "/api/v1/llm/keys", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteKey deletes an LLM provider key.
func (l *LLMClient) DeleteKey(ctx context.Context, id string) error {
	return l.client.request(ctx, http.MethodDelete, "/api/v1/llm/keys/"+id, nil, nil)
}

// GetChain returns the LLM fallback chain configuration.
func (l *LLMClient) GetChain(ctx context.Context) (*GetUserFallbackChainOutputBody, error) {
	var result GetUserFallbackChainOutputBody
	if err := l.client.request(ctx, http.MethodGet, "/api/v1/llm/chain", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChainEntry represents an entry in the fallback chain.
type ChainEntry struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	IsEnabled bool   `json:"is_enabled"`
}

// SetChain sets the LLM fallback chain configuration.
func (l *LLMClient) SetChain(ctx context.Context, entries []ChainEntry) error {
	return l.client.request(ctx, http.MethodPut, "/api/v1/llm/chain", map[string]any{"chain": entries}, nil)
}
