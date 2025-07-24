# Config Package

Central config management for backend services.

## Features
Environment-based config, dynamic loading, secrets management, type-safe access, hot reload.

## Usage
```go
configManager, err := config.NewManager("service-name")
cfg := configManager.GetConfig()
dynamicCfg := configManager.GetDynamicConfig()
```

## Example Structure
```yaml
environment: development
service: my-service
log: { level: debug, format: text }
server: { port: "8080" }
features: { enable_circuit_breaker: true, enable_retry: true }
components: { sqs: { job_queue_url: "http://localstack:4566/000000000000/job-queue" }, s3: { content_bucket: "content-east" }, neo4j: { uri: "bolt://neo4j:7687", username: "neo4j" } }
```

## Dynamic Config
```go
queueURL := dynamicCfg.GetStringFromComponent("sqs", "job_queue_url")
```

## Secrets
```go
secretsManager := config.NewSecretsManager()
secretsManager.LoadFromEnvironment()
password := secretsManager.Get("neo4j_password")
```

## Env Vars
CONFIG_ENV, CONFIG_PATH, CONFIG_RELOAD_INTERVAL

## Best Practices
Use env-specific files, keep secrets in env vars, validate config on startup, use type-safe accessors.