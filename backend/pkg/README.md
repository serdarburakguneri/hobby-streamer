# Shared Libraries

A collection of shared packages used across backend services. These libraries provide common functionality, improve consistency, and reduce code duplication.

## Usage

All shared packages are located under `backend/pkg/` and imported via local module paths using Go’s `replace` directive:

```go
require (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/config v0.0.0
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
    // ...
)

replace (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/config => ../pkg/config
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
    // ...
)
```

## Shared Packages

### [config](config/README.md)

Loads environment-based configuration with support for service-specific settings, feature flags, secrets, and live reloading.

### [auth](auth/README.md)

JWT token validation, Keycloak integration, and role-based access control.

### [logger](logger/README.md)

Structured logging with consistent formatting and log level support across services.

### [constants](constants/README.md)

Common constants for HTTP status codes, roles, and other shared enums.

### [errors](errors/README.md)

Typed error definitions with built-in support for retries, circuit breakers, and fallback logic.

### [messages](messages/README.md)

SQS message definitions used for inter-service communication. Includes type-safe payload structures and enums.

### [s3](s3/README.md)

Helpers for uploading, downloading, and managing files in S3-compatible storage (with LocalStack support for local development).

### [sqs](sqs/README.md)

Producer and consumer utilities for SQS. Includes a consumer registry, retry logic, and message routing support.