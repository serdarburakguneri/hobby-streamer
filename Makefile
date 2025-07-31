.PHONY: help \
	backend-lint backend-test backend-build backend-clean backend-generate backend-generate-graphql backend-generate-mocks backend-install-tools backend-fmt backend-vet backend-security backend-quality backend-pre-commit backend-restart \
	integration-lint integration-test integration-test-auth integration-test-asset-manager integration-test-streaming-api integration-test-all \
	neo4j-clean redis-flush

# Default target
help:
	@echo "\n== Backend Commands =="
	@echo "  backend-lint             - Run golangci-lint on all Go files in backend"
	@echo "  backend-test             - Run Go tests for all backend modules"
	@echo "  backend-build            - Build all backend services"
	@echo "  backend-clean            - Clean backend build artifacts"
	@echo "  backend-generate         - Generate code (GraphQL, etc.)"
	@echo "  backend-generate-graphql - Generate GraphQL code only"
	@echo "  backend-generate-mocks   - Generate mock files for testing"
	@echo "  backend-install-tools    - Install Go development tools"
	@echo "  backend-fmt              - Format backend Go code"
	@echo "  backend-vet              - Run go vet on backend packages"
	@echo "  backend-security         - Run gosec security checks"
	@echo "  backend-quality          - Run all quality checks (fmt, vet, lint, security, test)"
	@echo "  backend-pre-commit       - Run pre-commit checks (lint, test)"
	@echo "  backend-restart SERVICE=<service> - Build and restart a backend service container"
	@echo "\n== Integration Test Commands =="
	@echo "  integration-test         - Run all Karate integration tests"
	@echo "  integration-test-auth    - Run only auth service integration tests"
	@echo "  integration-test-asset-manager - Run only asset-manager integration tests"
	@echo "  integration-test-streaming-api - Run only streaming-api integration tests"
	@echo "  integration-test-all     - Run all backend and integration tests"
	@echo "\n== Database Commands =="
	@echo "  neo4j-clean              - Delete all assets and buckets from Neo4j database"
	@echo "  redis-flush              - Flush all data from Redis cache"
	@echo "\n== Other =="
	@echo "  help                     - Show this help message"

# Backend targets (from backend/Makefile, now namespaced)
backend-install-tools:
	@echo "Installing development tools..."
	cd backend && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	cd backend && go install github.com/99designs/gqlgen@latest
	cd backend && go install github.com/golang/mock/mockgen@latest
	cd backend && go install github.com/vektah/dataloaden@latest
	cd backend && go install github.com/securego/gosec/v2/cmd/gosec@latest

backend-lint:
	@echo "Running golangci-lint..."
	cd backend && for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Linting $$module_dir..."; \
		cd "$$module_dir" && $${HOME}/go/bin/golangci-lint run ./... && cd - > /dev/null; \
	done

backend-test:
	@echo "Running tests..."
	cd backend && for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Testing $$module_dir..."; \
		cd "$$module_dir" && go test ./... -v && cd - > /dev/null; \
	done

backend-build:
	@echo "Building all services..."
	cd backend/asset-manager && go build -o asset-manager cmd/main.go
	cd backend/auth-service && go build -o auth-service cmd/main.go
	cd backend/streaming-api && go build -o streaming-api cmd/main.go
	cd backend/transcoder && go build -o transcoder cmd/worker/main.go

backend-clean:
	@echo "Cleaning build artifacts..."
	rm -f backend/asset-manager/asset-manager
	rm -f backend/auth-service/auth-service
	rm -f backend/streaming-api/streaming-api
	rm -f backend/transcoder/transcoder
	cd backend && find . -name "*.test" -delete

backend-generate:
	@echo "Generating code..."
	cd backend/asset-manager && go generate ./...

backend-generate-graphql:
	@echo "Generating GraphQL code..."
	cd backend/asset-manager && $${HOME}/go/bin/gqlgen generate

backend-generate-mocks:
	@echo "Generating mocks..."
	cd backend && $${HOME}/go/bin/mockgen -source=pkg/auth/interface.go -destination=pkg/auth/mocks.go
	cd backend && $${HOME}/go/bin/mockgen -source=pkg/sqs/consumer.go -destination=pkg/sqs/mocks.go

backend-pre-commit: backend-lint backend-test

backend-fmt:
	@echo "Formatting code..."
	cd backend && for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		echo "Formatting $$module_dir..."; \
		cd "$$module_dir" && go fmt ./... && cd - > /dev/null; \
	done

backend-vet:
	@echo "Running go vet..."
	cd backend && for dir in $$(find . -name "go.mod" -type f); do \
		module_dir=$$(dirname "$$dir"); \
		if [ -d "$$module_dir" ]; then \
			echo "Vetting $$module_dir..."; \
			cd "$$module_dir" && go vet ./... && cd - > /dev/null; \
		fi \
	done

backend-security:
	@echo "Running security checks..."
	cd backend && $${HOME}/go/bin/gosec ./...

backend-quality: backend-fmt backend-vet backend-lint backend-security backend-test

backend-restart:
	@echo "Restarting $${SERVICE} with Docker Compose..."
	docker-compose build $${SERVICE}
	docker-compose up -d $${SERVICE}

# Integration test targets
integration-test:
	@echo "Running all Karate integration tests..."
	cd integration-tests && mvn test

integration-test-auth:
	@echo "Running Karate auth service integration tests..."
	cd integration-tests && mvn test -Dtest=KarateTestRunner#testAuthService

integration-test-asset-manager:
	@echo "Running Karate asset-manager integration tests..."
	cd integration-tests && mvn test -Dtest=KarateTestRunner#testAssetManager

integration-test-streaming-api:
	@echo "Running Karate streaming-api integration tests..."
	cd integration-tests && mvn test -Dtest=KarateTestRunner#testStreamingApi

integration-test-all: backend-test integration-test

# Database targets
neo4j-clean:
	@echo "Deleting all assets and buckets from Neo4j database..."
	docker exec -it $$(docker ps -q -f name=neo4j) cypher-shell -u neo4j -p password "MATCH (n) WHERE n:Bucket OR n:Asset DETACH DELETE n"

redis-flush:
	@echo "Flushing Redis cache..."
	docker exec -it $$(docker ps -q -f name=redis) redis-cli FLUSHALL 