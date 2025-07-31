#!/bin/bash
set -e

echo "Running pre-commit checks..."

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "golangci-lint not found. Installing..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

# Run formatting check
echo "Checking code formatting..."
if ! cd backend && go fmt ./...; then
    echo "Code formatting issues found. Please run 'cd backend && go fmt ./...' and commit again."
    exit 1
fi

# Run linting
echo "Running linting..."
if ! cd backend && $HOME/go/bin/golangci-lint run ./...; then
    echo "Linting issues found. Please fix them and commit again."
    exit 1
fi

# Run tests
echo "Running tests..."
if ! cd backend && go test ./...; then
    echo "Tests failed. Please fix them and commit again."
    exit 1
fi

# Run security checks if gosec is available
if command -v $HOME/go/bin/gosec &> /dev/null; then
    echo "Running security checks..."
    if ! cd backend && $HOME/go/bin/gosec ./...; then
        echo "Security issues found. Please review and fix them."
        exit 1
    fi
fi

echo "All pre-commit checks passed!" 