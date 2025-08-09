# Auth Service

Simple Go service for JWT authentication using Keycloak. Handles login, token validation, refresh, and health checks.

## Features
Username/password login, JWT validation, role extraction, token refresh, health check.

## API
`POST /auth/login`, `POST /auth/validate`, `POST /auth/refresh`, `GET /health`.

## Running
```bash
cd backend/auth-service && go run cmd/main.go
```
