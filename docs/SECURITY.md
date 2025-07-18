# Security Implementation

This document outlines the security measures implemented in the hobby-streamer project.

## Security Improvements

### 1. Rate Limiting

**Implementation**: `backend/pkg/security/middleware.go`

- **Configurable rate limiting** per client IP address
- **Sliding window algorithm** for accurate request counting
- **Default settings**: 100 requests per minute per IP
- **Configurable** via configuration files

```yaml
security:
  rate_limit:
    requests: 100
    window: "1m"
```

### 2. CORS Protection

**Implementation**: `backend/pkg/security/middleware.go`

- **Origin validation** instead of wildcard `*`
- **Configurable allowed origins** via configuration
- **Proper CORS headers** with credentials support
- **Preflight request handling**

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

### 3. Security Headers

**Implementation**: `backend/pkg/security/middleware.go`

Automatically adds essential security headers:

- `X-Content-Type-Options: nosniff` - Prevents MIME type sniffing
- `X-Frame-Options: DENY` - Prevents clickjacking attacks
- `X-XSS-Protection: 1; mode=block` - XSS protection
- `Referrer-Policy: strict-origin-when-cross-origin` - Controls referrer information
- `Content-Security-Policy` - Restricts resource loading

### 4. Input Validation

**Implementation**: `backend/pkg/security/middleware.go`

- **Request size limits** (default: 10MB)
- **Content type validation** for POST/PUT/PATCH requests
- **Path traversal protection** in file operations
- **Input sanitization** in lambda functions

### 5. WebSocket Security

**Implementation**: `backend/asset-manager/internal/config/graphql.go`

- **Origin validation** for WebSocket connections
- **Restricted to allowed origins** only
- **No wildcard origin acceptance**

### 6. Authentication & Authorization

**Implementation**: `backend/pkg/auth/`

- **JWT token validation** with Keycloak integration
- **Role-based access control** (RBAC)
- **Service-to-service authentication**
- **Token expiration handling**

### 7. API Gateway Security

**Implementation**: `local/setup-api-gateway.sh`

- **Restricted CORS origins** (no wildcards)
- **Proper CORS headers** configuration
- **Method and header restrictions**

### 8. Lambda Function Security

**Implementation**: `backend/lambdas/`

- **Input validation** and sanitization
- **Path traversal protection**
- **Proper error handling** without information leakage
- **Restricted CORS headers**

## Configuration

### Development Environment

```yaml
security:
  rate_limit:
    requests: 100
    window: "1m"
  max_request_size: 10485760  # 10MB
  cors:
    allowed_origins:
      - "http://localhost:3000"
      - "http://localhost:8081"
      - "http://localhost:8080"
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

### Production Considerations

For production deployment, consider:

1. **HTTPS only** - All endpoints should use HTTPS
2. **Strict CORS** - Limit origins to production domains only
3. **Higher rate limits** - Adjust based on expected traffic
4. **Monitoring** - Implement security event logging
5. **Secrets management** - Use proper secrets management for sensitive data

## Security Headers Explained

### Content Security Policy (CSP)

```
default-src 'self'; 
script-src 'self' 'unsafe-inline'; 
style-src 'self' 'unsafe-inline'; 
img-src 'self' data: https:; 
font-src 'self' data:; 
connect-src 'self' ws: wss:;
```

- **default-src 'self'** - Only allow resources from same origin
- **script-src 'self' 'unsafe-inline'** - Allow scripts from same origin and inline scripts
- **style-src 'self' 'unsafe-inline'** - Allow styles from same origin and inline styles
- **img-src 'self' data: https:** - Allow images from same origin, data URIs, and HTTPS sources
- **font-src 'self' data:** - Allow fonts from same origin and data URIs
- **connect-src 'self' ws: wss:** - Allow connections to same origin and WebSocket connections

## Rate Limiting Algorithm

The rate limiting uses a **sliding window** algorithm:

1. **Track requests** per client IP with timestamps
2. **Clean old requests** outside the time window
3. **Count valid requests** within the window
4. **Reject requests** when limit is exceeded
5. **Return 429 Too Many Requests** for rate-limited requests

## Monitoring & Logging

Security events are logged with:

- **Request details** (IP, method, path, user agent)
- **Rate limiting events** (when limits are exceeded)
- **Authentication failures** (invalid tokens, missing auth)
- **Input validation failures** (malformed requests)

## Future Enhancements

Consider implementing:

1. **IP whitelisting/blacklisting**
2. **Advanced threat detection**
3. **Request fingerprinting**
4. **Bot detection**
5. **DDoS protection**
6. **API key management**
7. **Audit logging**
8. **Security metrics collection** 