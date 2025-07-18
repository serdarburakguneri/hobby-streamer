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
	@echo "  fmt            - Format code with go fmt"
	@echo "  vet            - Run go vet on all packages"
	@echo "  security       - Run security checks with gosec"
	@echo "  quality        - Run all quality checks (fmt, vet, lint, security, test)"
	@echo "  pre-commit     - Run pre-commit checks (lint, test)"

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
	@for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Linting $$module_dir..."; \
		cd "$$module_dir" && $(HOME)/go/bin/golangci-lint run ./... && cd - > /dev/null; \
	done

# Run tests
test:
	@echo "Running tests..."
	@for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Testing $$module_dir..."; \
		cd "$$module_dir" && go test ./... -v && cd - > /dev/null; \
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
	@for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Formatting $$module_dir..."; \
		cd "$$module_dir" && go fmt ./... && cd - > /dev/null; \
	done

# Vet code
vet:
	@echo "Running go vet..."
	@for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Vetting $$module_dir..."; \
		cd "$$module_dir" && go vet ./... && cd - > /dev/null; \
	done

# Check for security issues
security:
	@echo "Running security checks..."
	$(HOME)/go/bin/gosec ./...

# Run all quality checks
quality: fmt vet lint security test 