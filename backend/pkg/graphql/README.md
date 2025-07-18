# GraphQL Package

A shared GraphQL client library for consuming GraphQL APIs across services.

## Features

- **Service-to-service authentication** using Keycloak service tokens
- **Circuit breaker integration** for resilience
- **Common error handling** for GraphQL responses
- **Configurable timeouts** and HTTP client settings
- **Structured logging** with service context

## Usage

### Basic Client Setup

```go
import (
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/graphql"
    "github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
)

// Create service client for authentication
serviceClient := auth.NewServiceClient(keycloakURL, realm, clientID, clientSecret)

// Create GraphQL client
gqlClient := graphql.NewClient(serviceClient, graphql.ClientConfig{
    Timeout: 10 * time.Second,
})
```

### Execute Query

```go
var response struct {
    Data struct {
        Asset *Asset `json:"asset"`
    } `json:"data"`
    Errors []struct {
        Message string `json:"message"`
    } `json:"errors,omitempty"`
}

query := `
    query {
        asset(id: "123") {
            id
            title
            description
        }
    }
`

err := gqlClient.ExecuteQuery(ctx, "http://asset-manager:8082/graphql", query, &response)
if err != nil {
    return err
}

if err := gqlClient.HandleGraphQLErrors(&response); err != nil {
    return err
}

asset := response.Data.Asset
```

### With Circuit Breaker

```go
circuitBreaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
    Name:      "asset-manager",
    Threshold: 5,
    Timeout:   30 * time.Second,
})

err := gqlClient.ExecuteQueryWithCircuitBreaker(ctx, circuitBreaker, url, query, &response)
```

## Interface

### ClientInterface

```go
type ClientInterface interface {
    ExecuteQuery(ctx context.Context, url, query string, response interface{}) error
    ExecuteQueryWithCircuitBreaker(ctx context.Context, circuitBreaker *errors.CircuitBreaker, url, query string, response interface{}) error
    HandleGraphQLErrors(response interface{}) error
}
```

### Configuration

```go
type ClientConfig struct {
    Timeout time.Duration
}
```

## Error Handling

The client automatically handles:
- HTTP request errors
- Authentication token retrieval
- GraphQL response parsing
- Circuit breaker integration

GraphQL errors are returned as `ExternalError` with the error messages from the GraphQL response. 