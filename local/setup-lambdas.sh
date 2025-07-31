#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

pushd ../backend/lambdas/cmd/generate_video_upload_url > /dev/null
echo "[INFO] Building video upload URL Lambda..."
echo "[INFO] Resolving dependencies..."
go mod tidy

GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name generate-video-upload-url > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: generate-video-upload-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name generate-video-upload-url --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: generate-video-upload-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name generate-video-upload-url \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={BUCKET_NAME=content-east,BUCKET_REGION=$AWS_REGION,AWS_ENDPOINT=$LOCALSTACK_INTERNAL_ENDPOINT}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

pushd ../backend/lambdas/cmd/delete_files > /dev/null
echo "[INFO] Building delete files Lambda..."

echo "[INFO] Resolving dependencies..."
go mod tidy

echo "[INFO] Building Lambda function..."
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

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
    --environment "Variables={AWS_ENDPOINT=$LOCALSTACK_INTERNAL_ENDPOINT,AWS_REGION=$AWS_REGION}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

pushd ../backend/lambdas/cmd/generate_image_upload_url > /dev/null
echo "[INFO] Building image upload URL Lambda..."
echo "[INFO] Resolving dependencies..."
go mod tidy

GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name generate-image-upload-url > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: generate-image-upload-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name generate-image-upload-url --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: generate-image-upload-url"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name generate-image-upload-url \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={BUCKET_NAME=content-east,BUCKET_REGION=$AWS_REGION,AWS_ENDPOINT=$LOCALSTACK_INTERNAL_ENDPOINT}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

pushd ../backend/lambdas/cmd/raw_video_uploaded > /dev/null
echo "[INFO] Building raw video uploaded Lambda..."
echo "[INFO] Resolving dependencies..."
go mod tidy

GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name raw-video-uploaded > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: raw-video-uploaded"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name raw-video-uploaded --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: raw-video-uploaded"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name raw-video-uploaded \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={KAFKA_BOOTSTRAP_SERVERS=kafka:29092}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

pushd ../backend/lambdas/cmd/hls_job_requested > /dev/null
echo "[INFO] Building HLS job requested Lambda..."
echo "[INFO] Resolving dependencies..."
go mod tidy

GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

if awslocal --no-cli-pager --region $AWS_REGION lambda get-function --function-name hls-job-requested > /dev/null 2>&1; then
  echo "[INFO] Updating existing Lambda function: hls-job-requested"
  awslocal --no-cli-pager --region $AWS_REGION lambda update-function-code --function-name hls-job-requested --zip-file fileb://function.zip > /dev/null
else
  echo "[INFO] Creating Lambda function: hls-job-requested"
  awslocal --no-cli-pager --region $AWS_REGION lambda create-function \
    --function-name hls-job-requested \
    --runtime go1.x \
    --handler main \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --environment "Variables={KAFKA_BOOTSTRAP_SERVERS=kafka:29092}" \
    --region $AWS_REGION > /dev/null
fi

popd > /dev/null

echo "[INFO] Setting up S3 event triggers..."
./setup-s3-lambda-triggers.sh

echo "[INFO] Lambda functions setup completed" 