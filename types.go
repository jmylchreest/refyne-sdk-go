package refyne

// ExtractRequest represents a data extraction request.
type ExtractRequest struct {
	// URL is the URL to extract data from.
	URL string `json:"url"`
	// Schema defines the data structure to extract.
	Schema map[string]any `json:"schema"`
	// FetchMode is the fetch mode: auto, static, or dynamic.
	FetchMode string `json:"fetchMode,omitempty"`
	// LLMConfig is the custom LLM configuration.
	LLMConfig *LLMConfig `json:"llmConfig,omitempty"`
}

// ExtractResponse represents the response from data extraction.
type ExtractResponse struct {
	// Data is the extracted data matching the schema.
	Data map[string]any `json:"data"`
	// URL is the URL that was extracted.
	URL string `json:"url"`
	// FetchedAt is when the page was fetched.
	FetchedAt string `json:"fetchedAt"`
	// Usage contains token usage information.
	Usage *TokenUsage `json:"usage,omitempty"`
	// Metadata contains extraction metadata.
	Metadata *ExtractionMetadata `json:"metadata,omitempty"`
}

// TokenUsage represents token usage from extraction.
type TokenUsage struct {
	// InputTokens is the number of input tokens used.
	InputTokens int `json:"inputTokens"`
	// OutputTokens is the number of output tokens used.
	OutputTokens int `json:"outputTokens"`
	// CostUSD is the total USD cost charged.
	CostUSD float64 `json:"costUsd"`
	// LLMCostUSD is the actual LLM cost from provider.
	LLMCostUSD float64 `json:"llmCostUsd"`
	// IsBYOK indicates if user's own API key was used.
	IsBYOK bool `json:"isByok"`
}

// ExtractionMetadata contains metadata from extraction.
type ExtractionMetadata struct {
	// FetchDurationMs is the time to fetch the page in milliseconds.
	FetchDurationMs int `json:"fetchDurationMs"`
	// ExtractDurationMs is the time to extract data in milliseconds.
	ExtractDurationMs int `json:"extractDurationMs"`
	// Model is the model used for extraction.
	Model string `json:"model"`
	// Provider is the LLM provider used.
	Provider string `json:"provider"`
}

// CrawlRequest represents a crawl job request.
type CrawlRequest struct {
	// URL is the seed URL to start crawling from.
	URL string `json:"url"`
	// Schema defines the data structure to extract.
	Schema map[string]any `json:"schema"`
	// Options contains crawl options.
	Options *CrawlOptions `json:"options,omitempty"`
	// WebhookURL is the URL to notify on completion.
	WebhookURL string `json:"webhookUrl,omitempty"`
	// LLMConfig is the custom LLM configuration.
	LLMConfig *LLMConfig `json:"llmConfig,omitempty"`
}

// CrawlOptions contains options for crawl jobs.
type CrawlOptions struct {
	// FollowSelector is the CSS selector for links to follow.
	FollowSelector string `json:"followSelector,omitempty"`
	// FollowPattern is the regex pattern for URLs to follow.
	FollowPattern string `json:"followPattern,omitempty"`
	// MaxDepth is the maximum crawl depth.
	MaxDepth int `json:"maxDepth,omitempty"`
	// NextSelector is the CSS selector for pagination.
	NextSelector string `json:"nextSelector,omitempty"`
	// MaxPages is the maximum pages to crawl.
	MaxPages int `json:"maxPages,omitempty"`
	// MaxURLs is the maximum total URLs to process.
	MaxURLs int `json:"maxUrls,omitempty"`
	// Delay is the delay between requests (e.g., "500ms").
	Delay string `json:"delay,omitempty"`
	// Concurrency is the number of concurrent requests.
	Concurrency int `json:"concurrency,omitempty"`
	// SameDomainOnly restricts to same-domain links.
	SameDomainOnly *bool `json:"sameDomainOnly,omitempty"`
	// ExtractFromSeeds enables extraction from seed URLs.
	ExtractFromSeeds *bool `json:"extractFromSeeds,omitempty"`
}

// CrawlJobCreated is returned when a crawl job is created.
type CrawlJobCreated struct {
	// JobID is the unique job identifier.
	JobID string `json:"jobId"`
	// Status is the initial job status.
	Status JobStatus `json:"status"`
	// StatusURL is the URL to check job status.
	StatusURL string `json:"statusUrl"`
}

// JobStatus represents the status of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// Job represents job details.
type Job struct {
	// ID is the job ID.
	ID string `json:"id"`
	// Type is the job type.
	Type string `json:"type"`
	// Status is the current status.
	Status JobStatus `json:"status"`
	// URL is the seed URL.
	URL string `json:"url"`
	// PageCount is the number of pages processed.
	PageCount int `json:"pageCount"`
	// TokenUsageInput is input tokens used.
	TokenUsageInput int `json:"tokenUsageInput"`
	// TokenUsageOutput is output tokens used.
	TokenUsageOutput int `json:"tokenUsageOutput"`
	// CostCredits is the cost in credits.
	CostCredits float64 `json:"costCredits"`
	// ErrorMessage is the error message if failed.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// StartedAt is when the job started.
	StartedAt string `json:"startedAt,omitempty"`
	// CompletedAt is when the job completed.
	CompletedAt string `json:"completedAt,omitempty"`
	// CreatedAt is when the job was created.
	CreatedAt string `json:"createdAt"`
}

// JobList is a list of jobs.
type JobList struct {
	Jobs []Job `json:"jobs"`
}

// JobResults contains job results.
type JobResults struct {
	// JobID is the job ID.
	JobID string `json:"jobId"`
	// Status is the job status.
	Status JobStatus `json:"status"`
	// PageCount is the number of pages processed.
	PageCount int `json:"pageCount"`
	// Results is the array of extraction results.
	Results []map[string]any `json:"results,omitempty"`
	// Merged is the merged results object.
	Merged map[string]any `json:"merged,omitempty"`
}

// AnalyzeRequest represents a website analysis request.
type AnalyzeRequest struct {
	// URL is the URL to analyze.
	URL string `json:"url"`
	// Depth is the analysis depth.
	Depth int `json:"depth,omitempty"`
}

// AnalyzeResponse is the response from website analysis.
type AnalyzeResponse struct {
	// URL is the URL that was analyzed.
	URL string `json:"url"`
	// SuggestedSchema is the suggested schema.
	SuggestedSchema map[string]any `json:"suggestedSchema"`
	// FollowPatterns are suggested URL patterns.
	FollowPatterns []string `json:"followPatterns"`
}

// Schema represents a schema definition.
type Schema struct {
	// ID is the schema ID.
	ID string `json:"id"`
	// Name is the schema name.
	Name string `json:"name"`
	// Description is the schema description.
	Description string `json:"description,omitempty"`
	// SchemaYAML is the schema definition in YAML.
	SchemaYAML string `json:"schemaYaml"`
	// Category is the category.
	Category string `json:"category,omitempty"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt"`
	// UpdatedAt is the last update timestamp.
	UpdatedAt string `json:"updatedAt"`
}

// SchemaList is a list of schemas.
type SchemaList struct {
	Schemas []Schema `json:"schemas"`
}

// CreateSchemaRequest is used to create a schema.
type CreateSchemaRequest struct {
	// Name is the schema name.
	Name string `json:"name"`
	// SchemaYAML is the schema definition in YAML.
	SchemaYAML string `json:"schemaYaml"`
	// Description is the schema description.
	Description string `json:"description,omitempty"`
	// Category is the category.
	Category string `json:"category,omitempty"`
}

// Site represents a saved site.
type Site struct {
	// ID is the site ID.
	ID string `json:"id"`
	// Name is the site name.
	Name string `json:"name"`
	// URL is the site URL.
	URL string `json:"url"`
	// SchemaID is the associated schema ID.
	SchemaID string `json:"schemaId,omitempty"`
	// CrawlOptions are the default crawl options.
	CrawlOptions *CrawlOptions `json:"crawlOptions,omitempty"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt"`
}

// SiteList is a list of sites.
type SiteList struct {
	Sites []Site `json:"sites"`
}

// CreateSiteRequest is used to create a site.
type CreateSiteRequest struct {
	// Name is the site name.
	Name string `json:"name"`
	// URL is the site URL.
	URL string `json:"url"`
	// SchemaID is the associated schema ID.
	SchemaID string `json:"schemaId,omitempty"`
	// CrawlOptions are the default crawl options.
	CrawlOptions *CrawlOptions `json:"crawlOptions,omitempty"`
}

// APIKey represents an API key (without the secret).
type APIKey struct {
	// ID is the key ID.
	ID string `json:"id"`
	// Name is the key name.
	Name string `json:"name"`
	// Prefix is the key prefix.
	Prefix string `json:"prefix"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt"`
	// LastUsedAt is the last used timestamp.
	LastUsedAt string `json:"lastUsedAt,omitempty"`
}

// APIKeyList is a list of API keys.
type APIKeyList struct {
	Keys []APIKey `json:"keys"`
}

// APIKeyCreated is returned when a key is created.
type APIKeyCreated struct {
	// ID is the key ID.
	ID string `json:"id"`
	// Name is the key name.
	Name string `json:"name"`
	// Key is the full API key (only shown once).
	Key string `json:"key"`
}

// UsageResponse contains usage statistics.
type UsageResponse struct {
	// Tier is the user's tier.
	Tier string `json:"tier"`
	// CreditsUsed is credits used this period.
	CreditsUsed float64 `json:"creditsUsed"`
	// CreditsLimit is the credit limit.
	CreditsLimit float64 `json:"creditsLimit"`
	// CreditsRemaining is credits remaining.
	CreditsRemaining float64 `json:"creditsRemaining"`
	// PeriodStart is the period start date.
	PeriodStart string `json:"periodStart"`
	// PeriodEnd is the period end date.
	PeriodEnd string `json:"periodEnd"`
}

// LLMConfig contains LLM provider configuration.
type LLMConfig struct {
	// Provider is the LLM provider.
	Provider string `json:"provider,omitempty"`
	// APIKey is the API key for the provider.
	APIKey string `json:"apiKey,omitempty"`
	// BaseURL is the custom base URL.
	BaseURL string `json:"baseUrl,omitempty"`
	// Model is the model to use.
	Model string `json:"model,omitempty"`
}

// LLMKey represents an LLM provider key.
type LLMKey struct {
	// ID is the key ID.
	ID string `json:"id"`
	// Provider is the provider name.
	Provider string `json:"provider"`
	// DefaultModel is the default model.
	DefaultModel string `json:"defaultModel"`
	// BaseURL is the custom base URL.
	BaseURL string `json:"baseUrl,omitempty"`
	// IsEnabled indicates if the key is enabled.
	IsEnabled bool `json:"isEnabled"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt"`
}

// LLMKeyList is a list of LLM keys.
type LLMKeyList struct {
	Keys []LLMKey `json:"keys"`
}

// UpsertLLMKeyRequest is used to upsert an LLM key.
type UpsertLLMKeyRequest struct {
	// Provider is the provider name.
	Provider string `json:"provider"`
	// APIKey is the API key.
	APIKey string `json:"apiKey"`
	// DefaultModel is the default model.
	DefaultModel string `json:"defaultModel"`
	// BaseURL is the custom base URL.
	BaseURL string `json:"baseUrl,omitempty"`
	// IsEnabled indicates if the key is enabled.
	IsEnabled *bool `json:"isEnabled,omitempty"`
}

// LLMChainEntry is an entry in the fallback chain.
type LLMChainEntry struct {
	// ID is the entry ID.
	ID string `json:"id,omitempty"`
	// Position is the position in chain.
	Position int `json:"position,omitempty"`
	// Provider is the provider name.
	Provider string `json:"provider"`
	// Model is the model name.
	Model string `json:"model"`
	// IsEnabled indicates if enabled.
	IsEnabled *bool `json:"isEnabled,omitempty"`
}

// LLMChain is the fallback chain.
type LLMChain struct {
	Chain []LLMChainEntry `json:"chain"`
}

// Model represents an available model.
type Model struct {
	// ID is the model ID.
	ID string `json:"id"`
	// Name is the display name.
	Name string `json:"name"`
}

// ModelList is a list of models.
type ModelList struct {
	Models []Model `json:"models"`
}

// ProvidersResponse contains available providers.
type ProvidersResponse struct {
	Providers []string `json:"providers"`
}
