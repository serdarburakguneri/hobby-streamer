# Logger Package

Structured logging helper built on `slog`. Supports sync/async logging, request tracking middleware, and context-aware logging.

## Features

JSON/text output, async logging, request logging middleware with trace IDs, gzip compression, context-aware logging, structured error handling.

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

// Sync logger
logger.Init(slog.LevelInfo, "text")
logger.Info("App started")
logger.Error("Something broke", "error", err)
```

## Async Logging

```go
// Async logger with buffer
logger.InitAsync(slog.LevelInfo, "json", 1000)
logger.Info("Non-blocking log")
defer logger.Close()

// Config
log:
  level: info
  format: json
  async:
    enabled: true
    buffer_size: 1000
```

## Middleware

```go
// Request logging with trace IDs
router.Use(logger.RequestLoggingMiddleware(logger.Get()))

// Gzip compression
router.Use(logger.CompressionMiddleware)
```

## Context Support

```go
ctx := context.WithValue(context.Background(), "user_id", "abc123")
logger.WithContext(ctx).Info("User action")
logger.WithService("streaming-api").Info("Service is up")
```

> ℹ️ Lightweight and flexible — good defaults with room to grow.