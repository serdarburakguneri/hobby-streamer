# Auth Package

Shared authentication package for backend services. Integrates with Keycloak for user/service token validation and role-based access control.

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"

// Create validator
validator := auth.NewKeycloakValidator("http://localhost:8080", "hobby", "asset-manager")

// Build middleware
middleware := auth.NewAuthMiddleware(validator)

// Apply to routes
router.Use(func(next http.Handler) http.Handler {
    return middleware.RequireUserAuth().RequireServiceAuth().Build()(next.ServeHTTP)
})
```

## Middleware Composition

```go
middleware.RequireUserAuth().Build()           // User token only
middleware.RequireServiceAuth().Build()        // Service token only
middleware.RequireUserAuth().RequireServiceAuth().Build()  // Both
```

## Role-Based Authorization

```go
middleware.RequireRole("admin")(handler)                    // Single role
middleware.RequireAnyRole([]string{"admin", "editor"})(handler)  // Any role
middleware.RequireAllRoles([]string{"admin", "moderator"})(handler)  // All roles
```

## Context Values

When auth succeeds, middleware injects:
- `"user"` – Authenticated user's JWT payload
- `"service_user"` – Calling service identity (if applicable)

> ⚠️ Internal use within Hobby Streamer services. Not for external use or production without review.