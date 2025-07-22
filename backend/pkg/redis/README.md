# Redis Package

This package provides Redis-based utilities for distributed environments.

## Rate Limiter

The `RedisRateLimiter` provides distributed rate limiting using Redis, suitable for production environments with multiple service instances.

### Usage

```go
import (
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/redis"
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
)

// Create Redis rate limiter
redisLimiter := redis.NewRedisRateLimiter("localhost:6379", "", 0)

// Use with security middleware
middleware := security.RateLimitMiddlewareWithLimiter(redisLimiter)

// Or use directly
allowed, err := redisLimiter.Allow(ctx, "client-ip", 100, time.Minute)
```

### Features

- **Distributed**: Works across multiple service instances
- **Sliding Window**: Uses Redis sorted sets for accurate sliding window rate limiting
- **Atomic Operations**: Uses Lua scripts for atomic rate limit checks
- **Automatic Cleanup**: Automatically expires rate limit keys
- **Error Handling**: Graceful fallback when Redis is unavailable

### Configuration

- `addr`: Redis server address (e.g., "localhost:6379")
- `password`: Redis password (empty string for no auth)
- `db`: Redis database number (0-15)

### Production Considerations

- Use Redis cluster for high availability
- Configure appropriate memory limits
- Monitor Redis performance
- Set up proper authentication and TLS 