# Auth Service

A lightweight Go service that handles JWT-based authentication using Keycloak. Provides login, token validation, and refresh endpoints, with structured logging and clean error handling.

---

## Features

- Username/password login via Keycloak
- JWT token validation with role extraction
- Token refresh endpoint
- Basic health check

---

## API Overview

### `POST /login`

Authenticate a user and retrieve access/refresh tokens.

**Request:**

```json
{
  "username": "testuser",
  "password": "testpass",
  "client_id": "asset-manager"
}
```

**Success Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

**Error Response:**

```json
{
  "error": "invalid_credentials",
  "message": "Invalid username or password"
}
```

---

### `POST /validate`

Validate a token and return decoded user info.

**Request:**

```json
{
  "token": "Bearer eyJ..."
}
```

**Success Response:**

```json
{
  "valid": true,
  "user": {
    "id": "user-id",
    "username": "testuser",
    "email": "test@example.com",
    "roles": ["user"]
  },
  "roles": ["user"]
}
```

**Error Response:**

```json
{
  "error": "invalid_token",
  "message": "Token is invalid or expired"
}
```

---

### `POST /refresh`

Refresh an access token using a refresh token.

**Request:**

```json
{
  "refresh_token": "eyJ..."
}
```

**Success Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

**Error Response:**

```json
{
  "error": "invalid_refresh_token",
  "message": "Refresh token is invalid or expired"
}
```

---

### `GET /health`

Simple health check.

**Response:**

```json
{
  "status": "ok"
}
```

---

## Running the Service

### Locally

```bash
cd backend/auth-service
go run cmd/main.go
```

### With Docker

```bash
docker build -t auth-service .
docker run -p 8080:8080 auth-service
```

---

> ⚠️ This service is designed for local development and testing in the context of the Hobby Streamer platform.