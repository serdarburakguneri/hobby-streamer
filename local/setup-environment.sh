#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "[INFO] Loading environment configuration..."

if [ -f "../.env" ]; then
  source ../.env
else
  echo "[WARNING] .env not found, using default values"
  AWS_REGION="us-east-1"
  AWS_ACCESS_KEY_ID="test"
  AWS_SECRET_ACCESS_KEY="test"
  LOCALSTACK_ENDPOINT="http://localstack:4566"
  LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"
  SQS_QUEUE_URL="http://localstack:4566/000000000000/transcoder-jobs"
  ANALYZE_QUEUE_URL="http://localstack:4566/000000000000/analyze-completed"
  DELETE_FILES_LAMBDA_ENDPOINT="http://localstack:4566/2015-03-31/functions/delete-files/invocations"
fi

echo "[INFO] Exporting environment variables..."
export AWS_REGION
export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export LOCALSTACK_ENDPOINT
export LOCALSTACK_EXTERNAL_ENDPOINT
export SQS_QUEUE_URL
export ANALYZE_QUEUE_URL
export DELETE_FILES_LAMBDA_ENDPOINT

echo "[INFO] Environment setup completed" 