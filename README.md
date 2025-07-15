# Hobby Streamer

A personal playground for experimenting with video streaming and content management. This project explores building a simple streaming platform with basic asset management capabilities.

## Overview

The project consists of several microservices working together to provide video streaming functionality:

- **Asset Manager**: GraphQL service for managing video assets and metadata
- **Auth Service**: Authentication and authorization using Keycloak
- **Transcoder**: Video processing and format conversion
- **Storage**: S3-based file storage with presigned URLs
- **Frontend**: React Native CMS for content management
- **Logging**: Centralized logging system with Fluentd, Elasticsearch, and Kibana

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

### Frontend
- React Native – Cross-platform mobile and web development
- Expo – Development platform for React Native applications

## Architecture

![Architecture Diagram](docs/arch.png)

## Service Documentation

### Backend Services
- [Asset Manager Service](backend/asset-manager/README.md): GraphQL API for managing assets and relationships
- [Auth Service](backend/auth-service/README.md): JWT-based authentication service with Keycloak integration
- [Transcoder Service](backend/transcoder/README.md): Background worker for video analysis and transcoding jobs

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

This script will:
- Start all infrastructure services (LocalStack, Neo4j, Keycloak, Elasticsearch, Kibana, Fluentd)
- Create required S3 buckets and SQS queues
- Build and start all backend services
- Install dependencies and start the Admin UI development server
- Set up centralized logging system
- Verify all services are running correctly

### Service Ports
- Auth Service: http://localhost:8080
- Asset Manager GraphQL: http://localhost:8082/query
- Neo4j Browser: http://localhost:7474
- Keycloak: http://localhost:9090
- LocalStack: http://localhost:4566
- CMS UI Web: http://localhost:8081
- Kibana (Logs): http://localhost:5601
- Elasticsearch: http://localhost:9200

### Keycloak Setup

Keycloak is used for authentication. Default credentials:

- Admin User: `admin` / `admin`
- Regular User: `user` / `user`

### Running the CMS UI

After running `./local/build.sh`, the CMS UI will be available at:

```
http://localhost:8081
```

Login with the Keycloak credentials above to start managing assets.

### GraphQL Playground

The Asset Manager GraphQL API can be explored at:

```
http://localhost:8082/graphql
```

### Neo4j Browser

The graph database relationships can be explored at:

```
http://localhost:7474
```

Default credentials: `neo4j` / `password`

## Development Workflow

### Restarting Services

```bash
# Rebuild and restart specific services
docker-compose up --build -d asset-manager
docker-compose up --build -d auth-service
docker-compose up --build -d transcoder

# Rebuild everything
./local/build.sh
```

### Viewing Container Logs

The project includes a centralized logging system that collects logs from all services and makes them searchable through Kibana.

#### Kibana UI (Recommended)
- Access logs at: http://localhost:5601
- Go to "Discover" for advanced log searching
- Use "Dashboard" for predefined views
- Search by service, log level, HTTP details, user context, etc.

#### Command Line
```bash
# View logs for specific services
./local/logs.sh auth
./local/logs.sh asset
./local/logs.sh transcoder

# View all logs
./local/logs.sh all

# Traditional Docker logs (still available)
docker-compose logs -f asset-manager
docker-compose logs -f auth-service
docker-compose logs -f transcoder
docker-compose logs -f
```

#### Log Search Examples
```bash
# Find all errors
level:error

# Find slow requests (>1 second)
duration_ms:>1000

# Find auth service logs
service_name:auth-service

# Find 5xx errors with error details
status_code:[500 TO 599] AND error:*
```

See [Logging Documentation](local/LOGGING.md) for detailed information about the logging system.

### Health Checks

```bash
# Service health checks
curl -s http://localhost:8080/health
curl -s http://localhost:8082/health
curl -s http://localhost:9090/health
curl -s http://localhost:4566/health

# GraphQL introspection
curl -s -X POST http://localhost:8082/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __schema { types { name } } }"}'
```

