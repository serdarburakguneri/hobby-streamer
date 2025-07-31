#!/bin/bash

cd "$(dirname "$0")"

echo "[INFO] Demo: Distributed Tracing with Tracking IDs"
echo ""

# Generate a tracking ID
TRACKING_ID=$(openssl rand -hex 8)
echo "[INFO] Generated tracking ID: $TRACKING_ID"
echo ""

echo "[INFO] Making requests with tracking ID..."
echo ""

# Test auth service
echo "[INFO] 1. Testing Auth Service..."
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tracking-ID: $TRACKING_ID" \
  -d '{"username":"user","password":"user"}' > /dev/null
echo "   ✓ Auth service request sent"

# Test asset manager
echo "[INFO] 2. Testing Asset Manager..."
curl -s -X POST http://localhost:8082/graphql \
  -H "Content-Type: application/json" \
  -H "X-Tracking-ID: $TRACKING_ID" \
  -d '{"query":"{ __schema { types { name } } }"}' > /dev/null
echo "   ✓ Asset manager request sent"

# Test lambda (via API Gateway)
echo "[INFO] 3. Testing Lambda (Presigned URL)..."
API_ID=$(cat .api-gateway-id 2>/dev/null || echo "demo")
curl -s -X POST "http://localhost:4566/restapis/$API_ID/dev/_user_request_/upload" \
  -H "Content-Type: application/json" \
  -H "X-Tracking-ID: $TRACKING_ID" \
  -d '{"fileName":"demo-video.mp4"}' > /dev/null
echo "   ✓ Lambda request sent"

echo ""
echo "[INFO] Demo completed!"
echo ""
echo "[INFO] To see all logs for this operation:"
echo "   Open Kibana: http://localhost:5601"
echo "   Go to Discover"
echo "   Search for: tracking_id:$TRACKING_ID"
echo ""
echo "[INFO] This will show you:"
echo "   - All services that processed the request"
echo "   - Timing information for each step"
echo "   - Any errors that occurred"
echo "   - Complete request flow across your system"
echo ""
echo "[INFO] Tracking ID: $TRACKING_ID" 