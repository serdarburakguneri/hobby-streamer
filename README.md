# Hobby Streamer

Hobby Streamer is a lightweight content management system (CMS) and streaming platform designed for experimenting end-to-end video workflows. The project enables you to:

- **Upload and manage video assets** through a simple API.
- **Process and transcode videos** for adaptive streaming (HLS/DASH) using FFMPEG.
- **Organize content** into buckets for easy management.
- **Deliver video content** in a way that mimics modern streaming platforms.

The goal is to provide a hands-on, cost-free environment for learning, prototyping, or building streaming solutions—without relying on any paid cloud infrastructure. All services (asset manager, transcoder, storage) run on your local machine, and AWS services (DynamoDB, SQS) are emulated using [LocalStack](https://github.com/localstack/localstack).

## Tech Stack
- LocalStack (DynamoDB, SQS, S3. Lambda) – Local AWS service emulation
- Go – Backend code for all services
- FFMPEG – For the transcoder service

## 🗂️ Architecture

![Architecture Diagram](docs/hobby-streamer.drawio.svg)

## 📚 Service Documentation

- [Asset Manager Service](backend/asset-manager/README.md): REST API for managing assets, images, videos, and buckets.
- [Auth Service](backend/auth-service/README.md): JWT-based authentication service with Keycloak integration.
- [Transcoder Service](backend/transcoder/README.md): Background worker for video analysis and transcoding jobs.
- [Storage Service](backend/storage/cmd/generate_presigned_upload_url/README.md): Lambda for generating S3 presigned URLs for direct uploads.

## 📦 Shared Libraries

- [Auth Package](backend/pkg/auth/README.md): Shared authentication library with JWT validation and role-based authorization.

## TODO

- Centralized logging
- A search mechanism for the asset manager service

## 🧪 Local Testing

### Prerequisites
- [Docker](https://www.docker.com/products/docker-desktop/) installed
- [Go](https://go.dev/doc/install) installed
- [Python](https://www.python.org/downloads/) installed
- [FFmpeg](https://ffmpeg.org/download.html) installed (required for video transcoding)
- [pipx](https://pypa.github.io/pipx/installation/) installed (for installing Python applications)
- [awscli-local (awslocal)](https://github.com/localstack/awscli-local) installed:
  ```sh
  pipx install awscli-local
  pipx ensurepath
  ```

### Local Environment Setup

To set up your entire local AWS-like environment (LocalStack, S3 buckets, DynamoDB tables, SQS queues) and start the core services, simply run:

```sh
./build.sh
```

This script will:
- Start LocalStack (via Docker Compose) if not already running
- Wait for LocalStack to be ready
- Create the required S3 buckets: `raw-storage`, `transcoded-storage`, `thumbnails-storage`
- Create the required DynamoDB tables: `asset`, `bucket`
- Create the required SQS queue: `transcoder-jobs`
- Start the Auth Service (on port 8080)
- Start the Asset Manager service (on port 8082)
- Start the Transcoder service (connected to the local SQS queue)

Logs for these services are written to `auth-service.log`, `asset-manager.log` and `transcoder.log` in the project root.

### Service Ports
- **Auth Service**: http://localhost:8080
- **Asset Manager**: http://localhost:8082
- **Keycloak**: http://localhost:9090
- **LocalStack**: http://localhost:4566

### Health Checks
Test if your services are running correctly:

```bash
# Auth Service Health Check
curl -s http://localhost:8080/health

# Asset Manager Health Check
curl -s http://localhost:8082/health

# Keycloak Health Check
curl -s http://localhost:9090/health

# LocalStack Health Check
curl -s http://localhost:4566/health
```

Expected responses:
- **Auth Service**: `{"status":"ok"}`
- **Asset Manager**: `{"status":"ok","service":"asset-manager"}`
- **Keycloak**: Should return a response (may be HTML)
- **LocalStack**: Should return a response (may be XML)

