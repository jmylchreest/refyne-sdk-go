# Refyne Go SDK Examples

This directory contains example code demonstrating how to use the Refyne Go SDK.

## Prerequisites

- Go 1.21+
- A valid Refyne API key

## Environment Setup

Set the required environment variables:

```bash
export REFYNE_API_KEY="your_api_key_here"
export REFYNE_BASE_URL="https://api.refyne.uk"  # Optional, defaults to production
```

## Examples

### Full Demo (`full-demo/main.go`)

A comprehensive demo that tests all major SDK functionality:
- Usage/subscription information retrieval
- Job listing
- Schema listing
- Site listing
- Website analysis (structure detection)
- Single page extraction
- Crawl job creation and monitoring
- Job result retrieval

**Run with:**

```bash
cd examples/full-demo
go run main.go
```

Or from the SDK root:

```bash
go run ./examples/full-demo
```

### Basic Example (`basic/main.go`)

A simple example showing basic extraction:

```bash
cd examples/basic
go run main.go
```

## Notes

- The demo uses `fatih/color` package for colorful terminal output
- All API calls accept a `context.Context` for cancellation
- Error handling uses standard Go error patterns with typed errors
- The client supports functional options for configuration
