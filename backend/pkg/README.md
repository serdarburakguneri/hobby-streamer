# Shared Libraries

This directory contains shared libraries used across all backend services in the hobby-streamer project. These libraries ensure consistency, reduce code duplication, and provide type-safe interfaces for common functionality.

## Architecture

The project uses a shared library approach to ensure consistency and reduce code duplication across services. All shared libraries are located in `backend/pkg/` and provide:

### Core Infrastructure Libraries
- **Auth**: JWT validation, role-based access control, and Keycloak integration
- **Logger**: Structured logging with consistent formatting and log levels
- **Constants**: Shared constants for HTTP status codes, user roles, and other common values

### AWS Service Libraries
- **S3**: File storage operations with LocalStack support for local development
- **SQS**: Message queue operations with producer/consumer patterns and registry management

### Communication Libraries
- **Messages**: Type-safe SQS message payloads and constants for inter-service communication

### Benefits
- **Type Safety**: Compile-time checking prevents runtime errors
- **Consistency**: Shared interfaces ensure all services behave the same way
- **Maintainability**: Changes to common functionality only need to be made once
- **Documentation**: Each library has comprehensive documentation and examples

### Usage Pattern
All backend services import these libraries as local modules using Go's replace directive:

```go
require (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
    // ... other shared libraries
)

replace (
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
    github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
    // ... other replace directives
)
```

## Available Libraries

### [Auth Package](auth/README.md)
Shared authentication library with JWT validation and role-based authorization.

### [Constants Package](constants/README.md)
Common constants for HTTP status codes, roles, and other shared values.

### [Logger Package](logger/README.md)
Centralized structured logging solution for all backend services.

### [Messages Package](messages/README.md)
Common SQS message payload structures and type constants for inter-service communication.

### [S3 Package](s3/README.md)
S3 client library for file upload, download, and directory operations with LocalStack support.

### [SQS Package](sqs/README.md)
AWS SQS client library with producer, consumer, and consumer registry functionality.

## Development Guidelines

### Adding New Libraries
1. Create a new directory under `backend/pkg/`
2. Include a `go.mod` file with the module name
3. Add comprehensive documentation in a `README.md` file
4. Update this main README to include the new library
5. Update all services that need the library to import it

### Library Standards
- Each library should be self-contained with minimal dependencies
- Provide clear interfaces and examples
- Include comprehensive error handling
- Support both local development (LocalStack) and production environments
- Follow Go best practices and conventions 