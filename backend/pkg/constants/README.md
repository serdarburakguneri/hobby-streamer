# Constants Package

Common constants for HTTP status codes, roles, and other shared values used across services.

## Features

- HTTP status codes and messages
- User roles and permissions
- Common application constants

## Usage

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

// HTTP status codes
w.WriteHeader(constants.StatusOK)
w.WriteHeader(constants.StatusCreated)
w.WriteHeader(constants.StatusBadRequest)

// User roles
if user.HasRole(constants.RoleAdmin) {
    // Admin functionality
}

// Common constants
const maxFileSize = constants.MaxUploadSize
```

## Available Constants

### HTTP Status Codes
- `StatusOK`: 200
- `StatusCreated`: 201
- `StatusBadRequest`: 400
- `StatusUnauthorized`: 401
- `StatusForbidden`: 403
- `StatusNotFound`: 404
- `StatusInternalServerError`: 500

### User Roles
- `RoleAdmin`: Administrator role
- `RoleUser`: Standard user role
- `RoleEditor`: Content editor role

### Application Constants
- `MaxUploadSize`: Maximum file upload size in bytes
- `DefaultPageSize`: Default pagination page size 