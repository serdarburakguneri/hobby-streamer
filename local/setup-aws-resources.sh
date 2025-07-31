#!/bin/bash

set -e

AWS_REGION="us-east-1"
LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"

echo "[INFO] Setting up AWS resources for Hobby Streamer..."

# Create S3 buckets
BUCKETS=("content-east" "raw-storage" "hls-storage" "dash-storage")

for bucket in "${BUCKETS[@]}"; do
    if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 ls s3://$bucket > /dev/null 2>&1; then
        echo "[INFO] Creating S3 bucket: $bucket"
        aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 mb s3://$bucket --region $AWS_REGION > /dev/null
    else
        echo "[INFO] S3 bucket $bucket already exists."
    fi
done

echo "[INFO] AWS resources setup completed successfully!" 