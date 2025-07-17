# Auth Package

A shared authentication and authorization library for backend services. Integrates with Keycloak and provides middleware for validating user and service tokens.

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"

// Create a token validator
validator := auth.NewKeycloakValidator(
    "http://localhost:8080",  // Keycloak base URL
    "hobby",                  // Realm
    "asset-manager",          // Client ID
)

// Create middleware
middleware := auth.NewAuthMiddleware(validator)

// Apply to routes using builder
router.Use(func(next http.Handler) http.Handler {
    return middleware.RequireUserAuth().RequireServiceAuth().Build()(next.ServeHTTP)
})
```

## Builder Usage

The middleware supports flexible composition:

```go
// Only require user authentication
middleware.RequireUserAuth().Build()

// Only require service authentication
middleware.RequireServiceAuth().Build()

// Require both
middleware.RequireUserAuth().RequireServiceAuth().Build()
```

## Role-Based Authorization

Handlers can be restricted by roles:

```go
// Require a specific role
middleware.RequireRole("admin")(handler)

// Require any one of the given roles
middleware.RequireAnyRole([]string{"admin", "editor"})(handler)

// Require all listed roles
middleware.RequireAllRoles([]string{"admin", "moderator"})(handler)
```

## Context Keys

When authentication is successful, the following keys are set in the request context:

- `"user"` – Authenticated user (JWT subject and claims)
- `"service_user"` – Service-level identity for internal communication