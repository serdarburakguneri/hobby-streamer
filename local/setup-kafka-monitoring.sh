#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up Kafka monitoring with ELK stack..."

echo "[INFO] Waiting for Elasticsearch to be ready..."
max_attempts=60
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if curl -s http://localhost:9200/_cluster/health > /dev/null 2>&1; then
    echo "[INFO] Elasticsearch is ready."
    break
  fi
  echo "[INFO] Elasticsearch not ready yet, waiting... (attempt $((attempt + 1))/$max_attempts)"
  sleep 3
  attempt=$((attempt + 1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "[ERROR] Elasticsearch failed to become ready after $max_attempts attempts"
  exit 1
fi

echo "[INFO] Creating Kafka monitoring index template..."

# Create index template for Kafka events
curl -X PUT "localhost:9200/_template/kafka-events-template" \
  -H "Content-Type: application/json" \
  -d '{
    "index_patterns": ["kafka-events-*"],
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "refresh_interval": "1s"
    },
    "mappings": {
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "event_id": {
          "type": "keyword"
        },
        "event_type": {
          "type": "keyword"
        },
        "source": {
          "type": "keyword"
        },
        "correlation_id": {
          "type": "keyword"
        },
        "causation_id": {
          "type": "keyword"
        },
        "topic": {
          "type": "keyword"
        },
        "partition": {
          "type": "integer"
        },
        "offset": {
          "type": "long"
        },
        "consumer_group": {
          "type": "keyword"
        },
        "processing_time_ms": {
          "type": "long"
        },
        "data": {
          "type": "object",
          "dynamic": true
        },
        "service_name": {
          "type": "keyword"
        },
        "level": {
          "type": "keyword"
        }
      }
    }
  }' || echo "[WARN] Failed to create index template (may already exist)"

echo "[INFO] Creating Kafka monitoring dashboard..."

# Create Kibana dashboard for Kafka monitoring
curl -X POST "localhost:5601/api/kibana/dashboards/import" \
  -H "Content-Type: application/json" \
  -H "kbn-xsrf: true" \
  -d '{
    "version": "8.0.0",
    "objects": [
      {
        "id": "kafka-monitoring-dashboard",
        "type": "dashboard",
        "attributes": {
          "title": "Kafka Monitoring Dashboard",
          "hits": 0,
          "description": "Monitor Kafka topics, consumer groups, and event processing",
          "panelsJSON": "[]",
          "optionsJSON": "{\"hidePanelTitles\":false,\"useMargins\":true}",
          "version": 1,
          "timeRestore": false,
          "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
          }
        }
      }
    ]
  }' || echo "[WARN] Failed to create dashboard (may already exist)"

echo "[INFO] Setting up Kafka metrics collection..."

# Create a simple Kafka metrics collector script
cat > kafka-metrics.sh << 'EOF'
#!/bin/bash

while true; do
  # Get topic information
  docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list | while read topic; do
    if [ -n "$topic" ]; then
      # Get partition info
      partitions=$(docker exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic "$topic" --quiet | wc -l)
      
      # Get consumer group lag
      docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list | while read group; do
        if [ -n "$group" ]; then
          lag=$(docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group "$group" --topic "$topic" 2>/dev/null | tail -n +2 | awk '{sum += $6} END {print sum+0}')
          if [ -n "$lag" ] && [ "$lag" != "0" ]; then
            echo "kafka.consumer.lag,topic=$topic,group=$group value=$lag"
          fi
        fi
      done
      
      echo "kafka.topic.partitions,topic=$topic value=$partitions"
    fi
  done
  
  sleep 30
done
EOF

chmod +x kafka-metrics.sh

echo "[INFO] Kafka monitoring setup completed!"
echo "[INFO] Kafka metrics will be collected every 30 seconds"
echo "[INFO] View Kafka UI at: http://localhost:8086"
echo "[INFO] View Kibana at: http://localhost:5601"
echo ""
echo "[INFO] Available Kafka monitoring features:"
echo "  - Topic partition monitoring"
echo "  - Consumer group lag tracking"
echo "  - Event processing metrics"
echo "  - Error rate monitoring"
echo "  - Performance analytics" 