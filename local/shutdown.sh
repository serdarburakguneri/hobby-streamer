#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "[INFO] Starting Hobby Streamer shutdown process..."

echo "[INFO] Phase 1: Stopping frontend services..."
pkill -f 'npm run web' || true
pkill -f 'expo start' || true

echo "[INFO] Phase 2: Stopping nginx proxy..."
docker-compose -f ../docker-compose.yml down nginx || true

echo "[INFO] Phase 3: Stopping backend services..."
docker-compose -f ../docker-compose.yml down streaming-api asset-manager auth-service transcoder || true

echo "[INFO] Phase 4: Stopping Redis..."
docker-compose -f ../docker-compose.yml down redis || true

echo "[INFO] Phase 5: Stopping infrastructure services..."
docker-compose -f ../docker-compose.yml down keycloak neo4j kibana elasticsearch fluentd || true

echo "[INFO] Phase 6: Stopping all remaining containers..."
docker-compose -f ../docker-compose.yml down || true

echo "[INFO] Phase 7: Cleaning up Docker networks..."
docker network prune -f || true

echo ""
echo "[INFO] Shutdown completed successfully!"
echo "[INFO] All services have been stopped."
echo "[INFO] To start all services again: ./build.sh" 