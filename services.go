package refyne

import (
	"context"
	"fmt"
	"net/http"
)

// JobsService handles job-related operations.
type JobsService struct {
	client *Client
}

// List returns all jobs.
func (s *JobsService) List(ctx context.Context, limit, offset int) (*JobList, error) {
	path := "/api/v1/jobs"
	if limit > 0 || offset > 0 {
		path += "?"
		if limit > 0 {
			path += fmt.Sprintf("limit=%d", limit)
		}
		if offset > 0 {
			if limit > 0 {
				path += "&"
			}
			path += fmt.Sprintf("offset=%d", offset)
		}
	}

	var resp JobList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a job by ID.
func (s *JobsService) Get(ctx context.Context, id string) (*Job, error) {
	var resp Job
	if err := s.client.requestWithOptions(ctx, http.MethodGet, "/api/v1/jobs/"+id, nil, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetResults returns job results.
func (s *JobsService) GetResults(ctx context.Context, id string, merge bool) (*JobResults, error) {
	path := "/api/v1/jobs/" + id + "/results"
	if merge {
		path += "?merge=true"
	}

	var resp JobResults
	if err := s.client.requestWithOptions(ctx, http.MethodGet, path, nil, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SchemasService handles schema operations.
type SchemasService struct {
	client *Client
}

// List returns all schemas.
func (s *SchemasService) List(ctx context.Context) (*SchemaList, error) {
	var resp SchemaList
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/schemas", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a schema by ID.
func (s *SchemasService) Get(ctx context.Context, id string) (*Schema, error) {
	var resp Schema
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/schemas/"+id, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new schema.
func (s *SchemasService) Create(ctx context.Context, req CreateSchemaRequest) (*Schema, error) {
	var resp Schema
	if err := s.client.request(ctx, http.MethodPost, "/api/v1/schemas", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a schema.
func (s *SchemasService) Update(ctx context.Context, id string, req CreateSchemaRequest) (*Schema, error) {
	var resp Schema
	if err := s.client.request(ctx, http.MethodPut, "/api/v1/schemas/"+id, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a schema.
func (s *SchemasService) Delete(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/schemas/"+id, nil, nil)
}

// SitesService handles site operations.
type SitesService struct {
	client *Client
}

// List returns all sites.
func (s *SitesService) List(ctx context.Context) (*SiteList, error) {
	var resp SiteList
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/sites", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a site by ID.
func (s *SitesService) Get(ctx context.Context, id string) (*Site, error) {
	var resp Site
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/sites/"+id, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new site.
func (s *SitesService) Create(ctx context.Context, req CreateSiteRequest) (*Site, error) {
	var resp Site
	if err := s.client.request(ctx, http.MethodPost, "/api/v1/sites", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a site.
func (s *SitesService) Update(ctx context.Context, id string, req CreateSiteRequest) (*Site, error) {
	var resp Site
	if err := s.client.request(ctx, http.MethodPut, "/api/v1/sites/"+id, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a site.
func (s *SitesService) Delete(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/sites/"+id, nil, nil)
}

// KeysService handles API key operations.
type KeysService struct {
	client *Client
}

// List returns all API keys.
func (s *KeysService) List(ctx context.Context) (*APIKeyList, error) {
	var resp APIKeyList
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/keys", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new API key.
func (s *KeysService) Create(ctx context.Context, name string) (*APIKeyCreated, error) {
	var resp APIKeyCreated
	if err := s.client.request(ctx, http.MethodPost, "/api/v1/keys", map[string]string{"name": name}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Revoke revokes an API key.
func (s *KeysService) Revoke(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/keys/"+id, nil, nil)
}

// LLMService handles LLM configuration.
type LLMService struct {
	client *Client
}

// ListProviders returns available providers.
func (s *LLMService) ListProviders(ctx context.Context) (*ProvidersResponse, error) {
	var resp ProvidersResponse
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/llm/providers", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListKeys returns configured provider keys.
func (s *LLMService) ListKeys(ctx context.Context) (*LLMKeyList, error) {
	var resp LLMKeyList
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/llm/keys", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpsertKey adds or updates a provider key.
func (s *LLMService) UpsertKey(ctx context.Context, req UpsertLLMKeyRequest) (*LLMKey, error) {
	var resp LLMKey
	if err := s.client.request(ctx, http.MethodPut, "/api/v1/llm/keys", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteKey deletes a provider key.
func (s *LLMService) DeleteKey(ctx context.Context, id string) error {
	return s.client.request(ctx, http.MethodDelete, "/api/v1/llm/keys/"+id, nil, nil)
}

// GetChain returns the fallback chain configuration.
func (s *LLMService) GetChain(ctx context.Context) (*LLMChain, error) {
	var resp LLMChain
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/llm/chain", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetChain sets the fallback chain configuration.
func (s *LLMService) SetChain(ctx context.Context, chain []LLMChainEntry) error {
	return s.client.request(ctx, http.MethodPut, "/api/v1/llm/chain", map[string]any{"chain": chain}, nil)
}

// ListModels returns available models for a provider.
func (s *LLMService) ListModels(ctx context.Context, provider string) (*ModelList, error) {
	var resp ModelList
	if err := s.client.request(ctx, http.MethodGet, "/api/v1/llm/models/"+provider, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
