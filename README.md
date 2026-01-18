# Refyne SDK for Go

Official Go SDK for the [Refyne API](https://docs.refyne.uk) - LLM-powered web extraction that transforms unstructured websites into clean, typed data.

[![Go Reference](https://pkg.go.dev/badge/github.com/jmylchreest/refyne-sdk-go.svg)](https://pkg.go.dev/github.com/jmylchreest/refyne-sdk-go)
[![CI](https://github.com/jmylchreest/refyne-sdk-go/actions/workflows/test.yml/badge.svg)](https://github.com/jmylchreest/refyne-sdk-go/actions/workflows/test.yml)

## Features

- **Idiomatic Go**: Context support, functional options, interfaces for testing
- **Zero Dependencies**: Uses only the standard library
- **Smart Caching**: Respects `Cache-Control` headers automatically
- **Auto-Retry**: Handles rate limits and transient errors with exponential backoff
- **API Version Compatibility**: Warns about breaking changes
- **Go 1.21+**: Uses modern Go features

## Installation

```bash
go get github.com/jmylchreest/refyne-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jmylchreest/refyne-sdk-go"
)

func main() {
    client := refyne.NewClient("your-api-key")

    result, err := client.Extract(context.Background(), refyne.ExtractRequest{
        URL: "https://example.com/product/123",
        Schema: map[string]any{
            "name":  "string",
            "price": "number",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.Data)
}
```

## Configuration

Use functional options to configure the client:

```go
client := refyne.NewClient(
    apiKey,
    refyne.WithBaseURL("https://api.refyne.uk"),
    refyne.WithTimeout(60*time.Second),
    refyne.WithMaxRetries(3),
    refyne.WithLogger(myLogger),
    refyne.WithCache(myCache),
    refyne.WithCacheEnabled(true),
    refyne.WithHTTPClient(myHTTPClient),
    refyne.WithUserAgentSuffix("MyApp/1.0"),
)
```

## Crawl Jobs

Extract data from multiple pages:

```go
// Start a crawl job
job, err := client.Crawl(ctx, refyne.CrawlRequest{
    URL:    "https://example.com/products",
    Schema: map[string]any{"name": "string", "price": "number"},
    Options: &refyne.CrawlOptions{
        FollowSelector: "a.product-link",
        MaxPages:       20,
        Delay:          "1s",
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Job started: %s\n", job.JobID)

// Poll for completion
for {
    status, err := client.Jobs.Get(ctx, job.JobID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %s (%d pages)\n", status.Status, status.PageCount)

    if status.Status == refyne.JobStatusCompleted ||
       status.Status == refyne.JobStatusFailed {
        break
    }

    time.Sleep(2 * time.Second)
}

// Get results
results, err := client.Jobs.GetResults(ctx, job.JobID, false)
if err != nil {
    log.Fatal(err)
}

for _, item := range results.Results {
    fmt.Println(item)
}
```

## Custom Logger

Implement the `Logger` interface:

```go
type MyLogger struct{}

func (l *MyLogger) Debug(msg string, meta map[string]any) {
    log.Printf("[DEBUG] %s %v\n", msg, meta)
}

func (l *MyLogger) Info(msg string, meta map[string]any) {
    log.Printf("[INFO] %s %v\n", msg, meta)
}

func (l *MyLogger) Warn(msg string, meta map[string]any) {
    log.Printf("[WARN] %s %v\n", msg, meta)
}

func (l *MyLogger) Error(msg string, meta map[string]any) {
    log.Printf("[ERROR] %s %v\n", msg, meta)
}

client := refyne.NewClient(apiKey, refyne.WithLogger(&MyLogger{}))
```

## Custom Cache

Implement the `Cache` interface:

```go
type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) Get(key string) (*refyne.CacheEntry, bool) {
    // Fetch from Redis
}

func (c *RedisCache) Set(key string, entry *refyne.CacheEntry) {
    // Store in Redis with TTL
}

func (c *RedisCache) Delete(key string) {
    // Delete from Redis
}

client := refyne.NewClient(apiKey, refyne.WithCache(&RedisCache{...}))
```

## Error Handling

```go
result, err := client.Extract(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *refyne.RateLimitError:
        fmt.Printf("Rate limited. Retry after %d seconds\n", e.RetryAfter)
    case *refyne.ValidationError:
        fmt.Printf("Validation errors: %v\n", e.Errors)
    case *refyne.AuthenticationError:
        fmt.Println("Invalid API key")
    case *refyne.RefyneError:
        fmt.Printf("API error: %s (%d)\n", e.Message, e.Status)
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## API Reference

### Main Client

| Method | Description |
|--------|-------------|
| `client.Extract(ctx, req)` | Extract data from a single page |
| `client.Crawl(ctx, req)` | Start an async crawl job |
| `client.Analyze(ctx, req)` | Analyze a site and suggest schema |
| `client.GetUsage(ctx)` | Get usage statistics |

### Sub-Services

| Service | Methods |
|---------|---------|
| `client.Jobs` | `List()`, `Get()`, `GetResults()` |
| `client.Schemas` | `List()`, `Get()`, `Create()`, `Update()`, `Delete()` |
| `client.Sites` | `List()`, `Get()`, `Create()`, `Update()`, `Delete()` |
| `client.Keys` | `List()`, `Create()`, `Revoke()` |
| `client.LLM` | `ListProviders()`, `ListKeys()`, `UpsertKey()`, `GetChain()`, `SetChain()` |

## Documentation

- [API Reference](https://docs.refyne.uk/docs/api-reference)
- [Go SDK Guide](https://docs.refyne.uk/docs/sdks/go)
- [GoDoc](https://pkg.go.dev/github.com/jmylchreest/refyne-sdk-go)

## Development

```bash
# Run tests
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run linter
golangci-lint run
```

## License

MIT License - see [LICENSE](LICENSE) for details.
