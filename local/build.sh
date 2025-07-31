#!/bin/bash
set -e

cd "$(dirname "$0")"

chmod +x setup-*

echo "[INFO] Starting Hobby Streamer build process..."

echo "[INFO] Setting up environment..."
./setup-environment.sh

echo "[INFO] Setting up infrastructure..."
./setup-infrastructure.sh

echo "[INFO] Setting up Kafka infrastructure..."
./setup-kafka.sh

echo "[INFO] Setting up Kafka monitoring..."
./setup-kafka-monitoring.sh

echo "[INFO] Setting up S3 buckets..."
./setup-s3-buckets.sh

echo "[INFO] Setting up Lambda functions..."
./setup-lambdas.sh

echo "[INFO] Setting up API Gateway..."
./setup-api-gateway.sh

echo "[INFO] Starting Redis..."
./setup-redis.sh

echo "[INFO] Setting up kibana dashboard..."
./setup-kibana-dashboard.sh

echo "[INFO] Setting up backend services..."
./setup-backend-services.sh

echo "[INFO] Setting up frontend..."
./setup-frontend.sh

echo ""
echo "[INFO] Build completed successfully!"
echo "[INFO] All services are running and ready for development."
echo "[INFO] To stop all services: docker-compose down"
echo "[INFO] To stop CMS UI: pkill -f 'npm run web'"
