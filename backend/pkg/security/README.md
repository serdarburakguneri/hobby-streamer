# Security Package

Middleware package for backend services that adds essential security protections, including rate limiting, security headers, CORS handling, and input validation.

## Features

- **Rate Limiting** – Per-IP rate limiting with sliding window
- **Security Headers** – Adds common headers for browser security
- **CORS Protection** – Configurable CORS with strict origin control
- **Input Validation** – Enforces request size and content-type rules

---

## Usage Example

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

Limit incoming requests by client IP:

```go
// Allow 100 requests per minute
router.Use(security.RateLimitMiddleware(100, time.Minute))
```

---

## CORS Configuration

Control allowed origins, methods, and headers:

```go
allowedOrigins := []string{
    "http://localhost:3000",
    "https://yourdomain.com",
}

allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
allowedHeaders := []string{"Content-Type", "Authorization", "X-Requested-With"}

router.Use(security.CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders))
```

---

## Security Headers

The middleware adds the following headers by default:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy` – Defaults to safe inline and script policies

---

