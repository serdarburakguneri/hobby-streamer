# Constants Package

Shared constants used across backend services. Includes standard HTTP status codes, user roles, and common application-level values.

## Features

- HTTP status codes with descriptive names
- User roles for access control
- Application-level constants (pagination, upload limits, etc.)

## Usage

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

// HTTP status codes
w.WriteHeader(constants.StatusOK)
w.WriteHeader(constants.StatusCreated)

// Role check
if user.HasRole(constants.RoleAdmin) {
    // Admin-specific logic
}

// Common usage
if file.Size > constants.MaxUploadSize {
    // Reject upload
}
```

## Available Constants

### HTTP Status Codes

| Name                        | Value |
|----------------------------|-------|
| `StatusOK`                 | 200   |
| `StatusCreated`            | 201   |
| `StatusBadRequest`         | 400   |
| `StatusUnauthorized`       | 401   |
| `StatusForbidden`          | 403   |
| `StatusNotFound`           | 404   |
| `StatusInternalServerError`| 500   |

### User Roles

| Constant      | Description          |
|---------------|----------------------|
| `RoleAdmin`   | Administrator role   |
| `RoleUser`    | Standard user role   |
| `RoleEditor`  | Content editor role  |

### Application Constants

| Constant            | Description                            |
|---------------------|----------------------------------------|
| `MaxUploadSize`     | Maximum allowed file size (in bytes)   |
| `DefaultPageSize`   | Default number of items per page       |