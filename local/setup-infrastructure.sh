#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Stopping all running containers for fresh start..."
docker-compose down

echo "[INFO] Starting Elasticsearch and Kibana..."
docker-compose up -d elasticsearch kibana

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

echo "[INFO] Setting up Kibana dashboard..."
./setup-kibana-dashboard.sh

echo "[INFO] Starting Fluentd..."
docker-compose up -d fluentd

echo "[INFO] Waiting for Fluentd to be healthy..."
max_attempts=60
attempt=0
while [ $attempt -lt $max_attempts ]; do
  container_id=$(docker-compose ps -q fluentd)
  if [ -n "$container_id" ]; then
    health_status=$(docker inspect --format='{{.State.Health.Status}}' "$container_id" 2>/dev/null || echo "unknown")
    if [ "$health_status" == "healthy" ]; then
      echo "[INFO] Fluentd is healthy."
      break
    fi
  fi
  echo "[INFO] Fluentd not healthy yet, waiting... (attempt $((attempt + 1))/$max_attempts)"
  sleep 2
  attempt=$((attempt + 1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "[ERROR] Fluentd failed to become healthy after $max_attempts attempts"
  exit 1
fi

echo "[INFO] Starting infrastructure services (LocalStack, Neo4j, Keycloak, Kafka)..."

if [ ! -f "keycloak-certs/cert.pem" ] || [ ! -f "keycloak-certs/key.pem" ]; then
  echo "[INFO] Generating Keycloak HTTPS certificates..."
  ./generate-keycloak-certs.sh
fi

docker-compose up -d localstack neo4j keycloak zookeeper kafka akhq

echo "[INFO] Waiting for Keycloak to be ready..."
max_attempts=60
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if curl -s http://localhost:9090/realms/master > /dev/null 2>&1; then
    echo "[INFO] Keycloak master realm is up."
    break
  fi
  echo "[INFO] Keycloak not ready yet, waiting... (attempt $((attempt + 1))/$max_attempts)"
  sleep 3
  attempt=$((attempt + 1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "[ERROR] Keycloak failed to become ready after $max_attempts attempts"
  exit 1
fi

sleep 10

if ! curl -s http://localhost:9090/realms/hobby | grep -q '"realm":"hobby"'; then
  echo "[INFO] Importing Keycloak hobby realm..."
  docker exec hobby-streamer-keycloak-1 /opt/keycloak/bin/kc.sh import --file=/opt/keycloak/data/import/hobby-realm.json --override=true
  if curl -s http://localhost:9090/realms/hobby | grep -q '"realm":"hobby"'; then
    echo "[INFO] Hobby realm imported successfully."
  else
    echo "[ERROR] Failed to import hobby realm!"
    exit 1
  fi
else
  echo "[INFO] Keycloak hobby realm already exists."
fi

echo "[INFO] Infrastructure setup completed"

echo "[INFO] Setting up Kafka topics..."
./setup-sqs-queues.sh

echo "[INFO] Setting up AWS resources..."
./setup-aws-resources.sh

echo "[INFO] Setting up Lambda functions..."
./setup-lambdas.sh

echo "[INFO] All infrastructure setup completed successfully!" 