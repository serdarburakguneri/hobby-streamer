# Shared Auth Package

A shared authentication package that provides token validation and role-based authorization for microservices.

## Overview

This package provides:
- **TokenValidator interface** for consistent token validation across services
- **KeycloakValidator implementation** for JWT validation with Keycloak
- **HTTP middleware** for easy integration into services
- **Role-based authorization** helpers

## Usage

### 1. Import the package

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
```

### 2. Create a validator

```go
validator := auth.NewKeycloakValidator(
    "http://localhost:8080",  // Keycloak URL
    "hobby",                  // Realm
    "asset-manager",          // Client ID
)
```

### 3. Create middleware

```go
middleware := auth.NewAuthMiddleware(validator)
```

### 4. Apply to routes

```go
// Require authentication only
r.HandleFunc("/assets", middleware.RequireAuth(handler.ListAssets)).Methods("GET")

// Require specific role
r.HandleFunc("/assets", middleware.RequireRole("admin")(handler.CreateAsset)).Methods("POST")

// Require any of multiple roles
r.HandleFunc("/assets", middleware.RequireAnyRole([]string{"admin", "editor"})(handler.UpdateAsset)).Methods("PUT")

// Require all roles
r.HandleFunc("/assets", middleware.RequireAllRoles([]string{"admin", "moderator"})(handler.DeleteAsset)).Methods("DELETE")
```

## Interface

### TokenValidator

```go
type TokenValidator interface {
    ValidateToken(ctx context.Context, token string) (*User, error)
    HasRole(user *User, role string) bool
    HasAnyRole(user *User, roles []string) bool
    HasAllRoles(user *User, roles []string) bool
}
```

### User

```go
type User struct {
    ID       string   `json:"id"`
    Username string   `json:"username"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles"`
}
```

## Middleware Functions

- `RequireAuth()` - Validates token and adds user to context
- `RequireRole(role)` - Requires specific role
- `RequireAnyRole(roles)` - Requires any of the specified roles
- `RequireAllRoles(roles)` - Requires all specified roles

## Example Integration

```go
func main() {
    validator := auth.NewKeycloakValidator(
        os.Getenv("KEYCLOAK_URL"),
        os.Getenv("KEYCLOAK_REALM"),
        os.Getenv("KEYCLOAK_CLIENT_ID"),
    )
    
    middleware := auth.NewAuthMiddleware(validator)
    
    r := mux.NewRouter()
    
    // Public endpoints
    r.HandleFunc("/health", healthHandler).Methods("GET")
    
    // Protected endpoints
    r.HandleFunc("/assets", middleware.RequireAuth(assetHandler.ListAssets)).Methods("GET")
    r.HandleFunc("/assets", middleware.RequireRole("admin")(assetHandler.CreateAsset)).Methods("POST")
    
    http.ListenAndServe(":8080", r)
}
``` 