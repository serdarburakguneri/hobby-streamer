# Error Handling & Resilience Package

Error handling and resilience patterns for backend services.

## Features

- **Custom Error Types**: Structured error classification with context
- **Retry Mechanisms**: Exponential backoff with jitter for transient failures
- **Circuit Breaker Pattern**: Prevents cascading failures from external services
- **Graceful Degradation**: Fallback mechanisms for service degradation
- **Error Context**: Enhanced error information for debugging

## Error Types

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Create typed errors
err := errors.NewValidationError("invalid input", nil)
err := errors.NewNotFoundError("asset not found", nil)
err := errors.NewExternalError("external service failed", cause)
err := errors.NewTransientError("temporary failure", cause)

// Check error types
if errors.IsTransient(err) {
    // Handle transient error
}

if errors.IsExternal(err) {
    // Handle external service error
}
```

## Retry Mechanisms

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Simple retry with default config
err := errors.Retry(ctx, func(ctx context.Context) error {
    return someOperation()
}, nil)

// Custom retry configuration
config := &errors.RetryConfig{
    MaxAttempts:   5,
    InitialDelay:  100 * time.Millisecond,
    MaxDelay:      10 * time.Second,
    BackoffFactor: 2.0,
    JitterFactor:  0.1,
    RetryableErrors: []errors.ErrorType{
        errors.ErrorTypeTransient,
        errors.ErrorTypeTimeout,
    },
}

err := errors.Retry(ctx, func(ctx context.Context) error {
    return someOperation()
}, config)

// Quick retry for fast operations
err := errors.RetryFast(ctx, func(ctx context.Context) error {
    return someOperation()
})
```

## Circuit Breaker

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Create circuit breaker
breaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
    Name:      "external-api",
    Threshold: 5,
    Timeout:   30 * time.Second,
    OnStateChange: func(name string, from, to errors.CircuitState) {
        log.Printf("Circuit breaker %s: %v -> %v", name, from, to)
    },
})

// Use circuit breaker
err := breaker.Execute(ctx, func() error {
    return callExternalAPI()
})

// Check circuit state
if breaker.State() == errors.CircuitOpen {
    // Handle open circuit
}
```

## Circuit Breaker Registry

```go
// Global registry for managing multiple circuit breakers
registry := errors.NewCircuitBreakerRegistry()

// Get or create circuit breaker
breaker := registry.GetOrCreate("asset-manager", errors.CircuitBreakerConfig{
    Threshold: 3,
    Timeout:   60 * time.Second,
})

// Use in services
err := breaker.Execute(ctx, func() error {
    return assetManagerService.Call()
})
```

## Graceful Degradation

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"

// Create fallback chain
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

// Execute with fallback
result := fallback.Execute(ctx)
if result.Success {
    log.Printf("Operation succeeded using: %s", result.Used)
} else {
    log.Printf("All fallbacks failed: %v", result.Error)
}
```

## Cache Fallback

```go
// Simple cache fallback
cacheFallback := errors.NewCacheFallback(cachedData, fallbackData)
data := cacheFallback.Get()

// Update cache
cacheFallback.SetCache(newData)
```

## Degradation Manager

```go
// Manage service degradation levels
degradationManager := errors.NewDegradationManager()

degradationManager.OnLevelChange(errors.DegradationPartial, func() {
    log.Println("Service degraded to partial mode")
})

degradationManager.OnLevelChange(errors.DegradationFull, func() {
    log.Println("Service fully degraded")
})

// Set degradation level
degradationManager.SetLevel(errors.DegradationPartial)

// Check degradation status
if degradationManager.IsDegraded() {
    // Use fallback mechanisms
}
```

## Integration with HTTP Handlers

```go
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    asset, err := h.service.GetAsset(ctx, slug)
    if err != nil {
        if errors.IsAppError(err) {
            appErr := err.(*errors.AppError)
            switch appErr.Type {
            case errors.ErrorTypeNotFound:
                h.writeError(w, http.StatusNotFound, appErr.Message)
                return
            case errors.ErrorTypeValidation:
                h.writeError(w, http.StatusBadRequest, appErr.Message)
                return
            case errors.ErrorTypeUnauthorized:
                h.writeError(w, http.StatusUnauthorized, appErr.Message)
                return
            case errors.ErrorTypeForbidden:
                h.writeError(w, http.StatusForbidden, appErr.Message)
                return
            case errors.ErrorTypeConflict:
                h.writeError(w, http.StatusConflict, appErr.Message)
                return
            default:
                h.writeError(w, http.StatusInternalServerError, "Internal server error")
                return
            }
        }
        
        h.writeError(w, http.StatusInternalServerError, "Internal server error")
        return
    }
    
    h.writeJSON(w, http.StatusOK, asset)
}
```

## Error Context

```go
// Add context to errors
err = errors.WithContext(err, map[string]interface{}{
    "user_id": userID,
    "asset_id": assetID,
    "operation": "create_asset",
})

// Access context in error handlers
if appErr, ok := err.(*errors.AppError); ok {
    for key, value := range appErr.Context {
        log.Printf("Error context %s: %v", key, value)
    }
}
``` 