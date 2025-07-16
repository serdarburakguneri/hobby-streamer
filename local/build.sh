#!/bin/bash
set -e

cd "$(dirname "$0")"

chmod +x setup-*

echo "[INFO] Starting Hobby Streamer build process..."

echo "[INFO] Phase 1: Setting up environment..."
./setup-environment.sh

echo "[INFO] Phase 2: Setting up infrastructure..."
./setup-infrastructure.sh

echo "[INFO] Phase 3: Setting up AWS resources..."
./setup-aws-resources.sh

#echo "[INFO] Phase 4: Setting up CloudFront distributions..."
#./setup-cloudfront.sh

echo "[INFO] Phase 4: Setting up Lambda functions..."
./setup-lambdas.sh

echo "[INFO] Phase 5: Setting up API Gateway..."
./setup-api-gateway.sh

echo "[INFO] Phase 6: Starting Redis..."
./setup-redis.sh

echo "[INFO] Phase 7: Setting up backend services..."
./setup-backend-services.sh

echo "[INFO] Phase 8: Setting up nginx proxy..."
./setup-nginx.sh

echo "[INFO] Phase 9: Setting up frontend..."
./setup-frontend.sh

echo ""
echo "[INFO] Build completed successfully!"
echo "[INFO] All services are running and ready for development."
echo "[INFO] To stop all services: docker-compose down"
echo "[INFO] To stop CMS UI: pkill -f 'npm run web'"
