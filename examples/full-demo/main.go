// Full SDK Demo - Tests all major functionality
//
// Run with: go run examples/full-demo/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	refyne "github.com/jmylchreest/refyne-sdk-go"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	bgBlue       = "\033[44m"
	bgGreen      = "\033[42m"
)

// Spinner provides animated progress indication
type Spinner struct {
	frames  []string
	current int
	message string
	done    chan bool
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  []string{"|", "/", "-", "\\"},
		message: message,
		done:    make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	fmt.Print("\033[?25l") // Hide cursor
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				fmt.Printf("\r%s%s%s %s", colorCyan, s.frames[s.current], colorReset, s.message)
				s.current = (s.current + 1) % len(s.frames)
			}
		}
	}()
}

// Succeed stops the spinner with a success message
func (s *Spinner) Succeed(message string) {
	s.stop()
	if message == "" {
		message = s.message
	}
	fmt.Printf("\r\033[K%s[OK]%s %s\n", colorGreen, colorReset, message)
}

// Fail stops the spinner with a failure message
func (s *Spinner) Fail(message string) {
	s.stop()
	if message == "" {
		message = s.message
	}
	fmt.Printf("\r\033[K%s[FAIL]%s %s\n", colorRed, colorReset, message)
}

func (s *Spinner) stop() {
	close(s.done)
	fmt.Print("\033[?25h") // Show cursor
}

func header(text string) {
	fmt.Println()
	fmt.Printf("%s%s %s %s\n", bgBlue, colorBold, text, colorReset)
	fmt.Println()
}

func subheader(text string) {
	fmt.Printf("%s%s> %s%s\n", colorBold, colorBlue, text, colorReset)
}

func info(label, value string) {
	fmt.Printf("  %s%s:%s %s\n", colorDim, label, colorReset, value)
}

func success(text string) {
	fmt.Printf("%s[OK]%s %s\n", colorGreen, colorReset, text)
}

func warn(text string) {
	fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, text)
}

func errorMsg(text string) {
	fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, text)
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s%s%s\n", colorDim, string(data), colorReset)
}

func ptr[T any](v T) *T {
	return &v
}

func main() {
	// Configuration - Override with environment variables for local development
	apiKey := os.Getenv("REFYNE_API_KEY")
	if apiKey == "" {
		fmt.Println("REFYNE_API_KEY environment variable not set")
		os.Exit(1)
	}

	opts := []refyne.ClientOption{}
	if baseURL := os.Getenv("REFYNE_BASE_URL"); baseURL != "" {
		opts = append(opts, refyne.WithBaseURL(baseURL))
	}

	const testURL = "https://www.bbc.co.uk/news"

	ctx := context.Background()

	// Banner
	fmt.Println()
	fmt.Printf("%s%s+-----------------------------------------------------------+%s\n", colorBold, colorMagenta, colorReset)
	fmt.Printf("%s%s|%s           %sRefyne Go SDK - Full Demo%s                  %s%s|%s\n", colorBold, colorMagenta, colorReset, colorBold, colorReset, colorBold, colorMagenta, colorReset)
	fmt.Printf("%s%s+-----------------------------------------------------------+%s\n", colorBold, colorMagenta, colorReset)

	// ========== Configuration ==========
	header("Configuration")

	subheader("SDK Information")
	info("Version", refyne.SDKVersion)
	info("Runtime", fmt.Sprintf("Go %s", runtime.Version()))

	subheader("Client Settings")
	baseURL := refyne.DefaultBaseURL
	if envURL := os.Getenv("REFYNE_BASE_URL"); envURL != "" {
		baseURL = envURL
	}
	info("Base URL", baseURL)
	info("API Key", fmt.Sprintf("%s...%s", apiKey[:10], apiKey[len(apiKey)-4:]))

	// Create client
	client := refyne.NewClient(apiKey, opts...)

	// ========== Subscription Info ==========
	header("Usage Information")

	spinner := NewSpinner("Fetching usage details...")
	spinner.Start()

	usage, err := client.GetUsage(ctx)
	if err != nil {
		spinner.Fail("Failed to fetch usage")
		errorMsg(err.Error())
		os.Exit(1)
	}
	spinner.Succeed("Usage details retrieved")

	info("Total Jobs", fmt.Sprintf("%d", usage.TotalJobs))
	info("Total Charged", fmt.Sprintf("$%.2f USD", usage.TotalChargedUsd))
	info("BYOK Jobs", fmt.Sprintf("%d", usage.ByokJobs))

	// ========== Analyze ==========
	header("Website Analysis")

	subheader("Target")
	info("URL", testURL)

	spinner = NewSpinner("Analyzing website structure...")
	spinner.Start()

	var suggestedSchema map[string]any
	analysis, err := client.Analyze(ctx, refyne.AnalyzeInput{URL: testURL})
	if err != nil {
		spinner.Fail("Analysis unavailable")
		warn(err.Error())
		// Use a fallback schema
		suggestedSchema = map[string]any{
			"headline": "string",
			"summary":  "string",
		}
		info("Using fallback schema", "")
		printJSON(suggestedSchema)
	} else {
		spinner.Succeed("Website analysis complete")
		// Parse suggested schema from YAML/JSON string
		if err := json.Unmarshal([]byte(analysis.SuggestedSchema), &suggestedSchema); err != nil {
			// Try parsing as simple schema
			suggestedSchema = map[string]any{
				"headline": "string",
				"summary":  "string",
			}
		}
		info("Page Type", analysis.PageType)
		info("Suggested Schema", "")
		printJSON(suggestedSchema)

		if analysis.FollowPatterns != nil && len(*analysis.FollowPatterns) > 0 {
			var patterns []string
			for _, p := range *analysis.FollowPatterns {
				patterns = append(patterns, p.Pattern)
			}
			info("Follow Patterns", strings.Join(patterns, ", "))
		}
	}

	// ========== Single Page Extract ==========
	header("Single Page Extraction")

	subheader("Request")
	info("URL", testURL)
	info("Schema", "Using suggested schema from analysis")

	spinner = NewSpinner("Extracting data from page...")
	spinner.Start()

	extractResult, err := client.Extract(ctx, refyne.ExtractInput{
		URL:    testURL,
		Schema: suggestedSchema,
	})
	if err != nil {
		spinner.Fail("Extraction failed")
		warn(err.Error())
	} else {
		spinner.Succeed("Extraction complete")

		subheader("Result")
		info("Fetched At", extractResult.FetchedAt)
		info("Tokens", fmt.Sprintf("%d in / %d out", extractResult.Usage.InputTokens, extractResult.Usage.OutputTokens))
		info("Cost", fmt.Sprintf("$%.6f", extractResult.Usage.CostUsd))
		info("Model", fmt.Sprintf("%s/%s", extractResult.Metadata.Provider, extractResult.Metadata.Model))
		info("Duration", fmt.Sprintf("%dms fetch + %dms extract", extractResult.Metadata.FetchDurationMs, extractResult.Metadata.ExtractDurationMs))

		subheader("Extracted Data")
		printJSON(extractResult.Data)
	}

	// ========== Crawl Job ==========
	header("Crawl Job")

	subheader("Request")
	info("URL", testURL)
	info("Max URLs", "5")
	info("Schema", "Using suggested schema from analysis")

	spinner = NewSpinner("Starting crawl job...")
	spinner.Start()

	crawlResult, err := client.Crawl(ctx, refyne.CrawlInput{
		URL:    testURL,
		Schema: suggestedSchema,
		Options: &refyne.CrawlOptions{
			MaxUrls:  ptr(int64(5)),
			MaxDepth: ptr(int64(1)),
		},
	})
	if err != nil {
		spinner.Fail("Failed to start crawl")
		warn(err.Error())

		// Demo complete without crawl
		fmt.Println()
		fmt.Printf("%s%s Demo Complete %s\n", bgGreen, colorBold, colorReset)
		fmt.Println()
		return
	}
	spinner.Succeed("Crawl job started")

	jobID := crawlResult.JobId
	info("Job ID", jobID)
	info("Status", crawlResult.Status)

	// ========== Poll for Results ==========
	header("Monitoring Job Progress")

	subheader("Polling for status updates...")

	var lastStatus string
	var pageCount int64
	pollInterval := 2 * time.Second

	for {
		job, err := client.Jobs.Get(ctx, jobID)
		if err != nil {
			errorMsg(fmt.Sprintf("Failed to get job: %v", err))
			break
		}

		status := job.Status
		if status != lastStatus {
			fmt.Printf("  %s->%s Status: %s%s%s\n", colorCyan, colorReset, colorBold, status, colorReset)
			lastStatus = status
		}

		if job.PageCount > pageCount {
			newPages := job.PageCount - pageCount
			for i := int64(0); i < newPages; i++ {
				fmt.Printf("  %s[OK]%s Page %d extracted\n", colorGreen, colorReset, pageCount+i+1)
			}
			pageCount = job.PageCount
		}

		if status == "completed" || status == "failed" {
			if status == "completed" {
				success(fmt.Sprintf("Crawl completed - %d pages processed", job.PageCount))
			} else {
				msg := "Unknown error"
				if job.ErrorMessage != nil {
					msg = *job.ErrorMessage
				}
				errorMsg(fmt.Sprintf("Crawl failed: %s", msg))
			}
			break
		}

		time.Sleep(pollInterval)
	}

	// ========== Fetch Job Results ==========
	header("Job Results")

	spinner = NewSpinner("Fetching job details and results...")
	spinner.Start()

	job, err := client.Jobs.Get(ctx, jobID)
	if err != nil {
		spinner.Fail("Failed to fetch job")
		errorMsg(err.Error())
		os.Exit(1)
	}
	spinner.Succeed("Job details retrieved")

	subheader("Job Details")
	info("ID", job.Id)
	info("Type", job.Type)
	info("Status", job.Status)
	info("URL", job.Url)
	info("Pages Processed", fmt.Sprintf("%d", job.PageCount))
	info("Tokens", fmt.Sprintf("%d in / %d out", job.TokenUsageInput, job.TokenUsageOutput))
	info("Cost", fmt.Sprintf("$%.4f USD", job.CostUsd))
	if job.StartedAt != nil {
		info("Started", *job.StartedAt)
	}
	if job.CompletedAt != nil {
		info("Completed", *job.CompletedAt)
	}

	// Get results (merged)
	spinner = NewSpinner("Fetching extraction results...")
	spinner.Start()

	results, err := client.Jobs.GetResults(ctx, jobID, &refyne.ResultsOptions{Merge: true})
	if err != nil {
		spinner.Fail("Failed to fetch results")
		errorMsg(err.Error())
		os.Exit(1)
	}
	spinner.Succeed("Results retrieved")

	subheader("Extracted Data (Merged)")
	if len(results) > 0 {
		var data interface{}
		if err := json.Unmarshal(results, &data); err == nil {
			printJSON(data)
		} else {
			fmt.Println(string(results))
		}
	} else {
		warn("No results available")
	}

	// ========== Done ==========
	fmt.Println()
	fmt.Printf("%s%s Demo Complete %s\n", bgGreen, colorBold, colorReset)
	fmt.Println()
}
