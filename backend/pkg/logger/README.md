# Logger Package

A structured logging package built on top of Go's `slog` with support for both synchronous and asynchronous logging.

## Features

- Structured logging with JSON and text formats
- Request tracking with unique IDs
- HTTP request logging middleware
- Compression middleware
- **Async logging support** for improved performance
- Context-aware logging
- Service and error context

## Usage

### Basic Usage

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

// Initialize synchronous logger
logger.Init(slog.LevelInfo, "text")

// Log messages
logger.Info("Application started")
logger.Error("Something went wrong", "error", err)
```

### Async Logging

For high-throughput applications, use async logging to avoid blocking:

```go
// Initialize async logger with buffer size
logger.InitAsync(slog.LevelInfo, "json", 1000)

// Log messages (non-blocking)
logger.Info("High-volume logging")
logger.Error("Error occurred", "error", err)

// Clean shutdown
defer logger.Close()
```

### Configuration

Enable async logging in your config:

```yaml
log:
  level: info
  format: json
  async:
    enabled: true
    buffer_size: 1000
```


## Middleware

### Request Logging

```go
handler := logger.RequestLoggingMiddleware(logger.Get())(yourHandler)
```

### Compression

```go
handler := logger.CompressionMiddleware(yourHandler)
```

## Context Support

```go
// Add context to logger
ctx := context.WithValue(context.Background(), "user_id", "123")
logger.WithContext(ctx).Info("User action")
```

## Service Context

```go
logger.WithService("my-service").Info("Service started")
``` 