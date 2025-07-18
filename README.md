# Hobby Streamer

A personal playground for building a lightweight, end-to-end video streaming platform. This project demonstrates practical experience in streaming architecture, asset management, and distributed system design — all while honing Go expertise and infrastructure fluency.

## What It Does

Hobby Streamer is a modular streaming system designed to handle:
- Uploading and transcoding user-generated videos
- Secure and scalable video delivery (HLS/DASH)
- Asset relationship and metadata management
- Authenticated access and role-based permissioning
- Developer-friendly logging, monitoring, and local emulation

## Architecture Overview

![Architecture Diagram](docs/arch.png)

For a closer look at the media pipeline, see the [Upload & Transcode Sequence Diagram](docs/video-upload-transcode-sequence.md).

## Tech Stack

### Backend
- Go – All backend services use idiomatic Go with error boundaries and resilience patterns
- GraphQL – Asset management API
- Neo4j – Graph-based modeling for video assets and relationships
- Keycloak – OAuth2-based identity and access control
- FFmpeg – Transcoding and media probing
- Redis – Lightweight caching for stream-related queries

### Infrastructure
- Docker Compose – Local orchestration
- LocalStack – Emulated AWS (S3, SQS, Lambda)
- Fluentd + Elasticsearch + Kibana – Full log pipeline
- Nginx – Local replacement for CloudFront CDN

### Frontend
- React Native – Streaming and CMS frontend (Web support enabled)


## Code Organization

### Core Backend Services

- [Asset Manager](backend/asset-manager/README.md): GraphQL API for asset CRUD and relationships
- [Auth Service](backend/auth-service/README.md): JWT auth with Keycloak
- [Transcoder](backend/transcoder/README.md): Worker for FFmpeg-based transcoding
- [Streaming API](backend/streaming-api/README.md): REST API with Redis caching

### Lambdas

- [Generate Presigned Upload URL](backend/lambdas/cmd/generate_presigned_upload_url/README.md): Generates temporary S3 upload URLs
- [Delete Files on Asset Deletion](backend/lambdas/cmd/delete_files/README.md): Cleans up S3 when assets are deleted
- [Trigger Transcode Job](backend/lambdas/cmd/trigger_transcode_job/README.md): Triggers video processing on upload

### Frontend

- [CMS UI](frontend/HobbyStreamerCMS/README.md): React Native dashboard for asset management
- [Streaming UI](frontend/HobbyStreamerUI/README.md): React Native client for watching videos

### Shared Libraries

See [Shared Libraries Documentation](backend/pkg/README.md) for details.

- [auth](backend/pkg/auth): JWT validation and RBAC support
- [config](backend/pkg/config): Dynamic configuration with feature flags
- [constants](backend/pkg/constants): Shared enums and constants
- [errors](backend/pkg/errors): Typed error handling, retries, and circuit breakers
- [logger](backend/pkg/logger): Centralized structured logging
- [messages](backend/pkg/messages): Shared message structures for SQS
- [s3](backend/pkg/s3): File management utilities for S3/LocalStack
- [security](backend/pkg/security): Rate limiting, CORS protection, and security headers
- [sqs](backend/pkg/sqs): Client utilities for producing and consuming SQS events

### Observability

See [Logging Setup](local/LOGGING.md) for details on how Fluentd, Elasticsearch, and Kibana are integrated into the stack.

## Features

### Asynchronous Processing & Resilience
- **Event-driven architecture** with SQS for reliable message processing
- **Circuit breakers** and retry mechanisms for external service calls
- **Graceful degradation** with fallback strategies
- **Dead letter queues** for failed message handling
- **Distributed tracing** through structured logging

### Observability & Monitoring
- **Structured logging** with correlation IDs across all services
- **Centralized log aggregation** with Fluentd → Elasticsearch → Kibana
- **Health checks** and readiness probes for all services
- **Error tracking** with detailed context and stack traces

### Security & Access Control
- **OAuth2/JWT authentication** with Keycloak integration
- **Role-based access control** (RBAC) with fine-grained permissions
- **Rate limiting** and DDoS protection
- **Input validation** and sanitization
- **Secure file uploads** with presigned URLs
- **CORS protection** and security headers

### Scalability & Performance
- **Horizontal scaling** ready with stateless service design
- **Redis caching** for frequently accessed data
- **Efficient video transcoding** with parallel processing

### Developer Experience
- **Local development** with full AWS emulation
- **Hot reloading** for rapid iteration
- **Code quality checks** with linting, formatting, and security scanning

## Getting Started

### Prerequisites

- Docker installed
- Go (version 1.21+) installed
- FFmpeg installed (required for video transcoding)
- Python installed (for LocalStack)
- pipx installed (for installing Python applications)
- awscli-local installed:
  ```
  pipx install awscli-local
  pipx ensurepath
  ```
- Node.js (version 22+) installed for frontend development

### Quick Start

To spin up all services and dependencies:

```
./local/build.sh
```

This will:
- Start backend services
- Launch the CMS and streaming UIs
- Set up Redis, Neo4j, Keycloak, and LocalStack
- Configure the logging pipeline

### Development Workflow

For development and code quality:

```bash
# Navigate to backend directory
cd backend

# Install development tools
make install-tools

# Run quality checks
make lint
make test

# Run pre-commit checks (alternatively, add a git hook)
./scripts/pre-commit.sh

# Generate code
make generate

# Build all services
make build
```
