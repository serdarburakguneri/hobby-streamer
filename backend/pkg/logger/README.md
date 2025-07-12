# Structured Logging Package

This package provides a centralized, structured logging solution for all backend services in the hobby-streamer project.

## Features

- Structured Logging: Uses Go's built-in `log/slog` package for structured, JSON-formatted logs
- Log Levels: Support for DEBUG, INFO, WARN, ERROR levels
- Context Awareness: Automatic inclusion of request context, user information, and custom fields
- Service Identification: Each service gets its own logger with service name
- HTTP Request Logging: Middleware for automatic HTTP request/response logging
- Error Context: Enhanced error logging with stack traces and context

## Usage

### Basic Setup

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

func main() {
    // Initialize logger with level and format
    logger.Init(slog.LevelInfo, "json") // or "text"
    
    // Get service-specific logger
    log := logger.WithService("my-service")
    
    log.Info("Service started")
}
```

### Log Levels

```go
log.Debug("Debug information", "key", "value")
log.Info("Information message", "user_id", 123)
log.Warn("Warning message", "attempt", 3)
log.Error("Error occurred", "operation", "database_query")
```

### Error Logging

```go
err := someOperation()
if err != nil {
    log.WithError(err).Error("Operation failed", "operation", "database_query")
}
```

### Context-Aware Logging

```go
// In HTTP handlers
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    log := h.logger.WithContext(r.Context())
    log.Info("Request processed", "user_id", getUserID(r))
}
```

### HTTP Request Logging Middleware

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

func main() {
    log := logger.WithService("my-service")
    
    // Apply logging middleware
    handler := logger.RequestLoggingMiddleware(log)(yourHandler)
    
    http.ListenAndServe(":8080", handler)
}
```

## Environment Variables

Configure logging behavior using environment variables:

- `LOG_LEVEL`: Set log level (debug, info, warn, error) - default: info
- `LOG_FORMAT`: Set log format (json, text) - default: text

## Log Output Examples

### Text Format
```
time=2024-01-15T10:30:00.000Z level=INFO service=asset-manager msg="Asset created successfully" asset_id=123 title="My Video"
```

### JSON Format
```json
{
  "time": "2024-01-15T10:30:00.000Z",
  "level": "INFO",
  "service": "asset-manager",
  "msg": "Asset created successfully",
  "asset_id": 123,
  "title": "My Video"
}
```

### HTTP Request Log
```
time=2024-01-15T10:30:00.000Z level=INFO service=asset-manager msg="HTTP request completed" method=POST path=/assets status_code=201 duration_ms=45 user_id=123
```

## Integration with Services

All backend services have been updated to use this logging system:

- asset-manager: Asset management operations
- auth-service: Authentication and authorization
- transcoder: Video processing and transcoding

## Benefits

1. Consistency: All services use the same logging format and levels
2. Observability: Structured logs make it easier to monitor and debug
3. Performance: Efficient logging with minimal overhead
4. Context: Rich context information for better debugging
5. Standards: Uses Go's official structured logging package 