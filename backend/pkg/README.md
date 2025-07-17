# Shared Libraries

Shared libraries used across all backend services. Ensures consistency, reduces code duplication, and provides type-safe interfaces for common functionality.

## Architecture

The project uses a shared library approach to ensure consistency and reduce code duplication across services. All shared libraries are located in `backend/pkg/` and provide:

### Configuration & Infrastructure Libraries
- **Config**: Dynamic configuration management with service-specific components, feature flags, secrets management, and hot reloading
- **Auth**: JWT validation, role-based access control, and Keycloak integration
- **Logger**: Structured logging with consistent formatting and log levels
- **Constants**: Shared constants for HTTP status codes, user roles, and other common values
- **Errors**: Comprehensive error handling with typed errors, retry mechanisms, circuit breakers, and graceful degradation

### AWS Service Libraries
- **S3**: File storage operations with LocalStack support for local development
- **SQS**: Message queue operations with producer/consumer patterns and registry management

### Communication Libraries
- **Messages**: Type-safe SQS message payloads and constants for inter-service communication

### Usage Pattern
All backend services import these libraries as local modules using Go's replace directive:

```go
require (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/config v0.0.0
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
    // ... other shared libraries
)

replace (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/config => ../pkg/config
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../pkg/errors
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
    // ... other replace directives
)
```

## Available Libraries

### [Config Package](config/README.md)
Dynamic configuration management system with service-specific components, feature flags, secrets management, and hot reloading capabilities. Provides maximum flexibility for service configuration while maintaining type safety for core settings.

### [Auth Package](auth/README.md)
Shared authentication library with JWT validation and role-based authorization.

### [Constants Package](constants/README.md)
Common constants for HTTP status codes, roles, and other shared values.

### [Errors Package](errors/README.md)
Comprehensive error handling library with typed errors, retry mechanisms, circuit breakers, and graceful degradation patterns.

### [Logger Package](logger/README.md)
Centralized structured logging solution for all backend services.

### [Messages Package](messages/README.md)
Common SQS message payload structures and type constants for inter-service communication.

### [S3 Package](s3/README.md)
S3 client library for file upload, download, and directory operations with LocalStack support.

### [SQS Package](sqs/README.md)
AWS SQS client library with producer, consumer, and consumer registry functionality.

