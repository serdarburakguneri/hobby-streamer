# Auth Service

Lightweight Go service for JWT-based authentication using Keycloak. Provides login, token validation, and refresh endpoints.

## Features

Username/password login, JWT validation with role extraction, token refresh, health check.

## API

### `POST /login`
Authenticate user and get access/refresh tokens.

**Request:** `{"username": "user", "password": "pass", "client_id": "asset-manager"}`  
**Response:** `{"access_token": "...", "refresh_token": "...", "expires_in": 300}`

### `POST /validate`
Validate token and return user info.

**Request:** `{"token": "Bearer ..."}`  
**Response:** `{"valid": true, "user": {"id": "...", "username": "...", "roles": ["user"]}}`

### `POST /refresh`
Refresh access token.

**Request:** `{"refresh_token": "..."}`  
**Response:** `{"access_token": "...", "refresh_token": "...", "expires_in": 300}`

### `GET /health`
Health check. **Response:** `{"status": "ok"}`

## Running

```bash
# Local
cd backend/auth-service && go run cmd/main.go

# Docker
docker build -t auth-service . && docker run -p 8080:8080 auth-service
```

> ⚠️ Designed for local development and testing.