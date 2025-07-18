# Dynamic Configuration System

A flexible configuration loader designed for services that need both static and dynamic config structures. Each service can define its own custom configuration blocks without tight coupling or hardcoded types.

---

## Key Features

- **Dynamic components** – Add service-specific config under `components`
- **Typed core config** – Common fields like logging, timeouts, and feature flags are fully typed
- **Environment support** – Separate configs for development, staging, production
- **Secrets support** – Secure value placeholders supported in YAML
- **Feature flags** – Toggle behavior at runtime without redeploying
- **Hot reloading** – Supports live config reload without restarting
- **Flexible access API** – Read values as strings, ints, maps, or unmarshal into structs

---

## Base Configuration (Static & Typed)

These fields are available in all services by default:

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

---

## Dynamic Components

Custom configuration blocks live under the `components` field. These can be tailored to each service’s needs.

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

  external_api:
    base_url: "https://api.example.com"
    api_key: "your-api-key"
    rate_limit: 100
    timeout: "10s"
    retry_on_failure: true

  my_custom_service:
    url: "http://my-service:8080"
    timeout: "30s"
    retry_count: 3
    enabled: true
    features:
      - "feature1"
      - "feature2"
```

---

## Accessing Configuration in Code

### Typed fields (core config):

```go
cfg.Server.Port             // → "8080"
cfg.Features.EnableCaching // → true
```

### Dynamic fields (components):

```go
// Get a single string
region := cfg.Components.Get("aws").GetString("region")

// Get an int
batchSize := cfg.Components.Get("sqs").GetInt("batch_size")

// Get as map[string]interface{}
apiConfig := cfg.Components.Get("external_api").AsMap()

// Unmarshal into custom struct
var dbConfig struct {
    URI      string `mapstructure:"uri"`
    Username string `mapstructure:"username"`
}
cfg.Components.Unmarshal("database", &dbConfig)
```

---

> ⚠️ This system is designed for internal service configuration — not intended as a general-purpose config loader.