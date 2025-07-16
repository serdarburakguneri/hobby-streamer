# Streaming API

REST API for streaming UI with Redis caching.

## Overview

The streaming-api service provides fast, cached access to bucket and asset data for streaming applications. It uses Redis for caching and fetches data from the asset-manager GraphQL service when needed.


## API Endpoints

### Buckets
- `GET /api/v1/buckets` - List all buckets
- `GET /api/v1/buckets/{key}` - Get bucket by key
- `GET /api/v1/buckets/{key}/assets` - Get assets in bucket

### Assets
- `GET /api/v1/assets` - List all assets
- `GET /api/v1/assets/{slug}` - Get asset by slug

### Health
- `GET /health` - Health check endpoint

## Environment Variables

- `PORT` - Server port (default: 8080)
- `REDIS_URL` - Redis connection URL
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (optional)
- `ASSET_MANAGER_URL` - Asset manager GraphQL URL (default: http://localhost:8081)

## Cache Strategy

- **Bucket Details**: 30 minutes TTL
- **Asset Details**: 30 minutes TTL
- **Bucket List**: 15 minutes TTL
- **Asset List**: 15 minutes TTL

## Development

```bash
go run cmd/main.go
```

## Docker

```bash
docker build -t streaming-api .
docker run -p 8080:8080 streaming-api
``` 