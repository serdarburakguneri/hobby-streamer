# Security Overview

This document outlines the security mechanisms implemented in the **Hobby Streamer** project. While this is a personal project and not production-grade, the goal is to follow good practices and build a secure-by-default foundation during development and testing.

## Implemented Measures

### 1. Rate Limiting

**Location**: `backend/pkg/security/middleware.go`

- Limits requests per client IP using a **sliding window** algorithm  
- Default: **100 requests per minute per IP**  
- Fully configurable through YAML:

```yaml
security:
  rate_limit:
    requests: 100
    window: "1m"
```

Requests exceeding the threshold are rejected with a `429 Too Many Requests` status.

---

### 2. CORS Protection

**Location**: `backend/pkg/security/middleware.go`

- Only allows **whitelisted origins**
- Supports configuration of allowed origins, methods, and headers
- Properly handles preflight requests

Example configuration:

```yaml
security:
  cors:
    allowed_origins:
      - "http://localhost:3000"
      - "http://localhost:8081"
    allowed_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "OPTIONS"
    allowed_headers:
      - "Content-Type"
      - "Authorization"
      - "X-Requested-With"
```

---

### 3. Security Headers

**Location**: `backend/pkg/security/middleware.go`

Automatically adds the following headers to enhance browser-level security:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy` (see below)

---

### 4. Input Validation

**Location**: `backend/pkg/security/middleware.go`

- Enforces **maximum request size** (default: 10MB)
- Validates **Content-Type** headers for applicable HTTP methods
- Protects against **path traversal** in file operations
- Adds basic input sanitization in lambda functions

---

### 5. WebSocket Origin Control

**Location**: `backend/asset-manager/internal/config/graphql.go`

- WebSocket connections are restricted to **trusted origins only**
- Wildcard origins (`*`) are explicitly disallowed

---

### 6. Authentication & Authorization

**Location**: `backend/pkg/auth/`

- **JWT validation** using Keycloak
- **Role-based access control** (RBAC) with scoped permissions
- **Service-to-service authentication** with token verification
- Proper handling of **token expiration**

---

### 7. Lambda Function Security

**Location**: `backend/lambdas/`

- Input validation and sanitization at entry points
- Protection against path traversal vulnerabilities
- Clear and safe error handling to avoid leaking internal details
- CORS headers restricted to allowed origins

---

## Content Security Policy (CSP)

The CSP used in HTTP responses is defined as:

```
default-src 'self'; 
script-src 'self' 'unsafe-inline'; 
style-src 'self' 'unsafe-inline'; 
img-src 'self' data: https:; 
font-src 'self' data:; 
connect-src 'self' ws: wss:;
```

This restricts loading of external resources and helps reduce the attack surface for XSS and injection attacks.

---

## About the Rate Limiting Algorithm

The rate limiter is based on a **sliding window** strategy:

1. Each client IP’s requests are timestamped.
2. Requests outside the configured time window are discarded.
3. Remaining requests are counted to determine if the threshold is exceeded.
4. If the limit is breached, the request is rejected with HTTP `429`.

---

> **Note:** While these mechanisms are designed with security in mind, this project is still experimental and under active development. Features are continuously improved, and best practices are applied incrementally.