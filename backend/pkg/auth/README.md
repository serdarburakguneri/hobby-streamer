# Auth Package

Shared Go package for authentication, integrates with Keycloak for user/service token validation and RBAC.

## Features
Keycloak integration, JWT validation, role-based access, middleware composition, context injection.

## Quick Start
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
validator := auth.NewKeycloakValidator("http://localhost:8080", "hobby", "asset-manager")
middleware := auth.NewAuthMiddleware(validator)
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
middleware.RequireRole("admin")(handler)
middleware.RequireAnyRole([]string{"admin", "editor"})(handler)
middleware.RequireAllRoles([]string{"admin", "moderator"})(handler)
```

## Context Values
On success: "user" (JWT payload), "service_user" (service identity).
