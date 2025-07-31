#!/bin/bash

cd "$(dirname "$0")"

echo "[INFO] Setting up Kibana dashboard..."

# Wait for Kibana to be ready
until curl -s http://localhost:5601/api/status > /dev/null 2>&1; do
  echo "[INFO] Waiting for Kibana to be ready..."
  sleep 5
done

# Create index pattern with default columns
# Note: '@timestamp' is the standard time field in ES/Kibana
#       'columns' sets the default selected fields in Discover

echo "[INFO] Creating index pattern..."
curl -X POST "http://localhost:5601/api/saved_objects/index-pattern/docker-logs" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "attributes": {
      "title": "docker-logs-*",
      "timeFieldName": "@timestamp",
      "columns": ["@timestamp", "service_name", "level", "msg", "method", "path", "status_code"]
    }
  }' 2>/dev/null || echo "[INFO] Index pattern already exists or failed to create"

# Create search for error logs
echo "[INFO] Creating error logs search..."
curl -X POST "http://localhost:5601/api/saved_objects/search/error-logs" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "attributes": {
      "title": "Error Logs",
      "description": "All error level logs",
      "hits": 0,
      "columns": ["@timestamp", "service_name", "level", "msg", "error"],
      "sort": [["@timestamp", "desc"]],
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"query\":{\"query\":\"level:error\",\"language\":\"kuery\"},\"filter\":[],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
      }
    },
    "references": [
      {
        "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
        "type": "index-pattern",
        "id": "docker-logs"
      }
    ]
  }' 2>/dev/null || echo "[INFO] Error logs search already exists or failed to create"

# Create search for HTTP requests
echo "[INFO] Creating HTTP requests search..."
curl -X POST "http://localhost:5601/api/saved_objects/search/http-requests" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "attributes": {
      "title": "HTTP Requests",
      "description": "All HTTP request logs",
      "hits": 0,
      "columns": ["@timestamp", "service_name", "method", "path", "status_code", "duration_ms"],
      "sort": [["@timestamp", "desc"]],
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"query\":{\"query\":\"method:*\",\"language\":\"kuery\"},\"filter\":[],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
      }
    },
    "references": [
      {
        "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
        "type": "index-pattern",
        "id": "docker-logs"
      }
    ]
  }' 2>/dev/null || echo "[INFO] HTTP requests search already exists or failed to create"

# Create basic dashboard
echo "[INFO] Creating basic dashboard..."
curl -X POST "http://localhost:5601/api/saved_objects/dashboard/hobby-streamer-logs" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d '{
    "attributes": {
      "title": "Hobby Streamer Logs",
      "hits": 0,
      "description": "Centralized logs from all services",
      "panelsJSON": "[{\"type\":\"search\",\"id\":\"error-logs\",\"panelIndex\":\"1\",\"gridData\":{\"x\":0,\"y\":0,\"w\":24,\"h\":10,\"i\":\"1\"}},{\"type\":\"search\",\"id\":\"http-requests\",\"panelIndex\":\"2\",\"gridData\":{\"x\":0,\"y\":10,\"w\":24,\"h\":10,\"i\":\"2\"}}]",
      "optionsJSON": "{\"hidePanelTitles\":false,\"useMargins\":true}",
      "version": 1,
      "timeRestore": false,
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
      }
    },
    "references": [
      {
        "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
        "type": "index-pattern",
        "id": "docker-logs"
      }
    ]
  }' 2>/dev/null || echo "[INFO] Dashboard already exists or failed to create"

echo "[INFO] Kibana dashboard setup complete!"
echo "[INFO] Access Kibana at: http://localhost:5601"
echo "[INFO] Go to Discover to view logs or Dashboard to see the overview" 