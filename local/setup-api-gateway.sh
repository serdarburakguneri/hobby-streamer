#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

if [ -f ".api-gateway-id" ]; then
  API_ID=$(cat .api-gateway-id)
  if [ -z "$API_ID" ] || ! awslocal --no-cli-pager --region $AWS_REGION apigateway get-rest-api --rest-api-id $API_ID > /dev/null 2>&1; then
    echo "[INFO] Removing stale or empty API Gateway ID: $API_ID"
    rm -f .api-gateway-id
    API_ID=""
  fi
fi

echo "[INFO] Setting up API Gateway for all services..."

if [ ! -f ".api-gateway-id" ] || [ -z "$(cat .api-gateway-id)" ]; then
  echo "[INFO] Creating new API Gateway..."
  
  echo "[INFO] Creating REST API..."
  API_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-rest-api --name "hobby-streamer-api")
  API_ID=$(echo $API_RESPONSE | jq -r '.id')
  echo "[INFO] API ID: $API_ID"

  echo "[INFO] Getting root resource..."
  ROOT_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway get-resources --rest-api-id $API_ID)
  ROOT_ID=$(echo $ROOT_RESPONSE | jq -r '.items[] | select(.path == "/") | .id')
  echo "[INFO] Root Resource ID: $ROOT_ID"

  echo "[INFO] Creating auth resource..."
  AUTH_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part auth)
  AUTH_ID=$(echo $AUTH_RESPONSE | jq -r '.id')
  echo "[INFO] Auth Resource ID: $AUTH_ID"

  echo "[INFO] Creating auth/login resource..."
  LOGIN_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part login)
  LOGIN_ID=$(echo $LOGIN_RESPONSE | jq -r '.id')
  echo "[INFO] Login Resource ID: $LOGIN_ID"

  echo "[INFO] Creating auth/validate resource..."
  VALIDATE_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part validate)
  VALIDATE_ID=$(echo $VALIDATE_RESPONSE | jq -r '.id')
  echo "[INFO] Validate Resource ID: $VALIDATE_ID"

  echo "[INFO] Creating auth/refresh resource..."
  REFRESH_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $AUTH_ID --path-part refresh)
  REFRESH_ID=$(echo $REFRESH_RESPONSE | jq -r '.id')
  echo "[INFO] Refresh Resource ID: $REFRESH_ID"

  echo "[INFO] Creating graphql resource..."
  GRAPHQL_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part graphql)
  GRAPHQL_ID=$(echo $GRAPHQL_RESPONSE | jq -r '.id')
  echo "[INFO] GraphQL Resource ID: $GRAPHQL_ID"

  echo "[INFO] Creating upload resource..."
  UPLOAD_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part upload)
  UPLOAD_ID=$(echo $UPLOAD_RESPONSE | jq -r '.id')
  echo "[INFO] Upload Resource ID: $UPLOAD_ID"

  echo "[INFO] Creating image-upload resource..."
  IMAGE_UPLOAD_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part image-upload)
  IMAGE_UPLOAD_ID=$(echo $IMAGE_UPLOAD_RESPONSE | jq -r '.id')
  echo "[INFO] Image Upload Resource ID: $IMAGE_UPLOAD_ID"

  echo "[INFO] Creating hls-job-requested resource..."
  HLS_JOB_RESPONSE=$(aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT_ID --path-part hls-job-requested)
  HLS_JOB_ID=$(echo $HLS_JOB_RESPONSE | jq -r '.id')
  echo "[INFO] HLS Job Requested Resource ID: $HLS_JOB_ID"

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
              "method.response.header.Access-Control-Allow-Origin": "'\''http://localhost:8081'\''",
              "method.response.header.Access-Control-Allow-Headers": "'\''Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'\''",
              "method.response.header.Access-Control-Allow-Methods": "'\''GET,POST,PUT,DELETE,OPTIONS'\''"
            }'
      fi
  }

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
          "method.response.header.Access-Control-Allow-Origin": "'\''http://localhost:8081'\''"
        }' \
        --response-templates '{"application/json": ""}'
  }

  echo "[INFO] Setting up auth endpoints..."
  create_method_with_cors $LOGIN_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/login"
  create_options_method $LOGIN_ID

  create_method_with_cors $VALIDATE_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/validate"
  create_options_method $VALIDATE_ID

  create_method_with_cors $REFRESH_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8080/auth/refresh"
  create_options_method $REFRESH_ID

  echo "[INFO] Setting up GraphQL endpoint..."
  create_method_with_cors $GRAPHQL_ID "POST" "HTTP_PROXY" "http://host.docker.internal:8082/graphql"
  create_method_with_cors $GRAPHQL_ID "GET" "HTTP_PROXY" "http://host.docker.internal:8082/graphql"
  create_options_method $GRAPHQL_ID

  echo "[INFO] Setting up video upload endpoint..."
  create_method_with_cors $UPLOAD_ID "POST" "AWS_PROXY" "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_REGION:000000000000:function:generate-video-upload-url/invocations"
  create_options_method $UPLOAD_ID

  echo "[INFO] Setting up image upload endpoint..."
  create_method_with_cors $IMAGE_UPLOAD_ID "POST" "AWS_PROXY" "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_REGION:000000000000:function:generate-image-upload-url/invocations"
  create_options_method $IMAGE_UPLOAD_ID

  echo "[INFO] Setting up HLS job requested endpoint..."
  create_method_with_cors $HLS_JOB_ID "POST" "AWS_PROXY" "arn:aws:apigateway:$AWS_REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$AWS_REGION:000000000000:function:hls-job-requested/invocations"
  create_options_method $HLS_JOB_ID

  echo "[INFO] Deploying API..."
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT apigateway create-deployment \
    --rest-api-id $API_ID \
    --stage-name dev

  echo "[INFO] Adding Lambda permission for API Gateway..."
  
  add_lambda_permission() {
    local function_name=$1
    local source_arn=$2
    
    if aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT lambda get-function --function-name $function_name > /dev/null 2>&1; then
      echo "[INFO] Adding permission for $function_name..."
      aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT lambda add-permission \
        --function-name $function_name \
        --statement-id apigateway-invoke \
        --action lambda:InvokeFunction \
        --principal apigateway.amazonaws.com \
        --source-arn "$source_arn" 2>/dev/null || echo "[WARN] Permission already exists for $function_name"
    else
      echo "[WARN] Lambda function $function_name not found, skipping permission"
    fi
  }
  
  add_lambda_permission "generate-video-upload-url" "arn:aws:execute-api:$AWS_REGION:000000000000:$API_ID/*/POST/upload"
  add_lambda_permission "generate-image-upload-url" "arn:aws:execute-api:$AWS_REGION:000000000000:$API_ID/*/POST/image-upload"
  add_lambda_permission "hls-job-requested" "arn:aws:execute-api:$AWS_REGION:000000000000:$API_ID/*/POST/hls-job-requested"

  echo $API_ID > .api-gateway-id
  echo "[INFO] API ID saved to .api-gateway-id"
else
  API_ID=$(cat .api-gateway-id)
  echo "[INFO] Using existing API Gateway with ID: $API_ID"
fi

echo "[INFO] Redeploying API Gateway stage to ensure CORS is active..."
awslocal --no-cli-pager --region $AWS_REGION apigateway create-deployment \
  --rest-api-id $API_ID \
  --stage-name dev > /dev/null

echo "[INFO] API Gateway URL: http://localhost:4566/restapis/$API_ID/dev/_user_request_/upload"

echo "[INFO] Updating frontend with new API Gateway configuration..."
sed -i.bak "s/getEnvVar('REACT_APP_API_GATEWAY_ID', '[^']*')/getEnvVar('REACT_APP_API_GATEWAY_ID', '$API_ID')/" \
  ../frontend/HobbyStreamerCMS/src/config/api.ts
sed -i.bak "s|getEnvVar('REACT_APP_API_GATEWAY_BASE_URL', '[^']*')|getEnvVar('REACT_APP_API_GATEWAY_BASE_URL', 'http://localhost:4566/_aws/execute-api/$API_ID/dev')|" \
  ../frontend/HobbyStreamerCMS/src/config/api.ts
sed -i.bak "s|getEnvVar('REACT_APP_AUTH_BASE_URL', '[^']*')|getEnvVar('REACT_APP_AUTH_BASE_URL', 'http://localhost:4566/_aws/execute-api/$API_ID/dev/auth')|" \
  ../frontend/HobbyStreamerCMS/src/config/api.ts
sed -i.bak "s|getEnvVar('REACT_APP_GRAPHQL_BASE_URL', '[^']*')|getEnvVar('REACT_APP_GRAPHQL_BASE_URL', 'http://localhost:4566/_aws/execute-api/$API_ID/dev/graphql')|" \
  ../frontend/HobbyStreamerCMS/src/config/api.ts
rm -f ../frontend/HobbyStreamerCMS/src/config/api.ts.bak
echo "[INFO] Frontend updated with API Gateway ID: $API_ID"

echo "[INFO] API Gateway setup completed" 