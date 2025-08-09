# Hobby Streamer

A small playground to try streaming ideas, asset management, and distributed systems.

## What’s here
Video upload, analysis, HLS/DASH transcoding, CDN playback, GraphQL API, Neo4j, Keycloak auth, Redis cache, Kafka events, outbox, retries, health checks.

Services: [`asset-manager`](backend/asset-manager/README.md), [`auth-service`](backend/auth-service/README.md), [`transcoder`](backend/transcoder/README.md), [`streaming-api`](backend/streaming-api/README.md), Lambdas ([`raw_video_uploaded`](backend/lambdas/cmd/raw_video_uploaded/README.md), [`generate_video_upload_url`](backend/lambdas/cmd/generate_video_upload_url/README.md), [`generate_image_upload_url`](backend/lambdas/cmd/generate_image_upload_url/README.md), [`delete_files`](backend/lambdas/cmd/delete_files/README.md)), Frontends ([`HobbyStreamerCMS`](frontend/HobbyStreamerCMS/README.md), [`HobbyStreamerUI`](frontend/HobbyStreamerUI/README.md)). Shared library: `backend/pkg`.

## Architecture
![Architecture Diagram](docs/arch.png)
Docs: [CDN](docs/cdn-proposal.md), [Kafka](docs/kafka-architecture.md), [Distributed Techniques](docs/distributed-techniques.md)

## Quick start
Requirements: Docker, Go 1.21+, FFmpeg, Node.js 22+, awscli-local.

```bash
./local/build.sh
```

Open: AKHQ `http://localhost:8086`, Kibana `http://localhost:5601`.

## Dev helpers
```bash
make backend-lint
make backend-test
make backend-generate
make backend-build
```
