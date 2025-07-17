# Auth Service

A small Go service for handling JWT authentication using Keycloak. Provides login, token validation, and refresh functionality, with structured logging and error handling.

## Features

- User authentication (login)
- JWT token validation
- Token refresh

## API Endpoints

### POST /login

Authenticates a user and returns access and refresh tokens.

**Request**
```json
{
  "username": "testuser",
  "password": "testpass",
  "client_id": "asset-manager"
}
```

**Response**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

**Error**
```json
{
  "error": "invalid_credentials",
  "message": "Invalid username or password"
}
```

---

### POST /validate

Validates a JWT token and returns user info.

**Request**
```json
{
  "token": "Bearer eyJ..."
}
```

**Response**
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

**Error**
```json
{
  "error": "invalid_token",
  "message": "Token is invalid or expired"
}
```

---

### POST /refresh

Refreshes an access token using a refresh token.

**Request**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

**Error**
```json
{
  "error": "invalid_refresh_token",
  "message": "Refresh token is invalid or expired"
}
```

---

### GET /health

Returns basic service health status.

**Response**
```json
{
  "status": "ok"
}
```

## Running the Service

### Local

```bash
cd backend/auth-service
go run cmd/main.go
```

### Docker

```bash
docker build -t auth-service .
docker run -p 8080:8080 auth-service
```