#!/bin/bash
set -e

cd "$(dirname "$0")"

# Load configuration
if [ -f "../.env" ]; then
  source ../.env
else
  echo "[WARNING] config.env not found, using default values"
  AWS_REGION="us-east-1"
  AWS_ACCESS_KEY_ID="test"
  AWS_SECRET_ACCESS_KEY="test"
  LOCALSTACK_ENDPOINT="http://localstack:4566"
  LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"
  SQS_QUEUE_URL="http://localstack:4566/000000000000/transcoder-jobs"
  ANALYZE_QUEUE_URL="http://localstack:4566/000000000000/analyze-completed"
  DELETE_FILES_LAMBDA_ENDPOINT="http://localstack:4566/2015-03-31/functions/delete-files/invocations"
fi

# Hobby Streamer Build Script
# This script sets up the complete development environment including:
# - LocalStack (S3, SQS, Lambda)
# - Neo4j (Graph database for asset management)
# - Keycloak (Authentication)
# - Auth Service (REST API)
# - Asset Manager (GraphQL API)
# - Transcoder Service
# - CMS UI (React Native/Expo)

# Stop CMS UI if running
echo "[INFO] Stopping any existing CMS UI processes..."
if [ -d "../frontend/HobbyStreamerCMS" ]; then
  cd ../frontend/HobbyStreamerCMS
  
  # Stop any running npm processes for this project
  pkill -f "npm run web" || true
  pkill -f "expo start" || true
  
  # Also stop any processes on the web port
  lsof -ti:8081 | xargs kill -9 2>/dev/null || true
  
  cd ../../local
fi

# Generate Keycloak certificates if they don't exist
if [ ! -f "keycloak-certs/cert.pem" ] || [ ! -f "keycloak-certs/key.pem" ]; then
  echo "[INFO] Generating Keycloak HTTPS certificates..."
  ./generate-keycloak-certs.sh
fi

# Stop any existing Expo processes
echo "[INFO] Stopping any existing Expo processes..."
pkill -f "expo start" || true
sleep 2

# Export environment variables for docker-compose
echo "[INFO] Exporting environment variables..."
export AWS_REGION
export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export LOCALSTACK_ENDPOINT
export LOCALSTACK_EXTERNAL_ENDPOINT
export SQS_QUEUE_URL
export ANALYZE_QUEUE_URL
export DELETE_FILES_LAMBDA_ENDPOINT

# Stop all running containers
if docker-compose ps | grep -q 'Up'; then
  echo "[INFO] Stopping all running containers..."
  docker-compose down
fi

# Start infrastructure services (LocalStack, Keycloak, and Neo4j)
echo "[INFO] Starting infrastructure services..."
docker-compose up -d localstack keycloak neo4j

# Wait for LocalStack to be ready
until curl -s http://localhost:4566/health > /dev/null 2>&1; do
  echo "[INFO] Waiting for LocalStack to be ready..."
  sleep 2
done
echo "[INFO] LocalStack is up. Waiting for all services to be ready..."
sleep 10

# Wait for Keycloak to be ready
echo "[INFO] Waiting for Keycloak to be ready..."
until curl -s http://localhost:9090/realms/master > /dev/null 2>&1; do
  sleep 2
done

# Import the hobby realm if it doesn't exist
if ! curl -s http://localhost:9090/realms/hobby | grep -q '"realm":"hobby"'; then
  echo "[INFO] Importing Keycloak hobby realm..."
  docker exec hobby-streamer-keycloak-1 /opt/keycloak/bin/kc.sh import --file=/opt/keycloak/data/import/hobby-realm.json --override=true
else
  echo "[INFO] Keycloak hobby realm already exists."
fi

# Wait for Neo4j to be ready
echo "[INFO] Waiting for Neo4j to be ready..."
until curl -s http://localhost:7474 > /dev/null 2>&1; do
  echo "[INFO] Waiting for Neo4j to be ready..."
  sleep 3
done
echo "[INFO] Neo4j is up and running."


echo "[INFO] LocalStack is up. Creating resources..."

# Ensure LocalStack is restarted to pick up any env var changes (CORS)
echo "[INFO] Restarting LocalStack to apply CORS configuration..."
docker-compose restart localstack
sleep 10

# Create S3 buckets if they do not exist (must be before CORS)
for bucket in raw-storage transcoded-storage thumbnails-storage; do
  if aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 ls "s3://$bucket" 2>&1 | grep -q 'NoSuchBucket'; then
    echo "[INFO] Creating S3 bucket: $bucket"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api create-bucket --bucket $bucket --region $AWS_REGION
  else
    echo "[INFO] S3 bucket $bucket already exists."
  fi
done

# Re-apply S3 CORS config after LocalStack is up and buckets exist
echo "[INFO] Re-applying S3 CORS configuration..."
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket raw-storage --cors-configuration file://cors.json
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket transcoded-storage --cors-configuration file://cors.json
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket thumbnails-storage --cors-configuration file://cors.json

# Create SQS queues
if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name transcoder-jobs --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: transcoder-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name transcoder-jobs --region $AWS_REGION > /dev/null
else
  echo "[INFO] SQS queue transcoder-jobs already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name analyze-completed --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: analyze-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name analyze-completed --region $AWS_REGION > /dev/null
else
  echo "[INFO] SQS queue analyze-completed already exists."
fi

# Build and deploy the presigned upload URL Lambda
pushd ../backend/storage/cmd/generate_presigned_upload_url > /dev/null
echo "[INFO] Building presigned upload URL Lambda..."
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

# Create or update Lambda in LocalStack
if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name generate-presigned-url > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: generate-presigned-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name generate-presigned-url --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: generate-presigned-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name generate-presigned-url \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=$AWS_REGION,AWS_ENDPOINT=$LOCALSTACK_ENDPOINT}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null



# Build and deploy the delete files Lambda
pushd ../backend/storage/cmd/delete_files > /dev/null
echo "[INFO] Building delete files Lambda..."

# Ensure dependencies are resolved
echo "[INFO] Resolving dependencies..."
go mod tidy

echo "[INFO] Building Lambda function..."
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

# Create or update Lambda in LocalStack
if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name delete-files > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: delete-files"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name delete-files --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: delete-files"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name delete-files \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={AWS_ENDPOINT=$LOCALSTACK_ENDPOINT,AWS_REGION=$AWS_REGION}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

# Remove invalid API Gateway ID if the API does not exist
if [ -f ".api-gateway-id" ]; then
  API_ID=$(cat .api-gateway-id)
  if [ -z "$API_ID" ] || ! awslocal --no-cli-pager --region $AWS_REGION apigateway get-rest-api --rest-api-id $API_ID > /dev/null 2>&1; then
    echo "[INFO] Removing stale or empty API Gateway ID: $API_ID"
    rm -f .api-gateway-id
    API_ID=""
  fi
fi

# Setup API Gateway for Lambda
echo "[INFO] Setting up API Gateway for presigned URL Lambda..."

# Always ensure API Gateway exists and is properly configured
if [ ! -f ".api-gateway-id" ] || [ -z "$(cat .api-gateway-id)" ]; then
  echo "[INFO] Creating new API Gateway..."
  ./setup-api-gateway.sh
  API_ID=$(cat .api-gateway-id)
else
  API_ID=$(cat .api-gateway-id)
  echo "[INFO] Using existing API Gateway with ID: $API_ID"
fi

# Always redeploy the API Gateway stage to ensure CORS is active
echo "[INFO] Redeploying API Gateway stage to ensure CORS is active..."
awslocal --no-cli-pager --region $AWS_REGION apigateway create-deployment \
  --rest-api-id $API_ID \
  --stage-name dev > /dev/null

echo "[INFO] API Gateway URL: http://localhost:4566/restapis/$API_ID/dev/_user_request_/upload"

# Update frontend with new API Gateway ID and URLs
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

# Start and rebuild all backend services
echo "[INFO] Building and starting all backend services with the latest code..."
docker-compose up --build -d auth-service asset-manager transcoder

# Wait for services to be ready
echo "[INFO] Waiting for services to be ready..."
sleep 10

# Test GraphQL endpoint
echo "[INFO] Testing Asset Manager GraphQL endpoint..."
until curl -s -X POST http://localhost:8082/graphql -H "Content-Type: application/json" -d '{"query":"{ __schema { types { name } } }"}' > /dev/null 2>&1; do
  echo "[INFO] Waiting for Asset Manager GraphQL endpoint to be ready..."
  sleep 3
done
echo "[INFO] Asset Manager GraphQL endpoint is ready."

echo "[INFO] All services are up to date and running."
echo "[INFO] - Auth Service: http://localhost:8080"
echo "[INFO] - Asset Manager GraphQL: http://localhost:8082/query"
echo "[INFO] - Neo4j Browser: http://localhost:7474"
echo "[INFO] - Keycloak: http://localhost:9090"
echo "[INFO] - LocalStack: http://localhost:4566"

# Setup CMS UI
echo "[INFO] Setting up CMS UI..."
if [ -d "../frontend/HobbyStreamerCMS" ]; then
  cd ../frontend/HobbyStreamerCMS
  
  # Ensure correct Node.js version and clean install dependencies
  echo "[INFO] Ensuring correct Node.js version..."
  if command -v nvm &> /dev/null; then
    nvm use
  fi
  
  echo "[INFO] Cleaning and reinstalling CMS UI dependencies..."
  rm -rf node_modules package-lock.json
  npm install
  
  # Start the web application in the background
  echo "[INFO] Starting CMS UI web application..."
  nohup npm run web > web.log 2>&1 &
  
  # Wait for the web server to start
  echo "[INFO] Waiting for CMS UI web server to start..."
  sleep 10
  
  echo "[INFO] CMS UI web application started"
  echo "[INFO] - CMS UI Web: http://localhost:8081"
  echo "[INFO] To run on device/simulator:"
  echo "[INFO]   - Android: npm run android"
  echo "[INFO]   - iOS: npm run ios"
  
  cd ../../local
else
  echo "[WARNING] CMS UI directory not found at ../frontend/HobbyStreamerCMS"
fi

echo ""
echo "[INFO] Build completed successfully!"
echo "[INFO] All services are running and ready for development."
echo ""
echo "[INFO] Quick access:"
echo "[INFO] - CMS UI: http://localhost:8081"
echo "[INFO] - Asset Manager: http://localhost:8082/query"
echo "[INFO] - Auth Service: http://localhost:8080"
echo ""
echo "[INFO] To stop all services: docker-compose down"
echo "[INFO] To stop CMS UI: pkill -f 'npm run web'"
