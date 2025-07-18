#!/bin/bash
set -e

source "$(dirname "$0")/setup-environment.sh"

echo "[INFO] Testing SQS queues..."

echo "[INFO] Listing all queues:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs list-queues --region $AWS_REGION

echo "[INFO] Testing HLS queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/hls-jobs \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Testing DASH queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/dash-jobs \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Testing analyze-jobs queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/analyze-jobs \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Testing analyze-completed queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/analyze-completed \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Testing HLS completed queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/hls-completed \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Testing DASH completed queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/dash-completed \
  --attribute-names All \
  --region $AWS_REGION

echo "[INFO] Queue testing completed" 