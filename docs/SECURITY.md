# Security Overview

Security mechanisms in Hobby Streamer: rate limiting, CORS, security headers, input validation, WebSocket origin control, JWT auth, RBAC, lambda input validation.

## Measures

### Rate Limiting
- Sliding window per-IP, default 100 req/min, YAML configurable.

### CORS
- Whitelisted origins, configurable methods/headers, preflight support.

### Security Headers
- X-Content-Type-Options: nosniff, X-Frame-Options: DENY, X-XSS-Protection: 1; mode=block, Referrer-Policy: strict-origin-when-cross-origin, Content-Security-Policy.

### Input Validation
- Max request size (default 10MB), Content-Type validation, path traversal protection, lambda input sanitization.

### WebSocket Origin Control
- Trusted origins only, no wildcards.

### Auth & Authorization
- JWT validation (Keycloak), RBAC, service-to-service auth, token expiration handling.

### Lambda Security
- Input validation, path traversal protection, safe error handling, CORS headers restricted.

## Content Security Policy
```
default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss:;
```

## Rate Limiting Algorithm
- Sliding window, timestamps per IP, requests outside window discarded, threshold enforced, HTTP 429 on breach.