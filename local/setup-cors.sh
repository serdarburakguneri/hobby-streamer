#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up CORS configuration for all services..."

echo "[INFO] Applying CORS configuration to S3 buckets..."
for bucket in content-east content-west; do
  echo "[INFO] Applying CORS to $bucket"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors \
    --bucket $bucket \
    --cors-configuration file://cors.json
done

echo "[INFO] Redeploying API Gateway to ensure CORS headers are active..."
if [ -f ".api-gateway-id" ]; then
  API_ID=$(cat .api-gateway-id)
  if [ ! -z "$API_ID" ]; then
    echo "[INFO] Redeploying API Gateway with ID: $API_ID"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-deployment \
      --rest-api-id $API_ID \
      --stage-name dev > /dev/null
    echo "[INFO] API Gateway redeployed successfully"
  else
    echo "[WARN] API Gateway ID is empty, skipping redeployment"
  fi
else
  echo "[WARN] API Gateway ID file not found, skipping redeployment"
fi

echo "[INFO] CORS configuration completed successfully!"
echo "[INFO] S3 buckets: CORS enabled for localhost:8081"
echo "[INFO] API Gateway: CORS headers configured and deployed" 