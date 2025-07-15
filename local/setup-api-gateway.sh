#!/bin/bash
set -e

# Load configuration
if [ -f "config.env" ]; then
  source config.env
else
  echo "[WARNING] config.env not found, using default values"
  AWS_REGION="us-east-1"
  LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"
fi

echo "[INFO] Setting up API Gateway for all services..."

# Create REST API
echo "[INFO] Creating REST API..."
API_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-rest-api --name "hobby-streamer-api")
API_ID=$(echo $API_RESPONSE | jq -r '.id')
echo "[INFO] API ID: $API_ID"

# Get root resource ID
echo "[INFO] Getting root resource..."
ROOT_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway get-resources --rest-api-id $API_ID)
ROOT_ID=$(echo $ROOT_RESPONSE | jq -r '.items[] | select(.path == "/") | .id')
echo "[INFO] Root Resource ID: $ROOT_ID"

# Create auth resource
echo "[INFO] Creating auth resource..."
AUTH_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part auth)
AUTH_ID=$(echo $AUTH_RESPONSE | jq -r '.id')
echo "[INFO] Auth Resource ID: $AUTH_ID"

# Create auth/login resource
echo "[INFO] Creating auth/login resource..."
LOGIN_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part login)
LOGIN_ID=$(echo $LOGIN_RESPONSE | jq -r '.id')
echo "[INFO] Login Resource ID: $LOGIN_ID"

# Create auth/validate resource
echo "[INFO] Creating auth/validate resource..."
VALIDATE_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part validate)
VALIDATE_ID=$(echo $VALIDATE_RESPONSE | jq -r '.id')
echo "[INFO] Validate Resource ID: $VALIDATE_ID"

# Create auth/refresh resource
echo "[INFO] Creating auth/refresh resource..."
REFRESH_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part refresh)
REFRESH_ID=$(echo $REFRESH_RESPONSE | jq -r '.id')
echo "[INFO] Refresh Resource ID: $REFRESH_ID"

# Create graphql resource
echo "[INFO] Creating graphql resource..."
GRAPHQL_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part graphql)
GRAPHQL_ID=$(echo $GRAPHQL_RESPONSE | jq -r '.id')
echo "[INFO] GraphQL Resource ID: $GRAPHQL_ID"

# Create upload resource
echo "[INFO] Creating upload resource..."
UPLOAD_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part upload)
UPLOAD_ID=$(echo $UPLOAD_RESPONSE | jq -r '.id')
echo "[INFO] Upload Resource ID: $UPLOAD_ID"

# Create transcode resource
echo "[INFO] Creating transcode resource..."
TRANSCODE_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part transcode)
TRANSCODE_ID=$(echo $TRANSCODE_RESPONSE | jq -r '.id')
echo "[INFO] Transcode Resource ID: $TRANSCODE_ID"

# Function to create method with CORS
create_method_with_cors() {
    local resource_id=$1
    local http_method=$2
    local integration_type=$3
    local integration_uri=$4
    
    echo "[INFO] Creating $http_method method..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-method \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method $http_method \
      --authorization-type "NONE"
    
    echo "[INFO] Integrating $http_method method..."
    if [ "$integration_type" = "HTTP_PROXY" ]; then
        aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-integration \
          --rest-api-id $API_ID \
          --resource-id $resource_id \
          --http-method $http_method \
          --type HTTP_PROXY \
          --integration-http-method $http_method \
          --uri $integration_uri
    elif [ "$integration_type" = "AWS_PROXY" ]; then
        aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-integration \
          --rest-api-id $API_ID \
          --resource-id $resource_id \
          --http-method $http_method \
          --type AWS_PROXY \
          --integration-http-method POST \
          --uri $integration_uri
    fi
    
    echo "[INFO] Adding method response for $http_method..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-method-response \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method $http_method \
      --status-code 200 \
      --response-parameters '{
        "method.response.header.Access-Control-Allow-Origin": true,
        "method.response.header.Access-Control-Allow-Headers": true,
        "method.response.header.Access-Control-Allow-Methods": true
      }'
    
    if [ "$integration_type" = "HTTP_PROXY" ]; then
        echo "[INFO] Adding integration response for $http_method with CORS headers..."
        aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-integration-response \
          --rest-api-id $API_ID \
          --resource-id $resource_id \
          --http-method $http_method \
          --status-code 200 \
          --response-parameters '{
            "method.response.header.Access-Control-Allow-Origin": "'\''*'\''",
            "method.response.header.Access-Control-Allow-Headers": "'\''Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'\''",
            "method.response.header.Access-Control-Allow-Methods": "'\''GET,POST,PUT,DELETE,OPTIONS'\''"
          }'
    fi
}

# Function to create OPTIONS method for CORS
create_options_method() {
    local resource_id=$1
    
    echo "[INFO] Creating OPTIONS method for CORS..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-method \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method OPTIONS \
      --authorization-type "NONE"
    
    echo "[INFO] Creating mock integration for OPTIONS..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-integration \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method OPTIONS \
      --type MOCK \
      --request-templates '{"application/json": "{\"statusCode\": 200}"}'
    
    echo "[INFO] Adding method response for OPTIONS..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-method-response \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method OPTIONS \
      --status-code 200 \
      --response-parameters '{
        "method.response.header.Access-Control-Allow-Origin": true,
        "method.response.header.Access-Control-Allow-Headers": true,
        "method.response.header.Access-Control-Allow-Methods": true
      }' \
      --response-models '{"application/json": "Empty"}'
    
    echo "[INFO] Adding integration response for OPTIONS with CORS headers..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway put-integration-response \
      --rest-api-id $API_ID \
      --resource-id $resource_id \
      --http-method OPTIONS \
      --status-code 200 \
      --response-parameters '{
        "method.response.header.Access-Control-Allow-Headers": "'\''Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'\''",
        "method.response.header.Access-Control-Allow-Methods": "'\''GET,POST,PUT,DELETE,OPTIONS'\''",
        "method.response.header.Access-Control-Allow-Origin": "'\''*'\''"
      }' \
      --response-templates '{"application/json": ""}'
}

# Create auth endpoints
echo "[INFO] Setting up auth endpoints..."

# Login endpoint
create_method_with_cors $LOGIN_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/login"
create_options_method $LOGIN_ID

# Validate endpoint
create_method_with_cors $VALIDATE_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/validate"
create_options_method $VALIDATE_ID

# Refresh endpoint
create_method_with_cors $REFRESH_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/refresh"
create_options_method $REFRESH_ID

# Create GraphQL endpoint
echo "[INFO] Setting up GraphQL endpoint..."
create_method_with_cors $GRAPHQL_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8082/graphql"
create_method_with_cors $GRAPHQL_ID "GET" "HTTP_PROXY" "http://host.docker.internal:8082/graphql"
create_options_method $GRAPHQL_ID

# Create upload endpoint (Lambda integration)
echo "[INFO] Setting up upload endpoint..."
create_method_with_cors $UPLOAD_ID "POST" "AWS_PROXY" "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_REGION:000000000000:function:generate-presigned-url/invocations"
create_options_method $UPLOAD_ID

# Create transcode endpoint (Lambda integration)
echo "[INFO] Setting up transcode endpoint..."
create_method_with_cors $TRANSCODE_ID "POST" "AWS_PROXY" "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_REGION:000000000000:function:trigger-transcode-job/invocations"
create_options_method $TRANSCODE_ID

# Deploy the API
echo "[INFO] Deploying API..."
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-deployment \
  --rest-api-id $API_ID \
  --stage-name dev

# Add Lambda permission for API Gateway
echo "[INFO] Adding Lambda permission for API Gateway..."
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT lambda add-permission \
  --function-name generate-presigned-url \
  --statement-id apigateway-invoke \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:$AWS_REGION:000000000000:$API_ID/*/POST/upload"

aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT lambda add-permission \
  --function-name trigger-transcode-job \
  --statement-id apigateway-invoke \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:$AWS_REGION:000000000000:$API_ID/*/POST/transcode"

# Save API ID to a file for reference
echo $API_ID > .api-gateway-id

echo "[INFO] API Gateway setup complete!"
echo "[INFO] API Gateway Base URL: http://localhost:4566/_aws/execute-api/$API_ID/dev"
echo "[INFO] Available endpoints:"
echo "[INFO] - Auth Login: POST /auth/login"
echo "[INFO] - Auth Validate: POST /auth/validate"
echo "[INFO] - Auth Refresh: POST /auth/refresh"
echo "[INFO] - GraphQL: POST/GET /graphql"
echo "[INFO] - Upload: POST /upload"
echo "[INFO] - Transcode: POST /transcode"
echo "[INFO] API ID saved to .api-gateway-id" 