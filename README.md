# Hobby Streamer

> Just a personal sandbox — not production, not polished.

A side project experimenting with video streaming platform ideas. Tinkering with Go, streaming pipelines, infrastructure tools, and architectural patterns.

## What This Is

Video upload/transcoding, stream delivery (HLS/DASH), asset metadata, auth with Keycloak, and developer-focused logging/monitoring.

## Architecture

![Architecture Diagram](docs/arch.png) | [Sequence Diagram](docs/video-upload-transcode-sequence.md)

## Tech Stack

**Backend:** Go, GraphQL, Neo4j, Keycloak, FFmpeg, Redis  
**Infra:** Docker Compose, LocalStack (mock AWS), Fluentd + Elasticsearch + Kibana, Nginx  
**Frontend:** React Native (CMS + viewer UI)

## Services

### Backend
- [`asset-manager`](backend/asset-manager/README.md) - Asset metadata & relationships
- [`auth-service`](backend/auth-service/README.md) - OAuth2 + RBAC
- [`transcoder`](backend/transcoder/README.md) - FFmpeg-based video processing
- [`streaming-api`](backend/streaming-api/README.md) - Stream delivery API

### Lambdas
- [`generate_video_upload_url`](backend/lambdas/cmd/generate_video_upload_url/README.md) - Video upload presigned URLs
- [`generate_image_upload_url`](backend/lambdas/cmd/generate_image_upload_url/README.md) - Image upload presigned URLs
- [`delete_files`](backend/lambdas/cmd/delete_files/README.md) - Cleanup uploaded files

### Frontend
- [`HobbyStreamerCMS`](frontend/HobbyStreamerCMS/README.md) - Content management
- [`HobbyStreamerUI`](frontend/HobbyStreamerUI/README.md) - Video viewer

### Shared Libs
See `backend/pkg` for auth, config, constants, error handling, logging, etc.

## Features

SQS-driven async processing, circuit breakers/retries, dead letter queues, structured logging, health checks, rate limiting, Redis caching, dev-friendly local setup.

## Getting Started

### Requirements
- Docker, Go 1.21+, FFmpeg, Python + pipx, Node.js 22+
- `awscli-local`: `pipx install awscli-local && pipx ensurepath`

### Quick Start
```bash
./local/build.sh
```
Starts all services, UIs, dependencies (Redis, Neo4j, Keycloak, LocalStack), and logging pipeline.

### Development
```bash
cd backend
make install-tools
make lint && make test
./scripts/pre-commit.sh
make generate && make build
```

## Observability
[Logging Setup](local/LOGGING.md) - Fluentd → Elasticsearch → Kibana with structured logs and correlation IDs.