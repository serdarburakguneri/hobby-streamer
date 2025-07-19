# Streaming API

Lightweight REST API for fast, read-optimized access to bucket and asset data. Uses Redis caching with fallback to Asset Manager.

## Features

REST endpoints, Redis caching, circuit breakers, retry logic, graceful error handling.

## API

### Buckets
- `GET /api/v1/buckets` — List all buckets  
- `GET /api/v1/buckets/{key}` — Fetch bucket by key  
- `GET /api/v1/buckets/{key}/assets` — List assets in bucket

### Assets
- `GET /api/v1/assets` — List all assets  
- `GET /api/v1/assets/{slug}` — Get asset by slug

### Health
- `GET /health` — Health check

## Caching

- Bucket list/detail: **15/30 minutes**  
- Asset list/detail: **15/30 minutes**

## Running

```bash
# Local
cd backend/streaming-api && go run cmd/main.go

# Docker
docker build -t streaming-api . && docker run -p 8080:8080 streaming-api
```

> ⚠️ Local development and frontend testing. Cache logic may evolve.