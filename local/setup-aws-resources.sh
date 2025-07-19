#!/bin/bash
set -e

source "$(dirname "$0")/setup-environment.sh"

echo "[INFO] LocalStack is up. Creating resources..."

echo "[INFO] Restarting LocalStack to apply CORS configuration..."
docker-compose restart localstack
sleep 10

for bucket in raw-storage hls-storage dash-storage images-storage; do
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
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors --bucket images-storage --cors-configuration file://cors.json

echo "[INFO] Creating SQS queues..."

# Create DLQ queues first
for queue in analyze-jobs-dlq analyze-completed-dlq hls-completed-dlq dash-completed-dlq hls-jobs-dlq dash-jobs-dlq; do
  if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name $queue --region $AWS_REGION > /dev/null 2>&1; then
    echo "[INFO] Creating SQS DLQ: $queue"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name $queue --region $AWS_REGION > /dev/null
  else
    echo "[INFO] SQS DLQ $queue already exists."
  fi
done

# Get DLQ ARNs for redrive policy
ANALYZE_JOBS_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:analyze-jobs-dlq"
HLS_JOBS_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:hls-jobs-dlq"
DASH_JOBS_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:dash-jobs-dlq"
ANALYZE_COMPLETED_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:analyze-completed-dlq"
HLS_COMPLETED_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:hls-completed-dlq"
DASH_COMPLETED_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:dash-completed-dlq"

# Create main queues with redrive policy
if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name analyze-jobs --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: analyze-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name analyze-jobs --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for analyze-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/analyze-jobs \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$ANALYZE_JOBS_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue analyze-jobs already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name hls-jobs --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: hls-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name hls-jobs --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for hls-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/hls-jobs \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$HLS_JOBS_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue hls-jobs already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name dash-jobs --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: dash-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name dash-jobs --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for dash-jobs"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/dash-jobs \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$DASH_JOBS_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue dash-jobs already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name analyze-completed --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: analyze-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name analyze-completed --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for analyze-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/analyze-completed \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$ANALYZE_COMPLETED_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue analyze-completed already exists."
fi

# Create completion queues with redrive policy
if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name hls-completed --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: hls-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name hls-completed --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for hls-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/hls-completed \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$HLS_COMPLETED_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue hls-completed already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name dash-completed --region $AWS_REGION > /dev/null 2>&1; then
  echo "[INFO] Creating SQS queue: dash-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name dash-completed --region $AWS_REGION > /dev/null
  echo "[INFO] Setting redrive policy for dash-completed"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
    --queue-url http://localhost:4566/000000000000/dash-completed \
    --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'"$DASH_COMPLETED_DLQ_ARN"'\",\"maxReceiveCount\":3}"}'
else
  echo "[INFO] SQS queue dash-completed already exists."
fi

echo "[INFO] AWS resources setup completed" 