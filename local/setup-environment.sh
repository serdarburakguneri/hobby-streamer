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
  HLS_QUEUE_URL="http://localstack:4566/000000000000/hls-jobs"
  DASH_QUEUE_URL="http://localstack:4566/000000000000/dash-jobs"
  ANALYZE_JOBS_QUEUE_URL="http://localstack:4566/000000000000/analyze-jobs"
ANALYZE_COMPLETED_QUEUE_URL="http://localstack:4566/000000000000/analyze-completed"
  HLS_COMPLETED_QUEUE_URL="http://localstack:4566/000000000000/hls-completed"
  DASH_COMPLETED_QUEUE_URL="http://localstack:4566/000000000000/dash-completed"
  DELETE_FILES_LAMBDA_ENDPOINT="http://localstack:4566/2015-03-31/functions/delete-files/invocations"
fi

echo "[INFO] Exporting environment variables..."
export AWS_REGION
export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export LOCALSTACK_ENDPOINT
export LOCALSTACK_EXTERNAL_ENDPOINT
export HLS_QUEUE_URL
export DASH_QUEUE_URL
export ANALYZE_JOBS_QUEUE_URL
export ANALYZE_COMPLETED_QUEUE_URL
export HLS_COMPLETED_QUEUE_URL
export DASH_COMPLETED_QUEUE_URL
export DELETE_FILES_LAMBDA_ENDPOINT

echo "[INFO] Environment setup completed" 