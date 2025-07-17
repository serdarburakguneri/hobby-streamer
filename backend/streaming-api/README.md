# Streaming API

A REST API for streaming applications. Provides fast access to bucket and asset data with Redis caching and fallback mechanisms. Fetches data from the asset manager service when not available in cache.

## Overview

The service acts as a read-only interface optimized for frontend consumption. It reduces load on the asset manager by caching common queries in Redis. It includes basic resilience patterns like retries, circuit breakers, and graceful error handling.

## Features

- REST API for retrieving buckets and assets
- Redis caching with TTL configuration
- Circuit breakers for upstream calls
- Retry logic with exponential backoff

## API Endpoints

### Buckets

- `GET /api/v1/buckets` – List all buckets
- `GET /api/v1/buckets/{key}` – Get a single bucket by key
- `GET /api/v1/buckets/{key}/assets` – List assets in a specific bucket

### Assets

- `GET /api/v1/assets` – List all assets
- `GET /api/v1/assets/{slug}` – Get asset by slug

### Health

- `GET /health` – Returns service status


## Cache Strategy

- Bucket list: 15 minutes
- Bucket details: 30 minutes
- Asset list: 15 minutes
- Asset details: 30 minutes

## Development

```bash
cd backend/streaming-api
go run cmd/main.go
```

## Docker

```bash
docker build -t streaming-api .
docker run -p 8080:8080 streaming-api
```