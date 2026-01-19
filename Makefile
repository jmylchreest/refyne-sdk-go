.PHONY: generate generate-prod test lint fmt help

# Default target
help:
	@echo "Refyne Go SDK"
	@echo ""
	@echo "Usage:"
	@echo "  make generate          Generate types from OpenAPI spec (local server)"
	@echo "  make generate-prod     Generate types from production API"
	@echo "  make test              Run tests"
	@echo "  make test-race         Run tests with race detection"
	@echo "  make lint              Run linter"
	@echo "  make fmt               Format code"
	@echo ""
	@echo "Environment Variables:"
	@echo "  OPENAPI_SPEC_URL       Override the spec URL"

# Generate types from OpenAPI 3.0.3 spec (default: local dev server)
# Uses 3.0.3 for oapi-codegen compatibility
generate:
	@echo "Generating types from OpenAPI 3.0.3 spec..."
	oapi-codegen -config oapi-codegen.yaml http://localhost:8080/openapi-3.0.json
	@echo "Done. Generated types_generated.go"

# Generate types from production API
generate-prod:
	@echo "Generating types from production API..."
	oapi-codegen -config oapi-codegen.yaml https://api.refyne.uk/openapi-3.0.json
	@echo "Done. Generated types_generated.go"

# Run tests
test:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .
