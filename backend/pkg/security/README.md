# Security Package

Middleware for backend services: rate limiting, CORS, security headers, input validation.

## Features
Rate limiting (per-IP, sliding window), security headers, strict CORS, input validation (body size, content type).

## Quick Example
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
router.Use(security.SecurityHeadersMiddleware())
router.Use(security.RateLimitMiddleware(100, time.Minute))
router.Use(security.CORSMiddleware(
    []string{"http://localhost:3000", "https://yourdomain.com"},
    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    []string{"Content-Type", "Authorization", "X-Requested-With"},
))
router.Use(security.InputValidationMiddleware())
```

## Security Headers
- X-Content-Type-Options: nosniff, X-Frame-Options: DENY, X-XSS-Protection: 1; mode=block, Referrer-Policy: strict-origin-when-cross-origin, Content-Security-Policy: conservative defaults.
