# Dynamic Configuration System

A flexible configuration loader that supports both static and dynamic service configurations. Allows each service to define its own custom structure without requiring hardcoded types.

## Key Features

- **Dynamic Components** – Each service defines its own configuration blocks under `components`
- **Typed Core Config** – Core fields like server, logging, and feature flags are fully typed
- **Environment Support** – Separate configurations for development, staging, and production
- **Secrets Management** – Built-in support for secure values
- **Feature Flags** – Toggle features at runtime without code changes
- **Hot Reloading** – Reload configuration without restarting the service
- **Flexible Access** – Read dynamic config values as strings, ints, bools, maps, etc.

## Base Configuration (Typed)

The common fields shared across all services are defined with strict types:

```yaml
environment: development
service: my-service

log:
  level: debug
  format: text

server:
  port: "8080"
  read_timeout: "15s"
  write_timeout: "15s"
  idle_timeout: "60s"

features:
  enable_circuit_breaker: true
  enable_retry: true
  enable_caching: true
  enable_metrics: false
  enable_tracing: false

circuit_breaker:
  threshold: 5
  timeout: "30s"

retry:
  max_attempts: 3
  base_delay: "100ms"
  max_delay: "5s"

cache:
  ttl: "30m"
```

## Dynamic Components

Each service can define custom sections under the `components` field:

```yaml
components:
  aws:
    region: "us-east-1"
    endpoint: "http://localstack:4566"
    force_path_style: true
    custom_setting: "value"

  sqs:
    transcoder_queue_url: "http://localstack:4566/000000000000/transcoder-jobs"
    analyze_queue_url: "http://localstack:4566/000000000000/analyze-completed"
    custom_event_queue: "http://localstack:4566/000000000000/custom-events"
    batch_size: 10
    visibility_timeout: 30

  database:
    type: "neo4j"
    uri: "bolt://neo4j:7687"
    username: "neo4j"
    password: "password"
    max_connections: 10
    connection_timeout: "5s"

  my_custom_service:
    url: "http://my-service:8080"
    timeout: "30s"
    retry_count: 3
    enabled: true
    features:
      - "feature1"
      - "feature2"

  external_api:
    base_url: "https://api.example.com"
    api_key: "your-api-key"
    rate_limit: 100
    timeout: "10s"
    retry_on_failure: true
```

## Accessing Configuration in Code

Static configuration is mapped to typed structs. Dynamic components can be accessed using helper functions:

```go
// Get a value as string
cfg.Components.Get("aws").GetString("region")

// Get a nested value as int
cfg.Components.Get("sqs").GetInt("batch_size")

// Get a map
cfg.Components.Get("external_api").AsMap()

// Unmarshal into custom struct
var dbConfig struct {
    URI      string `mapstructure:"uri"`
    Username string `mapstructure:"username"`
}
cfg.Components.Unmarshal("database", &dbConfig)
```