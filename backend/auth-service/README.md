# Auth Service

Simple Go service for JWT authentication using Keycloak. Handles login, token validation, refresh, and health checks.

## Features
Username/password login, JWT validation, role extraction, token refresh, health check.

## API

### POST /auth/login
Authenticate and get tokens.
Request: `{ "username": "user", "password": "pass", "client_id": "asset-manager" }`
Response: `{ "access_token": "...", "refresh_token": "...", "expires_in": 300 }`

### POST /auth/validate
Validate token, get user info.
Request: `{ "token": "Bearer ..." }`
Response: `{ "valid": true, "user": { "id": "...", "username": "...", "roles": ["user"] } }`

### POST /auth/refresh
Refresh access token.
Request: `{ "refresh_token": "..." }`
Response: `{ "access_token": "...", "refresh_token": "...", "expires_in": 300 }`

### GET /health
Health check. Response: `{ "status": "ok" }`

## Running
```bash
cd backend/auth-service && go run cmd/main.go
# or with Docker
docker build -t auth-service . && docker run -p 8080:8080 auth-service
```
