#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Starting backend services (auth-service, asset-manager, transcoder, streaming-api, nginx)..."
docker-compose build auth-service asset-manager transcoder streaming-api nginx
docker-compose up -d auth-service asset-manager transcoder streaming-api nginx

echo "[INFO] All services are running successfully!"
echo "[INFO] You can view logs using: docker-compose logs <service-name>"
echo "[INFO] Or access Kibana at http://localhost:5601"

echo "[INFO] Waiting for services to be ready..."
sleep 10

echo "[INFO] Testing Asset Manager GraphQL endpoint..."
until curl -s -X POST http://localhost:8082/graphql -H "Content-Type: application/json" -d '{"query":"{ __schema { types { name } } }"}' > /dev/null 2>&1; do
  echo "[INFO] Waiting for Asset Manager GraphQL endpoint to be ready..."
  sleep 3
done
echo "[INFO] Asset Manager GraphQL endpoint is ready."

echo "[INFO] Testing Nginx CDN endpoint..."
until curl -s http://localhost:8083/health > /dev/null 2>&1; do
  echo "[INFO] Waiting for Nginx CDN endpoint to be ready..."
  sleep 3
done
echo "[INFO] Nginx CDN endpoint is ready."