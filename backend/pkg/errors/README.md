# Error Handling & Resilience Package

Shared library for consistent error handling and resilience across backend services. Includes support for typed errors, retries with backoff, circuit breakers, graceful degradation, and contextual logging.

---

## Features

- Typed application-level errors (validation, transient, external, etc.)
- Retry utilities with exponential backoff and jitter
- Circuit breaker support for external dependencies
- Fallback chains for graceful degradation
- Cache-aware fallback helpers
- Degradation state tracking (partial/full)
- Error context attachment for structured logs and debugging

---

## Typed Errors

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Create typed errors
err := errors.NewValidationError("invalid input", nil)
err := errors.NewNotFoundError("asset not found", nil)
err := errors.NewExternalError("external call failed", cause)
err := errors.NewTransientError("temporary issue", cause)

// Type checking
if errors.IsTransient(err) {
    // Trigger retry/fallback logic
}
```

---

## Retry Logic

```go
// Basic retry
err := errors.Retry(ctx, func(ctx context.Context) error {
    return doSomething()
}, nil)

// Custom configuration
cfg := &errors.RetryConfig{
    MaxAttempts:     5,
    InitialDelay:    100 * time.Millisecond,
    MaxDelay:        10 * time.Second,
    BackoffFactor:   2.0,
    JitterFactor:    0.1,
    RetryableErrors: []errors.ErrorType{errors.ErrorTypeTransient},
}

err := errors.Retry(ctx, operation, cfg)

// Shortcut for fast retry
err := errors.RetryFast(ctx, operation)
```

---

## Circuit Breakers

```go
breaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
    Name:      "external-api",
    Threshold: 5,
    Timeout:   30 * time.Second,
    OnStateChange: func(name string, from, to errors.CircuitState) {
        log.Printf("Breaker %s: %v → %v", name, from, to)
    },
})

err := breaker.Execute(ctx, callExternalService)

if breaker.State() == errors.CircuitOpen {
    // Handle fallback logic
}
```

### With Registry

```go
registry := errors.NewCircuitBreakerRegistry()

breaker := registry.GetOrCreate("asset-manager", errors.CircuitBreakerConfig{
    Threshold: 3,
    Timeout:   60 * time.Second,
})

err := breaker.Execute(ctx, func() error {
    return assetManager.Fetch(...)
})
```

---

## Graceful Degradation

```go
fallback := errors.NewFallbackChain(
    primaryService.Call, "primary",
).AddFallback(
    secondaryService.Call, "secondary",
).AddFallback(
    cacheService.Get, "cache",
)

result := fallback.Execute(ctx)

if result.Success {
    log.Printf("Used fallback: %s", result.Used)
} else {
    log.Printf("All fallbacks failed: %v", result.Error)
}
```

---

## Cache Fallback Helper

```go
cacheFallback := errors.NewCacheFallback(cachedData, fallbackData)
value := cacheFallback.Get()

// Set new cache data
cacheFallback.SetCache(updatedData)
```

---

## Degradation State Manager

```go
manager := errors.NewDegradationManager()

manager.OnLevelChange(errors.DegradationPartial, func() {
    log.Println("Partial degradation detected")
})

manager.OnLevelChange(errors.DegradationFull, func() {
    log.Println("System in full degradation")
})

manager.SetLevel(errors.DegradationPartial)

if manager.IsDegraded() {
    // Activate fallback path
}
```

---

## Integration with HTTP Handlers

```go
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
    asset, err := h.service.GetAsset(r.Context(), slug)
    if err != nil {
        h.handleError(w, err, "Could not load asset")
        return
    }

    h.writeJSON(w, http.StatusOK, asset)
}

func (h *Handler) handleError(w http.ResponseWriter, err error, fallbackMsg string) {
    if errors.IsAppError(err) {
        appErr := err.(*errors.AppError)
        h.logger.WithError(err).Error("App error", "type", appErr.Type, "context", appErr.Context)

        status := appErr.HTTPStatus()
        msg := appErr.Message

        if appErr.Type == errors.ErrorTypeCircuitBreaker {
            msg = "Service temporarily unavailable"
        }

        h.writeError(w, status, msg)
        return
    }

    h.logger.WithError(err).Error("Unexpected error")
    h.writeError(w, http.StatusInternalServerError, fallbackMsg)
}
```

---

## Error Context

Attach custom context for debugging/logging:

```go
err = errors.WithContext(err, map[string]interface{}{
    "user_id":    userID,
    "asset_id":   assetID,
    "operation":  "create_asset",
})

// Access context
if appErr, ok := err.(*errors.AppError); ok {
    for k, v := range appErr.Context {
        log.Printf("Context %s: %v", k, v)
    }
}
```

---

## Message Consumer Example

```go
func (c *Consumer) HandleMessage(ctx context.Context, payload map[string]interface{}) error {
    err := c.service.Process(ctx, payload)
    if err != nil {
        log.WithError(err).Error("Failed to process message")
        return errors.WrapWithContext(err, "message handling failed")
    }
    return nil
}
```

The `WrapWithContext` helper preserves original error types:

- `ValidationError` → stays a validation error  
- `NotFoundError` → remains typed  
- `TransientError` → remains retryable  
- Other errors → wrapped as `InternalError`

---

> ✅ Use this package across services to enforce consistency in error modeling, retries, fallbacks, and observability.