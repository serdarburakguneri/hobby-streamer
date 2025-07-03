#!/bin/bash
set -e

# Start LocalStack
if ! docker-compose ps | grep -q localstack; then
  echo "[INFO] Starting LocalStack via docker-compose..."
  docker-compose up -d
else
  echo "[INFO] LocalStack already running."
fi

# Wait for LocalStack to be ready
until curl -s http://localhost:4566/health | grep '"s3": *"running"' > /dev/null; do
  echo "[INFO] Waiting for LocalStack to be ready..."
  sleep 2
done
sleep 2

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

# Build and deploy the presigned upload URL Lambda
pushd services/storage/cmd/generate_presigned_upload_url > /dev/null
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
    --role arn:aws:iam::000000000000:role/lambda-role
fi

# Set environment variables for Lambda
awslocal lambda update-function-configuration \
  --function-name generate-presigned-url \
  --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"

popd > /dev/null

# Create DynamoDB tables
for table in asset bucket; do
  if ! aws --endpoint-url=http://localhost:4566 dynamodb describe-table --table-name "$table" --region us-west-2 > /dev/null 2>&1; then
    echo "[INFO] Creating DynamoDB table: $table"
    aws --endpoint-url=http://localhost:4566 dynamodb create-table \
      --table-name "$table" \
      --attribute-definitions AttributeName=id,AttributeType=S \
      --key-schema AttributeName=id,KeyType=HASH \
      --billing-mode PAYPERREQUEST \
      --region us-west-2
  else
    echo "[INFO] DynamoDB table $table already exists."
  fi
done

# Create SQS queue
if ! aws --endpoint-url=http://localhost:4566 sqs get-queue-url --queue-name transcoder-jobs --region us-west-2 > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: transcoder-jobs"
  aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name transcoder-jobs --region us-west-2
else
  echo "[INFO] SQS queue transcoder-jobs already exists."
fi

echo "[INFO] Local environment setup complete."

# Run tests before starting services
pushd services/asset-manager > /dev/null
echo "[INFO] Running asset-manager tests..."
go test ./... -v || { echo '[ERROR] Asset-manager tests failed.'; exit 1; }
popd > /dev/null

pushd services/transcoder > /dev/null
echo "[INFO] Running transcoder tests..."
go test ./... -v || { echo '[ERROR] Transcoder tests failed.'; exit 1; }
popd > /dev/null

# Start Asset Manager service
pushd services/asset-manager > /dev/null
if pgrep -f "go run ./cmd/main.go" > /dev/null; then
  echo "[INFO] Asset Manager service already running."
else
  echo "[INFO] Starting Asset Manager service on port 8080..."
  PORT=8080 nohup go run ./cmd/main.go > ../../asset-manager.log 2>&1 &
fi
popd > /dev/null

# Start Transcoder service
pushd services/transcoder > /dev/null
if pgrep -f "go run ./cmd/worker/main.go" > /dev/null; then
  echo "[INFO] Transcoder service already running."
else
  echo "[INFO] Starting Transcoder service..."
  SQS_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs nohup go run ./cmd/worker/main.go > ../../transcoder.log 2>&1 &
fi
popd > /dev/null

echo "[INFO] All services started. Check asset-manager.log and transcoder.log for output."
