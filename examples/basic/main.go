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
	"time"

	"github.com/jmylchreest/refyne-sdk-go"
)

func main() {
	apiKey := os.Getenv("REFYNE_API_KEY")
	if apiKey == "" {
		log.Fatal("REFYNE_API_KEY environment variable not set")
	}

	// Create a client with options
	client := refyne.NewClient(
		apiKey,
		refyne.WithTimeout(60*time.Second),
		refyne.WithUserAgentSuffix("BasicExample/1.0"),
	)

	ctx := context.Background()

	fmt.Println("Extracting product data...")
	fmt.Println()

	// Extract structured data from a page
	result, err := client.Extract(ctx, refyne.ExtractRequest{
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
	for key, value := range result.Data {
		fmt.Printf("  %s: %v\n", key, value)
	}

	fmt.Printf("\nURL: %s\n", result.URL)
	fmt.Printf("Fetched at: %s\n", result.FetchedAt)

	if result.Usage != nil {
		fmt.Println("\nUsage:")
		fmt.Printf("  Input tokens: %d\n", result.Usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", result.Usage.OutputTokens)
		fmt.Printf("  Cost: $%.4f\n", result.Usage.CostUSD)
	}

	if result.Metadata != nil {
		fmt.Println("\nPerformance:")
		fmt.Printf("  Fetch time: %dms\n", result.Metadata.FetchDurationMs)
		fmt.Printf("  Extract time: %dms\n", result.Metadata.ExtractDurationMs)
		fmt.Printf("  Model: %s/%s\n", result.Metadata.Provider, result.Metadata.Model)
	}
}
