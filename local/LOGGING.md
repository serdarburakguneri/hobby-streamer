# Logging Setup

This project includes a centralized logging system using Fluentd, Elasticsearch, and Kibana to collect and visualize logs from all Docker containers.

## Architecture

- **Fluentd**: Log collector that receives logs from Docker containers and forwards them to Elasticsearch
- **Elasticsearch**: Search and analytics engine that stores all logs
- **Kibana**: Web interface for searching, analyzing, and visualizing logs

## Distributed Tracing

The logging system includes distributed tracing capabilities using tracking IDs that flow through all services:

### How It Works

1. **Request Initiation**: When a request hits any service, a unique tracking ID is generated (if not present)
2. **Header Propagation**: The tracking ID is passed via `X-Tracking-ID` header to all downstream services
3. **Context Injection**: Each service includes the tracking ID in all log entries
4. **Cross-Service Tracing**: You can trace a single operation across your entire system


### Benefits

- **End-to-End Visibility**: See the complete journey of any operation
- **Error Correlation**: Quickly identify which service caused an issue
- **Performance Analysis**: Track timing across service boundaries
- **Debugging**: Follow user actions through the entire system

## Services

### Logging Stack
- **Elasticsearch**: http://localhost:9200
- **Kibana**: http://localhost:5601
- **Fluentd**: Collects logs on port 24224

### Application Services (with logging enabled)
- Auth Service
- Asset Manager
- Transcoder
- Keycloak
- Neo4j
- LocalStack

## Usage

### Viewing Logs

1. **Kibana UI** (Recommended):
   - Open http://localhost:5601
   - Go to "Discover" to search and filter logs
   - Use "Dashboard" for predefined views

2. **Command Line**:
   ```bash
   # View logs for specific service
   ./local/logs.sh auth
   ./local/logs.sh asset
   ./local/logs.sh transcoder
   
   # View all logs
   ./local/logs.sh all
   ```

3. **Docker Compose**:
   ```bash
   # View logs for specific service
   docker-compose logs -f auth-service
   
   # View all logs
   docker-compose logs -f
   ```

### Log Structure

Each log entry includes your structured logger fields:
- `service_name`: Name of the service (auth-service, asset-manager, etc.)
- `timestamp`: When the log was generated
- `level`: Log level (debug, info, warn, error)
- `msg`: The log message
- `error`: Error details (when present)
- `method`: HTTP method (for request logs)
- `path`: HTTP path (for request logs)
- `status_code`: HTTP status code (for request logs)
- `duration_ms`: Request duration in milliseconds
- `request_id`: Request ID for tracing
- `tracking_id`: **Distributed tracing ID** - follows operations across all services
- `user_id`: User ID (when authenticated)
- `username`: Username (when authenticated)
- `remote_addr`: Client IP address
- `user_agent`: Client user agent
- `asset_id`: Asset ID (for asset-related operations)
- `message_type`: SQS message type (for transcoding jobs)
- `output_bucket`: Output S3 bucket name
- `output_key`: Output S3 key
- `bucket`: S3 bucket name
- `key`: S3 key
- `file_count`: Number of files processed
- `deleted_count`: Number of files deleted
- `error_count`: Number of errors encountered
- `expires_in_minutes`: URL expiration time

### Searching in Kibana

Common search queries for your logger format:
- `service_name:auth-service` - Filter by service
- `level:error` - Find all error logs
- `level:warn` - Find all warning logs
- `error:*` - Find logs with error details
- `method:POST` - Find POST requests
- `status_code:500` - Find 5xx errors
- `status_code:[400 TO 499]` - Find 4xx errors
- `duration_ms:>1000` - Find slow requests (>1s)
- `request_id:*` - Find logs with request IDs
- `tracking_id:*` - **Find all logs for a specific operation**
- `user_id:*` - Find authenticated user logs
- `msg:*error*` - Find messages containing "error"
- `timestamp:[now-1h TO now]` - Last hour's logs

### Advanced Search Examples

**Find all errors from auth service in the last 30 minutes:**
```
service_name:auth-service AND level:error AND timestamp:[now-30m TO now]
```

**Find slow requests (>2 seconds):**
```
duration_ms:>2000 AND method:*
```

**Find all 5xx errors with error details:**
```
status_code:[500 TO 599] AND error:*
```

**Find requests from a specific user:**
```
user_id:12345 AND method:*
```

**Find GraphQL queries:**
```
path:*graphql* AND method:POST
```

**Find authentication failures:**
```
msg:*auth* AND level:error
```

**Find lambda function logs:**
```
service_name:generate-presigned-url OR service_name:trigger-transcode-job OR service_name:delete-files
```

**Find transcoding job triggers:**
```
message_type:*transcode* AND level:info
```

**Find file deletion operations:**
```
service_name:delete-files AND deleted_count:>0
```

**Find S3 operations:**
```
bucket:* OR output_bucket:*
```

**Trace a complete operation across all services:**
```
tracking_id:a1b2c3d4e5f6
```

**Find all operations for a specific asset:**
```
asset_id:12345 AND tracking_id:*
```

### Setup

The logging stack is automatically started with the main build script:
```bash
./local/build.sh
```

To manually setup Kibana dashboard:
```bash
./local/setup-kibana-dashboard.sh
```

## Configuration

### Fluentd Configuration
Located at `local/fluentd/fluent.conf`:
- Collects logs from Docker containers via forward protocol
- Parses JSON logs
- Adds service name and timestamp
- Forwards to Elasticsearch with buffering

### Docker Logging
Each service in `docker-compose.yml` is configured with:
```yaml
logging:
  driver: "fluentd"
  options:
    fluentd-address: "localhost:24224"
    tag: "docker.service-name"
```

