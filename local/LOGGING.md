# Logging

Fluentd → Elasticsearch → Kibana stack collects logs from all containers.

## Access
Kibana: `http://localhost:5601`. Elasticsearch: `http://localhost:9200`.

## CLI
```bash
./local/logs.sh auth|asset|transcoder|all
```
