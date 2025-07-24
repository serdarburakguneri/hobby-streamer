# Constants Package

Shared constants for backend services: HTTP status codes, user roles, app-level values.

## Features
HTTP status codes, user roles, pagination, upload limits, common values.

## Usage
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
w.WriteHeader(constants.StatusOK)
if user.HasRole(constants.RoleAdmin) { /* ... */ }
if file.Size > constants.MaxUploadSize { /* ... */ }
```

## Available Constants
Status: StatusOK (200), StatusCreated (201), StatusBadRequest (400), StatusUnauthorized (401), StatusForbidden (403), StatusNotFound (404), StatusInternalServerError (500). Roles: RoleAdmin, RoleUser, RoleEditor. App: MaxUploadSize, DefaultPageSize.