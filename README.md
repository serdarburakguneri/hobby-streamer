# Hobby Streamer

> A personal playground for video streaming ideas. Not production, just for fun and learning.

This project explores video upload, transcoding, streaming delivery, asset metadata, authentication, and developer-friendly observability.

## Features
Video upload, transcoding (HLS/DASH), stream delivery, asset metadata, Keycloak-based auth, developer logging, monitoring, Docker Compose setup, Redis caching, SQS async processing, circuit breakers, retries, dead letter queues, health checks, rate limiting, local-first development.

## Architecture

![Architecture Diagram](docs/arch.png)

Docs: [Transcoding Sequence](docs/video-upload-transcode-sequence.md), [CDN Proposal](docs/cdn-proposal.md), [Domain Events](docs/domain-events.md)

## Tech Stack
Backend: Go, GraphQL, Neo4j, Keycloak, FFmpeg, Redis, Docker Compose, LocalStack, Fluentd, Elasticsearch, Kibana, Nginx. Frontend: React Native (CMS, viewer UI).

## Services
Backend: [`asset-manager`](backend/asset-manager/README.md), [`auth-service`](backend/auth-service/README.md), [`transcoder`](backend/transcoder/README.md), [`streaming-api`](backend/streaming-api/README.md). Lambdas: [`generate_video_upload_url`](backend/lambdas/cmd/generate_video_upload_url/README.md), [`generate_image_upload_url`](backend/lambdas/cmd/generate_image_upload_url/README.md), [`delete_files`](backend/lambdas/cmd/delete_files/README.md). Frontend: [`HobbyStreamerCMS`](frontend/HobbyStreamerCMS/README.md), [`HobbyStreamerUI`](frontend/HobbyStreamerUI/README.md). Shared: see `backend/pkg` for common code.

## Getting Started

Requirements: Docker, Go 1.21+, FFmpeg, Python + pipx, Node.js 22+, `awscli-local` (`pipx install awscli-local && pipx ensurepath`).

Quick start:
```bash
./local/build.sh
```
Starts all services, UIs, dependencies, and logging pipeline.

Development:
```bash
cd backend
make install-tools
make lint && make test
./scripts/pre-commit.sh
make generate && make build
```

## Observability
See [Logging Setup](local/LOGGING.md) for Fluentd → Elasticsearch → Kibana with structured logs and correlation IDs.