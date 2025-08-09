# Redis Package

Redis utilities for distributed rate limiting in Go services.

## Features
Distributed rate limiting, sliding window, atomic operations (Lua), automatic cleanup, error handling, production-ready.

## Usage
```go
import (
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/redis"
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
)
redisLimiter := redis.NewRedisRateLimiter("localhost:6379", "", 0)
middleware := security.RateLimitMiddlewareWithLimiter(redisLimiter)
allowed, err := redisLimiter.Allow(ctx, "client-ip", 100, time.Minute)
```

## Config
addr: Redis address, password: Redis password, db: Redis DB number
