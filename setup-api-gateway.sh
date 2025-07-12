#!/bin/bash
set -e

echo "[INFO] Setting up API Gateway for presigned URL Lambda..."

# Create REST API
echo "[INFO] Creating REST API..."
API_RESPONSE=$(aws --endpoint-url=http://localhost:4566 apigateway create-rest-api --name "presigned-url-api")
API_ID=$(echo $API_RESPONSE | jq -r '.id')
echo "[INFO] API ID: $API_ID"

# Get root resource ID
echo "[INFO] Getting root resource..."
ROOT_RESPONSE=$(aws --endpoint-url=http://localhost:4566 apigateway get-resources --rest-api-id $API_ID)
ROOT_ID=$(echo $ROOT_RESPONSE | jq -r '.items[] | select(.path == "/") | .id')
echo "[INFO] Root Resource ID: $ROOT_ID"

# Create upload resource
echo "[INFO] Creating upload resource..."
UPLOAD_RESPONSE=$(aws --endpoint-url=http://localhost:4566 apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part upload)
UPLOAD_ID=$(echo $UPLOAD_RESPONSE | jq -r '.id')
echo "[INFO] Upload Resource ID: $UPLOAD_ID"

# Create POST method
echo "[INFO] Creating POST method..."
aws --endpoint-url=http://localhost:4566 apigateway put-method \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method POST \
  --authorization-type "NONE"

# Integrate POST method with Lambda
echo "[INFO] Integrating POST method with Lambda..."
aws --endpoint-url=http://localhost:4566 apigateway put-integration \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:eu-west-1:lambda:path/2015-03-31/functions/arn:aws:lambda:eu-west-1:000000000000:function:generate-presigned-url/invocations

# Create OPTIONS method for CORS
echo "[INFO] Creating OPTIONS method for CORS..."
aws --endpoint-url=http://localhost:4566 apigateway put-method \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method OPTIONS \
  --authorization-type "NONE"

# Create mock integration for OPTIONS
echo "[INFO] Creating mock integration for OPTIONS..."
aws --endpoint-url=http://localhost:4566 apigateway put-integration \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method OPTIONS \
  --type MOCK \
  --request-templates '{"application/json": "{\"statusCode\": 200}"}'

# Add method response for OPTIONS
echo "[INFO] Adding method response for OPTIONS..."
aws --endpoint-url=http://localhost:4566 apigateway put-method-response \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method OPTIONS \
  --status-code 200 \
  --response-models '{"application/json": "Empty"}'

# Add integration response for OPTIONS with CORS headers
echo "[INFO] Adding integration response for OPTIONS with CORS headers..."
aws --endpoint-url=http://localhost:4566 apigateway put-integration-response \
  --rest-api-id $API_ID \
  --resource-id $UPLOAD_ID \
  --http-method OPTIONS \
  --status-code 200 \
  --response-parameters '{
    "method.response.header.Access-Control-Allow-Headers": "'\''Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'\''",
    "method.response.header.Access-Control-Allow-Methods": "'\''POST,OPTIONS'\''",
    "method.response.header.Access-Control-Allow-Origin": "'\''*'\''"
  }' \
  --response-templates '{"application/json": ""}'

# Deploy the API
echo "[INFO] Deploying API..."
aws --endpoint-url=http://localhost:4566 apigateway create-deployment \
  --rest-api-id $API_ID \
  --stage-name dev

# Add Lambda permission for API Gateway
echo "[INFO] Adding Lambda permission for API Gateway..."
aws --endpoint-url=http://localhost:4566 lambda add-permission \
  --function-name generate-presigned-url \
  --statement-id apigateway-invoke \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn arn:aws:execute-api:eu-west-1:000000000000:$API_ID/*/POST/upload

# Save API ID to a file for reference
echo $API_ID > .api-gateway-id

echo "[INFO] API Gateway setup complete!"
echo "[INFO] API Gateway URL: http://localhost:4566/restapis/$API_ID/dev/_user_request_/upload"
echo "[INFO] API ID saved to .api-gateway-id" 