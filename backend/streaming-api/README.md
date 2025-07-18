# Streaming API

A lightweight REST API designed for fast, read-optimized access to bucket and asset data. Uses Redis for caching and includes fallback logic for fetching from the Asset Manager when needed.

---

## Overview

This service acts as a frontend-friendly layer over the asset data, optimized for speed and reliability. It helps reduce load on the main GraphQL service by caching frequent queries and adding resilience via retries and circuit breakers.

---

## Features

- REST endpoints for fetching buckets and assets
- Redis caching with configurable TTL
- Circuit breaker support for upstream failures
- Retry logic with exponential backoff
- Graceful error handling for degraded scenarios

---

## API Endpoints

### Buckets

- `GET /api/v1/buckets` — List all buckets  
- `GET /api/v1/buckets/{key}` — Fetch a single bucket by key  
- `GET /api/v1/buckets/{key}/assets` — List assets in a given bucket

### Assets

- `GET /api/v1/assets` — List all assets  
- `GET /api/v1/assets/{slug}` — Get a single asset by slug

### Health

- `GET /health` — Basic health check

---

## Caching Strategy

- Bucket list: **15 minutes**  
- Bucket detail: **30 minutes**  
- Asset list: **15 minutes**  
- Asset detail: **30 minutes**

All cache entries are keyed by request parameters and automatically expire based on the TTLs above.

---

## Running Locally

```bash
cd backend/streaming-api
go run cmd/main.go
```

---

## Docker

```bash
docker build -t streaming-api .
docker run -p 8080:8080 streaming-api
```

---

> ⚠️ This API is meant for local development and frontend testing. Behavior and cache logic may evolve as the project matures.