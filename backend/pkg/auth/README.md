# Auth Package

Authentication and authorization package for microservices with Keycloak integration.

## Quick Start

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"

// Create validator
validator := auth.NewKeycloakValidator(
    "http://localhost:8080",  // Keycloak URL
    "hobby",                  // Realm  
    "asset-manager",          // Client ID
)

// Create middleware with builder pattern
middleware := auth.NewAuthMiddleware(validator)

// Apply to routes
router.Use(func(next http.Handler) http.Handler {
    return middleware.RequireUserAuth().RequireServiceAuth().Build()(next.ServeHTTP)
})
```

## Builder Pattern

```go
// User authentication only
middleware.RequireUserAuth().Build()

// Service authentication only  
middleware.RequireServiceAuth().Build()

// Both user and service authentication
middleware.RequireUserAuth().RequireServiceAuth().Build()
```

## Role-Based Authorization

```go
// Require specific role
middleware.RequireRole("admin")(handler)

// Require any of multiple roles
middleware.RequireAnyRole([]string{"admin", "editor"})(handler)

// Require all roles
middleware.RequireAllRoles([]string{"admin", "moderator"})(handler)
```

## Context Values

- `"user"` - Regular user authentication
- `"service_user"` - Service-to-service authentication 