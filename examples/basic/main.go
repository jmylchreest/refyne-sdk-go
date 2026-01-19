// Example: basic extraction
//
// This example demonstrates how to extract structured data from a web page.
//
// Usage:
//
//	export REFYNE_API_KEY=your_api_key_here
//	go run examples/basic/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmylchreest/refyne-sdk-go"
)

func main() {
	apiKey := os.Getenv("REFYNE_API_KEY")
	if apiKey == "" {
		log.Fatal("REFYNE_API_KEY environment variable not set")
	}

	// Create client with optional custom base URL
	opts := []refyne.ClientOption{}
	if baseURL := os.Getenv("REFYNE_BASE_URL"); baseURL != "" {
		opts = append(opts, refyne.WithBaseURL(baseURL))
	}

	client := refyne.NewClient(apiKey, opts...)
	ctx := context.Background()

	fmt.Println("Extracting product data...")
	fmt.Println()

	// Extract structured data from a page
	result, err := client.Extract(ctx, refyne.ExtractInput{
		URL: "https://example.com/product/123",
		Schema: map[string]any{
			"name":        map[string]any{"type": "string", "description": "Product name"},
			"price":       map[string]any{"type": "number", "description": "Price in USD"},
			"description": map[string]any{"type": "string", "description": "Product description"},
			"inStock":     map[string]any{"type": "boolean", "description": "Whether in stock"},
		},
	})
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	fmt.Println("Extracted data:")
	if data, ok := result.Data.(map[string]any); ok {
		for key, value := range data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Printf("\nURL: %s\n", result.Url)
	fmt.Printf("Fetched at: %s\n", result.FetchedAt)

	fmt.Println("\nUsage:")
	fmt.Printf("  Input tokens: %d\n", result.Usage.InputTokens)
	fmt.Printf("  Output tokens: %d\n", result.Usage.OutputTokens)
	fmt.Printf("  Cost: $%.4f\n", result.Usage.CostUsd)

	fmt.Println("\nPerformance:")
	fmt.Printf("  Fetch time: %dms\n", result.Metadata.FetchDurationMs)
	fmt.Printf("  Extract time: %dms\n", result.Metadata.ExtractDurationMs)
	fmt.Printf("  Model: %s/%s\n", result.Metadata.Provider, result.Metadata.Model)
}
