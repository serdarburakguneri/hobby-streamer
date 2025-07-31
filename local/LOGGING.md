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

1. **Kibana UI**:
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


