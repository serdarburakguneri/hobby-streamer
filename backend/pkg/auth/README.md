# Auth Package

A shared authentication and authorization package used across backend services. Integrates with Keycloak and provides middleware for validating user and service tokens, along with support for role-based access control.

---

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"

// Create a validator instance
validator := auth.NewKeycloakValidator(
    "http://localhost:8080",  // Keycloak base URL
    "hobby",                  // Realm name
    "asset-manager",          // Client ID
)

// Build middleware
middleware := auth.NewAuthMiddleware(validator)

// Apply to routes
router.Use(func(next http.Handler) http.Handler {
    return middleware.
        RequireUserAuth().
        RequireServiceAuth().
        Build()(next.ServeHTTP)
})
```

---

## Middleware Composition

You can configure the middleware depending on your route requirements:

```go
// Require user token only
middleware.RequireUserAuth().Build()

// Require service token only
middleware.RequireServiceAuth().Build()

// Require both user and service tokens
middleware.RequireUserAuth().RequireServiceAuth().Build()
```

---

## Role-Based Authorization

Handlers can be gated based on user roles:

```go
// Require a single role
middleware.RequireRole("admin")(handler)

// Require any one of the listed roles
middleware.RequireAnyRole([]string{"admin", "editor"})(handler)

// Require all specified roles
middleware.RequireAllRoles([]string{"admin", "moderator"})(handler)
```

---

## Context Values

When authentication succeeds, the middleware injects the following into the request context:

- `"user"` – The authenticated user’s JWT payload
- `"service_user"` – Identity of the calling service (if applicable)

---

> ⚠️ This package is designed for internal use within Hobby Streamer services. It's not meant for external use or production deployment without review.