# Logger Package

Structured logging for Go, built on slog. Supports sync/async logging, request tracking, context-aware logs.

## Features
JSON/text output, async logging, request logging with trace IDs, gzip compression, context support, structured error handling.

## Quick Start
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
logger.Init(slog.LevelInfo, "text")
logger.Info("App started")
logger.Error("Something broke", "error", err)
```

## Async Logging
```go
logger.InitAsync(slog.LevelInfo, "json", 1000)
logger.Info("Non-blocking log")
defer logger.Close()
```

## Middleware
```go
router.Use(logger.RequestLoggingMiddleware(logger.Get()))
router.Use(logger.CompressionMiddleware)
```

## Context Support
```go
ctx := context.WithValue(context.Background(), "user_id", "abc123")
logger.WithContext(ctx).Info("User action")
logger.WithService("streaming-api").Info("Service is up")
```