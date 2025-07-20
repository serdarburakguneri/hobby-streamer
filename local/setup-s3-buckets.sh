#!/bin/bash
set -e

# Ensure we're in the local directory
cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up Netflix-style single bucket storage with cross-region replication..."

echo "[INFO] Creating primary region bucket (us-east-1)..."
if aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 ls "s3://content-east" 2>&1 | grep -q 'NoSuchBucket'; then
  echo "[INFO] Creating S3 bucket: content-east"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api create-bucket --bucket content-east --region $AWS_REGION
else
  echo "[INFO] S3 bucket content-east already exists."
fi

echo "[INFO] Creating secondary region bucket (us-east-1)..."
if aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3 ls "s3://content-west" 2>&1 | grep -q 'NoSuchBucket'; then
  echo "[INFO] Creating S3 bucket: content-west"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api create-bucket --bucket content-west --region $AWS_REGION
else
  echo "[INFO] S3 bucket content-west already exists."
fi

echo "[INFO] Enabling versioning on both buckets..."
for bucket in content-east content-west; do
  echo "[INFO] Enabling versioning for $bucket"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-versioning \
    --bucket $bucket \
    --versioning-configuration Status=Enabled
done

echo "[INFO] Setting up cross-region replication..."

cat > replication-config.json << EOF
{
  "Role": "arn:aws:iam::000000000000:role/s3-replication-role",
  "Rules": [
    {
      "ID": "CrossRegionReplication",
      "Status": "Enabled",
      "Priority": 1,
      "DeleteMarkerReplication": { "Status": "Enabled" },
      "Destination": {
        "Bucket": "arn:aws:s3:::content-west",
        "StorageClass": "STANDARD"
      }
    }
  ]
}
EOF

echo "[INFO] Applying replication config for content-east -> content-west"
aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-replication \
  --bucket content-east \
  --replication-configuration file://replication-config.json

rm -f replication-config.json

echo "[INFO] Applying CORS configuration to both buckets..."
for bucket in content-east content-west; do
  echo "[INFO] Applying CORS to $bucket"
  aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT s3api put-bucket-cors \
    --bucket $bucket \
    --cors-configuration file://cors.json
done

echo "[INFO] Single bucket storage setup completed successfully!"
echo "[INFO] Primary bucket: content-east ($AWS_REGION)"
echo "[INFO] Secondary bucket: content-west ($AWS_REGION)"
echo "[INFO] Cross-region replication is enabled (simulated in local development)"
echo "[INFO] Content structure: {assetId}/{type}/{quality}/{filename}" 