# Error Handling & Resilience Package

Shared library for error handling and resilience patterns used across backend services. Includes support for typed errors, retries, circuit breakers, and graceful degradation.

## Features

- Typed application errors with context
- Retry logic with backoff and jitter
- Circuit breaker support for external dependencies
- Fallback chains for graceful degradation
- Cache-aware fallback strategies
- Degradation state management
- Error context for better logging and debugging

---

## Error Types

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Create common typed errors
err := errors.NewValidationError("invalid input", nil)
err := errors.NewNotFoundError("asset not found", nil)
err := errors.NewExternalError("external service failed", cause)
err := errors.NewTransientError("temporary failure", cause)

// Check error types
if errors.IsTransient(err) {
    // Retry or fallback
}
```

---

## Retry Mechanisms

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Default retry
err := errors.Retry(ctx, func(ctx context.Context) error {
    return someOperation()
}, nil)

// Custom retry config
config := &errors.RetryConfig{
    MaxAttempts:      5,
    InitialDelay:     100 * time.Millisecond,
    MaxDelay:         10 * time.Second,
    BackoffFactor:    2.0,
    JitterFactor:     0.1,
    RetryableErrors:  []errors.ErrorType{errors.ErrorTypeTransient},
}

err := errors.Retry(ctx, someOperation, config)

// Fast retry for short-lived operations
err := errors.RetryFast(ctx, someOperation)
```

---

## Circuit Breaker

```go
breaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
    Name:      "external-api",
    Threshold: 5,
    Timeout:   30 * time.Second,
    OnStateChange: func(name string, from, to errors.CircuitState) {
        log.Printf("Circuit breaker %s: %v -> %v", name, from, to)
    },
})

err := breaker.Execute(ctx, callExternalAPI)

if breaker.State() == errors.CircuitOpen {
    // Fallback or error
}
```

### Circuit Breaker Registry

```go
registry := errors.NewCircuitBreakerRegistry()

breaker := registry.GetOrCreate("asset-manager", errors.CircuitBreakerConfig{
    Threshold: 3,
    Timeout:   60 * time.Second,
})

err := breaker.Execute(ctx, func() error {
    return assetManagerService.Call()
})
```

---

## Graceful Degradation

```go
fallback := errors.NewFallbackChain(
    func(ctx context.Context) error {
        return primaryService.Call()
    },
    "primary",
).AddFallback(
    func(ctx context.Context) error {
        return secondaryService.Call()
    },
    "secondary",
).AddFallback(
    func(ctx context.Context) error {
        return cacheService.Get()
    },
    "cache",
)

result := fallback.Execute(ctx)
if result.Success {
    log.Printf("Used: %s", result.Used)
} else {
    log.Printf("All fallbacks failed: %v", result.Error)
}
```

---

## Cache Fallback

```go
cacheFallback := errors.NewCacheFallback(cachedData, fallbackData)
data := cacheFallback.Get()

// Update cache
cacheFallback.SetCache(newData)
```

---

## Degradation Manager

```go
degradationManager := errors.NewDegradationManager()

degradationManager.OnLevelChange(errors.DegradationPartial, func() {
    log.Println("Service degraded: partial")
})

degradationManager.OnLevelChange(errors.DegradationFull, func() {
    log.Println("Service degraded: full")
})

degradationManager.SetLevel(errors.DegradationPartial)

if degradationManager.IsDegraded() {
    // Use fallback logic
}
```

---

## Integration with HTTP Handlers

```go
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    asset, err := h.service.GetAsset(ctx, slug)

    if err != nil {
        h.handleError(w, err, "Failed to get asset")
        return
    }

    h.writeJSON(w, http.StatusOK, asset)
}

func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
    if errors.IsAppError(err) {
        appErr := err.(*errors.AppError)
        h.logger.WithError(err).Error("Application error", "error_type", appErr.Type, "context", appErr.Context)

        status := appErr.HTTPStatus()
        message := appErr.Message
        if appErr.Type == errors.ErrorTypeCircuitBreaker {
            message = "Service temporarily unavailable"
        }

        h.writeError(w, status, message)
        return
    }

    h.logger.WithError(err).Error("Unexpected error")
    h.writeError(w, http.StatusInternalServerError, defaultMessage)
}
```

---

## Error Context

```go
err = errors.WithContext(err, map[string]interface{}{
    "user_id": userID,
    "asset_id": assetID,
    "operation": "create_asset",
})

if appErr, ok := err.(*errors.AppError); ok {
    for key, value := range appErr.Context {
        log.Printf("Error context %s: %v", key, value)
    }
}
```

## Message Consumer Integration

```go
func (c *Consumer) HandleMessage(ctx context.Context, payload map[string]interface{}) error {
    err := c.service.ProcessMessage(ctx, payload)
    if err != nil {
        log.WithError(err).Error("Failed to process message")
        return errors.WrapWithContext(err, "failed to process message")
    }
    return nil
}
```

The `WrapWithContext` function intelligently wraps errors with appropriate context while preserving the original error type:

- **Validation errors** → Wrapped as validation errors
- **Not found errors** → Wrapped as not found errors  
- **Transient errors** → Wrapped as transient errors
- **Other errors** → Wrapped as internal errors