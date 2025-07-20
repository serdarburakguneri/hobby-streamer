#!/bin/bash

set -e

AWS_REGION="us-east-1"
LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"

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

echo "[INFO] SQS queue testing completed successfully!" 