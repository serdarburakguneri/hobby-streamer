#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up S3 event triggers for Lambda functions..."

echo "[INFO] Adding S3 event notification for raw-video-uploaded Lambda..."

# Create S3 event notification configuration
cat > /tmp/s3-notification-config.json << EOF
{
  "LambdaFunctionConfigurations": [
    {
      "Id": "raw-video-uploaded-trigger",
      "LambdaFunctionArn": "arn:aws:lambda:$AWS_REGION:000000000000:function:raw-video-uploaded",
      "Events": ["s3:ObjectCreated:*"],
            "Filter": {
        "Key": {
          "FilterRules": [
            {
              "Name": "suffix",
              "Value": ".mp4"
            }
          ]
        }
      }
    }
  ]
}
EOF

# Add notification configuration to S3 bucket
echo "[INFO] Adding notification configuration to content-east bucket..."
awslocal --no-cli-pager --region $AWS_REGION s3api put-bucket-notification-configuration \
  --bucket content-east \
  --notification-configuration file:///tmp/s3-notification-config.json

# Add S3 permission to invoke Lambda
echo "[INFO] Adding S3 permission to invoke raw-video-uploaded Lambda..."
if awslocal --no-cli-pager --region $AWS_REGION lambda get-policy --function-name raw-video-uploaded 2>/dev/null | grep -q "s3-invoke"; then
  echo "[INFO] S3 permission already exists for raw-video-uploaded Lambda"
else
  awslocal --no-cli-pager --region $AWS_REGION lambda add-permission \
    --function-name raw-video-uploaded \
    --statement-id s3-invoke \
    --action lambda:InvokeFunction \
    --principal s3.amazonaws.com \
    --source-arn "arn:aws:s3:::content-east"
  echo "[INFO] S3 permission added successfully"
fi

echo "[INFO] S3 event triggers setup completed successfully!"
echo "[INFO] Raw video uploads to s3://content-east/raw-storage/ will trigger the raw-video-uploaded Lambda" 