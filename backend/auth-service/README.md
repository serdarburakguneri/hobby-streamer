# Auth Service

Go-based authentication service that handles JWT flows with Keycloak.

## Features

- User authentication (login)
- JWT token validation
- Token refresh
- Health checks

## Dependencies

- Keycloak: Identity and access management
- JWT: Token handling and validation
- Gorilla Mux: HTTP routing

## Configuration

Environment variables:
- `KEYCLOAK_URL`: Keycloak server URL (default: http://localhost:8080)
- `KEYCLOAK_REALM`: Keycloak realm name (default: hobby)
- `KEYCLOAK_CLIENT_ID`: Keycloak client ID (default: asset-manager)
- `KEYCLOAK_CLIENT_SECRET`: Keycloak client secret (optional for public clients)

## API Endpoints

### POST /login
Authenticates a user and returns a JWT token.

**Request:**
```json
{
  "username": "testuser",
  "password": "testpass",
  "client_id": "asset-manager"
}
```

**Response:**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

### POST /validate
Validates a JWT token and returns user information.

**Request:**
```json
{
  "token": "Bearer eyJ..."
}
```

**Response:**
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

### POST /refresh
Refreshes an expired JWT token.

**Request:**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response:**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 300,
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

### GET /health
Returns service health status.

**Response:**
```json
{
  "status": "ok"
}
```

## Running the Service

### Local Development
```bash
cd backend/auth-service
go run cmd/main.go
```

### With Docker
```bash
docker build -t auth-service .
docker run -p 8080:8080 auth-service
```
