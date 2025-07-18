# Logger Package

A structured logging helper built on top of Go’s `slog`. Supports both sync and async logging, and includes middleware for tracking requests and compressing responses. Meant to keep logs clean, contextual, and useful — without getting in your way.

---

## Features

- JSON and text output formats
- Async logging (optional, but handy under load)
- Request logging middleware with trace IDs
- Gzip compression middleware
- Context-aware logging (user ID, service name, etc.)
- Works well with structured error handling

---

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

// Set up sync logger (stdout)
logger.Init(slog.LevelInfo, "text")

logger.Info("App started")
logger.Error("Something broke", "error", err)
```

---

## Async Logging

For higher-throughput scenarios, async logging helps avoid blocking on log writes:

```go
// Use async logger with a buffer size
logger.InitAsync(slog.LevelInfo, "json", 1000)

logger.Info("This log is non-blocking")
logger.Error("Async error", "error", err)

// Make sure logs flush before shutdown
defer logger.Close()
```

### Config Sample

```yaml
log:
  level: info
  format: json
  async:
    enabled: true
    buffer_size: 1000
```

---

## Middleware

### Request Logging

Adds trace IDs and logs incoming HTTP requests:

```go
router.Use(logger.RequestLoggingMiddleware(logger.Get()))
```

### Compression

Adds gzip compression when the client supports it:

```go
router.Use(logger.CompressionMiddleware)
```

---

## Context Support

Add context to log entries:

```go
ctx := context.WithValue(context.Background(), "user_id", "abc123")
logger.WithContext(ctx).Info("User action")
```

---

## Add Service Info

Tag logs with your service name:

```go
logger.WithService("streaming-api").Info("Service is up")
```

---

> ℹ️ This logger is meant to be lightweight and flexible — good defaults, with room to grow if you need more control.