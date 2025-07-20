#!/bin/bash
set -e

# Ensure we're in the local directory
cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Testing SQS queues for Hobby Streamer..."

echo "[INFO] Testing job-queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/job-queue \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] Testing completion-queue attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/completion-queue \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] Testing asset-events attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/asset-events \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] Testing job-queue-dlq attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/job-queue-dlq \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] Testing completion-queue-dlq attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/completion-queue-dlq \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] Testing asset-events-dlq attributes:"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-attributes \
    --queue-url http://localhost:4566/000000000000/asset-events-dlq \
    --attribute-names All \
    --region $AWS_REGION

echo "[INFO] SQS queue testing completed successfully!" 