# Dynamic Configuration System

Configuration management system that allows each service to define its own components dynamically without hardcoded types.

## Key Features

- **Dynamic Components**: No hardcoded service-specific types - each service defines its own component structure
- **Type Safety for Core Config**: Common types (LogConfig, ServerConfig, etc.) are strongly typed
- **Flexible Access**: Multiple ways to access component data (string, int, bool, float, map)
- **Environment Support**: Different configs for development, staging, production
- **Secrets Management**: Secure handling of sensitive data
- **Feature Flags**: Runtime feature toggles
- **Hot Reloading**: Configuration changes without restart

## Architecture

### Base Configuration (Statically Typed)
Common configuration that all services share:

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

### Dynamic Components
Services can define any structure under `components`:

```yaml
components:
  # AWS configuration - any keys you want
  aws:
    region: "us-east-1"
    endpoint: "http://localstack:4566"
    force_path_style: true
    custom_setting: "value"

  # SQS configuration - service-specific queue names
  sqs:
    transcoder_queue_url: "http://localstack:4566/000000000000/transcoder-jobs"
    analyze_queue_url: "http://localstack:4566/000000000000/analyze-completed"
    custom_event_queue: "http://localstack:4566/000000000000/custom-events"
    batch_size: 10
    visibility_timeout: 30

  # Database configuration - any database type
  database:
    type: "neo4j"
    uri: "bolt://neo4j:7687"
    username: "neo4j"
    password: "password"
    max_connections: 10
    connection_timeout: "5s"

  # Custom service configuration
  my_custom_service:
    url: "http://my-service:8080"
    timeout: "30s"
    retry_count: 3
    enabled: true
    features:
      - "feature1"
      - "feature2"

  # External API configuration
  external_api:
    base_url: "https://api.example.com"
    api_key: "your-api-key"
    rate_limit: 100
    timeout: "10s"
    retry_on_failure: true
```

