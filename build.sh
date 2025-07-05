#!/bin/bash
set -e

# Stop all running containers
if docker-compose ps | grep -q 'Up'; then
  echo "[INFO] Stopping all running containers..."
  docker-compose down
fi

# Start infrastructure services (LocalStack and Keycloak)
echo "[INFO] Starting infrastructure services..."
docker-compose up -d localstack keycloak

# Wait for LocalStack to be ready
until curl -s http://localhost:4566/health > /dev/null 2>&1; do
  echo "[INFO] Waiting for LocalStack to be ready..."
  sleep 2
done
echo "[INFO] LocalStack is up. Waiting for all services to be ready..."
sleep 10

# Test if DynamoDB is working
echo "[INFO] Testing DynamoDB connectivity..."
until aws --endpoint-url=http://localhost:4566 dynamodb list-tables --region us-west-2 > /dev/null 2>&1; do
  echo "[INFO] Waiting for DynamoDB to be ready..."
  sleep 3
done

echo "[INFO] LocalStack is up. Creating resources..."

# Create S3 buckets
for bucket in raw-storage transcoded-storage thumbnails-storage; do
  if aws --endpoint-url=http://localhost:4566 s3 ls "s3://$bucket" 2>&1 | grep -q 'NoSuchBucket'; then
    echo "[INFO] Creating S3 bucket: $bucket"
    aws --endpoint-url=http://localhost:4566 s3 mb s3://$bucket
  else
    echo "[INFO] S3 bucket $bucket already exists."
  fi
done

# Create DynamoDB tables
for table in asset bucket; do
  echo "[INFO] Checking DynamoDB table: $table"
  if ! aws --endpoint-url=http://localhost:4566 dynamodb describe-table --table-name "$table" --region us-west-2 > /dev/null 2>&1; then
    echo "[INFO] Creating DynamoDB table: $table"
    aws --endpoint-url=http://localhost:4566 dynamodb create-table \
      --table-name "$table" \
      --attribute-definitions AttributeName=id,AttributeType=S \
      --key-schema AttributeName=id,KeyType=HASH \
      --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
      --region us-west-2 > /dev/null
    echo "[INFO] Waiting for table $table to be active..."
    sleep 5
  else
    echo "[INFO] DynamoDB table $table already exists."
  fi
done

# Create SQS queue
if ! aws --endpoint-url=http://localhost:4566 sqs get-queue-url --queue-name transcoder-jobs --region us-west-2 > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: transcoder-jobs"
  aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name transcoder-jobs --region us-west-2 > /dev/null
else
  echo "[INFO] SQS queue transcoder-jobs already exists."
fi

# Build and deploy the presigned upload URL Lambda
pushd backend/storage/cmd/generate_presigned_upload_url > /dev/null
echo "[INFO] Building presigned upload URL Lambda..."
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

# Create or update Lambda in LocalStack
if awslocal lambda get-function --function-name generate-presigned-url > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: generate-presigned-url"
  awslocal lambda update-function-code --function-name generate-presigned-url --zip-file fileb://function.zip
else
  echo "[INFO] Creating Lambda function: generate-presigned-url"
  awslocal lambda create-function \
    --function-name generate-presigned-url \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
fi

popd > /dev/null

# Start and rebuild all services
echo "[INFO] Building and starting all services with the latest code..."
docker-compose up --build -d

echo "[INFO] All services are up to date and running."
echo "[INFO] - Auth Service: http://localhost:8080"
echo "[INFO] - Asset Manager: http://localhost:8082"
echo "[INFO] - Keycloak: http://localhost:9090"
echo "[INFO] - LocalStack: http://localhost:4566"
