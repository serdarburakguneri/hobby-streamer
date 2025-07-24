# Error Handling & Resilience Package

Consistent error handling and resilience for backend services.

## Features
Typed errors (validation, transient, external), retry with backoff, circuit breakers, fallback chains, cache-aware helpers, degradation tracking, error context for logs.

## Usage
### Typed Errors
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
err := errors.NewValidationError("invalid input", nil)
if errors.IsTransient(err) { /* retry/fallback */ }
```

### Retry
```go
err := errors.Retry(ctx, func(ctx context.Context) error { return doSomething() }, nil)
```

### Circuit Breaker
```go
breaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{Name: "external-api", Threshold: 5, Timeout: 30 * time.Second})
err := breaker.Execute(ctx, callExternalService)
```

### Fallback Chain
```go
fallback := errors.NewFallbackChain(primaryService.Call, "primary").AddFallback(secondaryService.Call, "secondary")
result := fallback.Execute(ctx)
```

## Integration
Use in HTTP handlers, message consumers, and service logic for consistent error modeling, retries, fallbacks, and observability.