#!/bin/bash
set -e

cd "$(dirname "$0")"

source ./setup-environment.sh

echo "[INFO] LocalStack is up. Creating resources..."

echo "[INFO] Restarting LocalStack to apply CORS configuration..."
docker-compose restart localstack
sleep 10

for bucket in raw-storage hls-storage dash-storage thumbnails-storage; do
  if aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 ls "s3://$bucket" 2>&1 | grep -q 'NoSuchBucket'; then
    echo "[INFO] Creating S3 bucket: $bucket"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api create-bucket --bucket $bucket --region $AWS_REGION
  else
    echo "[INFO] S3 bucket $bucket already exists."
  fi
done

echo "[INFO] Re-applying S3 CORS configuration..."
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket raw-storage --cors-configuration file://cors.json
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket hls-storage --cors-configuration file://cors.json
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket dash-storage --cors-configuration file://cors.json
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket thumbnails-storage --cors-configuration file://cors.json

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

echo "[INFO] AWS resources setup completed" 