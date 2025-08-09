# Streaming API

REST API for fast, read-optimized access to bucket and asset data, with Redis caching and fallback to Asset Manager.

## Features
REST endpoints, Redis caching, circuit breakers, retry logic, graceful error handling.

## API
Buckets: `GET /api/v1/buckets`, `GET /api/v1/buckets/{key}`, `GET /api/v1/buckets/{key}/assets`. Assets: `GET /api/v1/assets`, `GET /api/v1/assets/{slug}`. Health: `GET /health`.

## Caching
Bucket and asset list/detail: 15/30 minutes.

## Running
```bash
cd backend/streaming-api && go run cmd/main.go
```
