# Dynamic Configuration System

Flexible configuration loader for services needing both static and dynamic config structures. Each service can define custom configuration blocks without tight coupling.

## Key Features

Dynamic components, typed core config, environment support, secrets support, feature flags, hot reloading, flexible access API.

## Base Configuration (Static & Typed)

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

Custom configuration blocks under `components`:

```yaml
components:
  aws:
    region: "us-east-1"
    endpoint: "http://localstack:4566"
    force_path_style: true

  sqs:
    transcoder_queue_url: "http://localstack:4566/000000000000/transcoder-jobs"
    analyze_queue_url: "http://localstack:4566/000000000000/analyze-completed"
    batch_size: 10
    visibility_timeout: 30

  database:
    type: "neo4j"
    uri: "bolt://neo4j:7687"
    username: "neo4j"
    password: "password"
    max_connections: 10
    connection_timeout: "5s"
```

## Accessing Configuration

### Typed fields (core config):
```go
cfg.Server.Port             // → "8080"
cfg.Features.EnableCaching // → true
```

### Dynamic fields (components):
```go
region := cfg.Components.Get("aws").GetString("region")
batchSize := cfg.Components.Get("sqs").GetInt("batch_size")
apiConfig := cfg.Components.Get("external_api").AsMap()

// Unmarshal into custom struct
var dbConfig struct {
    URI      string `mapstructure:"uri"`
    Username string `mapstructure:"username"`
}
cfg.Components.Unmarshal("database", &dbConfig)
```

> ⚠️ Internal service configuration system — not intended as general-purpose config loader.