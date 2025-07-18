.PHONY: help lint test build clean generate

# Default target
help:
	@echo "Available targets:"
	@echo "  lint           - Run golangci-lint on all Go files"
	@echo "  test           - Run tests for all packages"
	@echo "  build          - Build all services"
	@echo "  clean          - Clean build artifacts"
	@echo "  generate       - Generate code (GraphQL, etc.)"
	@echo "  generate-graphql - Generate GraphQL code only"
	@echo "  generate-mocks - Generate mock files for testing"
	@echo "  install-tools  - Install development tools"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/99designs/gqlgen@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/vektah/dataloaden@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run linting
lint:
	@echo "Running golangci-lint..."
	@for service in asset-manager auth-service streaming-api transcoder; do \
		if [ -d "$$service" ]; then \
			echo "Linting $$service..."; \
			cd $$service && $(HOME)/go/bin/golangci-lint run ./... && cd ..; \
		fi; \
	done
	@for pkg in pkg/*; do \
		if [ -d "$$pkg" ]; then \
			echo "Linting $$pkg..."; \
			cd $$pkg && $(HOME)/go/bin/golangci-lint run ./... && cd ../..; \
		fi; \
	done

# Run tests
test:
	@echo "Running tests..."
	@for service in asset-manager auth-service streaming-api transcoder; do \
		if [ -d "$$service" ]; then \
			echo "Testing $$service..."; \
			cd $$service && go test ./... -v && cd ..; \
		fi; \
	done
	@for pkg in pkg/*; do \
		if [ -d "$$pkg" ]; then \
			echo "Testing $$pkg..."; \
			cd $$pkg && go test ./... -v && cd ../..; \
		fi; \
	done

# Build all services
build:
	@echo "Building all services..."
	cd asset-manager && go build -o asset-manager cmd/main.go
	cd auth-service && go build -o auth-service cmd/main.go
	cd streaming-api && go build -o streaming-api cmd/main.go
	cd transcoder && go build -o transcoder cmd/worker/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f asset-manager/asset-manager
	rm -f auth-service/auth-service
	rm -f streaming-api/streaming-api
	rm -f transcoder/transcoder
	find . -name "*.test" -delete

# Generate code
generate:
	@echo "Generating code..."
	cd asset-manager && go generate ./...

# Generate GraphQL code
generate-graphql:
	@echo "Generating GraphQL code..."
	cd asset-manager && $(HOME)/go/bin/gqlgen generate

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	$(HOME)/go/bin/mockgen -source=pkg/auth/interface.go -destination=pkg/auth/mocks.go
	$(HOME)/go/bin/mockgen -source=pkg/sqs/consumer.go -destination=pkg/sqs/mocks.go

# Pre-commit checks
pre-commit: lint test

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Running go vet..."
	go vet ./...

# Check for security issues
security:
	@echo "Running security checks..."
	$(HOME)/go/bin/gosec ./...

# Run all quality checks
quality: fmt vet lint security test 