# Config Package

A centralized configuration management package for the Hobby Streamer backend services.

## Features

- Environment-based configuration
- Dynamic configuration loading
- Secrets management
- Type-safe configuration access
- Hot reloading support

## Usage

### Basic Configuration

```go
configManager, err := config.NewManager("service-name")
if err != nil {
    log.Fatal(err)
}
defer configManager.Close()

cfg := configManager.GetConfig()
dynamicCfg := configManager.GetDynamicConfig()
```

### Configuration Structure

```yaml
environment: development
service: my-service

log:
  level: debug
  format: text
  async:
    enabled: true
    buffer_size: 1000

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

components:
  sqs:
    job_queue_url: "http://localstack:4566/000000000000/job-queue"
    completion_queue_url: "http://localstack:4566/000000000000/completion-queue"
  
  s3:
    content_bucket: "content-east"
  
  neo4j:
    uri: "bolt://neo4j:7687"
    username: "neo4j"
    max_connections: 50
    connection_timeout: "30s"
    max_lifetime: "1h"
```

### Dynamic Configuration

Dynamic configuration allows runtime updates without service restart:

```go
dynamicCfg := configManager.GetDynamicConfig()

// Get string values
queueURL := dynamicCfg.GetStringFromComponent("sqs", "job_queue_url")

// Get int values
maxConnections := dynamicCfg.GetIntFromComponent("neo4j", "max_connections")

// Get bool values
enabled := dynamicCfg.GetBoolFromComponent("features", "enable_circuit_breaker")
```

### Secrets Management

```go
secretsManager := config.NewSecretsManager()
secretsManager.LoadFromEnvironment()

password := secretsManager.Get("neo4j_password")
```

## Configuration Files

Configuration files should be placed in the `config/` directory of each service:

- `config.development.yaml` - Development environment
- `config.production.yaml` - Production environment (if needed)

## Environment Variables

The following environment variables can be used to override configuration:

- `CONFIG_ENV` - Environment name (default: development)
- `CONFIG_PATH` - Path to config files (default: ./config)
- `CONFIG_RELOAD_INTERVAL` - Hot reload interval (default: 30s)

## Components

### SQS Configuration

The SQS component uses a simplified two-queue architecture:

- `job_queue_url` - Queue for job triggers (analyze, transcode)
- `completion_queue_url` - Queue for job completions

### Neo4j Configuration

- `uri` - Neo4j connection URI
- `username` - Database username
- `max_connections` - Connection pool size
- `connection_timeout` - Connection timeout
- `max_lifetime` - Connection lifetime

### S3 Configuration

- `content_bucket` - S3 bucket for content storage

## Best Practices

1. Use environment-specific configuration files
2. Keep sensitive data in environment variables
3. Use dynamic configuration for frequently changing values
4. Validate configuration on startup
5. Use type-safe accessors when possible