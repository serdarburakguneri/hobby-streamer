#!/bin/bash
set -e

cd "$(dirname "$0")"

source ./setup-environment.sh

echo "[INFO] Starting backend services (auth-service, asset-manager, transcoder, nginx)..."
docker-compose up -d auth-service asset-manager transcoder nginx

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

echo "[INFO] All services are up to date and running."
echo "[INFO] - Auth Service: http://localhost:8080"
echo "[INFO] - Asset Manager GraphQL: http://localhost:8082/query"
echo "[INFO] - Nginx (HLS/DASH Proxy): http://localhost:8083"
echo "[INFO] - Neo4j Browser: http://localhost:7474"
echo "[INFO] - Keycloak: http://localhost:9090"
echo "[INFO] - LocalStack: http://localhost:4566"
echo "[INFO] - Kibana (Logs): http://localhost:5601"
echo "[INFO] - Elasticsearch: http://localhost:9200"

echo "[INFO] Backend services setup completed" 