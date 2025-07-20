# Security Package

Middleware for backend services that adds essential security protections — including rate limiting, CORS controls, security headers, and input validation. Helps you cover the basics without a ton of boilerplate.

---

## Features

-  **Rate limiting** — Per-IP limits using a sliding window
-  **Security headers** — Sensible defaults to protect against common browser threats
-  **CORS** — Strict, configurable cross-origin rules
-  **Input validation** — Enforces max body size and content type

---

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

---

## Rate Limiting

Limit requests per IP using a sliding time window:

```go
// Max 100 requests per minute per client IP
router.Use(security.RateLimitMiddleware(100, time.Minute))
```

---

## CORS Configuration

Set allowed origins, methods, and headers:

```go
router.Use(security.CORSMiddleware(
    []string{
        "http://localhost:3000",
        "https://yourdomain.com",
    },
    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    []string{"Content-Type", "Authorization", "X-Requested-With"},
))
```

---

## Security Headers

The middleware automatically applies these:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy` — Conservative defaults with inline/script restrictions

---

>  This package is built for internal use — lightweight and flexible enough to work in both local and containerized environments.