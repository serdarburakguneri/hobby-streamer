# Hobby Streamer

A personal playground for experimenting with video streaming and content management. This project explores building a simple streaming platform with basic asset management capabilities.

## Architecture

![Architecture Diagram](docs/arch.png)

### Video Upload and Transcoding Flow

For a detailed view of how video uploading and transcoding works in the system, see the [Video Upload and Transcoding Sequence Diagram](docs/video-upload-transcode-sequence.md). This diagram shows the complete flow from user upload through analysis and transcoding to HLS/DASH formats.

## Tech Stack

### Backend Services
- Go – Backend services
- GraphQL – API layer for asset management
- Neo4j – Graph database for asset relationships and metadata
- Keycloak – Identity and access management
- FFMPEG – Video processing and transcoding

### Infrastructure & DevOps
- Docker Compose – Containerized development environment
- LocalStack – Local AWS service emulation (S3, SQS, Lambda)
- Fluentd – Log collection and forwarding
- Elasticsearch – Log storage and search
- Kibana – Log visualization and analysis
- Nginx (To replace cloudfront locally)

### Frontend
- React Native – Cross-platform mobile and web development

## Service Documentation

### Backend Services
- [Asset Manager Service](backend/asset-manager/README.md): GraphQL API for managing assets and relationships
- [Auth Service](backend/auth-service/README.md): JWT-based authentication service with Keycloak integration
- [Transcoder Service](backend/transcoder/README.md): Background worker for video analysis and transcoding jobs
- [Streaming API Service](backend/streaming-api/README.md): REST API with Redis caching for streaming applications

### Lambdas
- [Generate Presigned Upload URL Lambda](backend/lambdas/cmd/generate_presigned_upload_url/README.md): Lambda for generating S3 presigned URLs for direct uploads
- [Delete Files Lambda](backend/lambdas/cmd/delete_files/README.md): Lambda for cleaning up S3 files when assets are deleted
- [Trigger Transcode Job Lambda](backend/lambdas/cmd/trigger_transcode_job/README.md): Lambda for triggering video transcoding jobs

### Frontend Services
- [CMS UI](frontend/HobbyStreamerCMS/README.md): React Native CMS interface for managing assets

### Infrastructure & Logging
- [Logging System](local/LOGGING.md): Centralized logging with Fluentd, Elasticsearch, and Kibana

### Shared Libraries
See [Shared Libraries Documentation](backend/pkg/README.md) for detailed information about the shared library architecture and available packages.

**Available Libraries:**
- [Auth Package](backend/pkg/auth/README.md): Shared authentication library with JWT validation and role-based authorization
- [Constants Package](backend/pkg/constants/README.md): Common constants for HTTP status codes, roles, and other shared values
- [Logger Package](backend/pkg/logger/README.md): Centralized structured logging solution for all backend services
- [Messages Package](backend/pkg/messages/README.md): Common SQS message payload structures and type constants for inter-service communication
- [S3 Package](backend/pkg/s3/README.md): S3 client library for file upload, download, and directory operations with LocalStack support
- [SQS Package](backend/pkg/sqs/README.md): AWS SQS client library with producer, consumer, and consumer registry functionality

## Getting Started

### Prerequisites
- Docker installed
- Go (version 1.21+) installed
- FFmpeg installed (required for video transcoding)
- Python installed (For localstack)
- pipx installed (for installing Python applications)
- awscli-local (awslocal) installed:
  ```sh
  pipx install awscli-local
  pipx ensurepath
  ```
- Node.js (version 22+) installed for frontend development

### Quick Start

To set up the development environment with all services, simply run:

```sh
./local/build.sh
```

### Service Ports
- Auth Service: http://localhost:8080
- Asset Manager GraphQL: http://localhost:8082/query
- Streaming API: http://localhost:8084
- Redis: redis://localhost:6379
- Neo4j Browser: http://localhost:7474
- Keycloak: http://localhost:9090
- LocalStack: http://localhost:4566
- CMS UI Web: http://localhost:8081
- Kibana (Logs): http://localhost:5601
- Elasticsearch: http://localhost:9200
- Nginx HLS Proxy: http://localhost:8083

