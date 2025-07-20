# Delete Files Lambda

Go-based AWS Lambda function for deleting S3 files related to an asset. Used for cleanup when assets are removed.

## Features

Deletes files from multiple S3 buckets (raw, transcoded, thumbnails), returns detailed success/failure info.

## Request

```json
{
  "assetId": "asset123",
  "files": [
    {"bucket": "raw-storage", "key": "asset123/main.mp4"},
    {"bucket": "transcoded-storage", "key": "asset123/main.m3u8"},
    {"bucket": "images-storage", "key": "asset123/thumbnail.jpg"}
  ]
}
```

## Response

```json
{
  "message": "Deleted 3 files for asset asset123",
  "deleted": [{"bucket": "raw-storage", "key": "asset123/main.mp4"}],
  "errors": ["Failed to delete asset123/thumbnail.jpg: NoSuchKey"]
}
```

## Local Development

```bash
# Build
cd backend/lambdas/cmd/delete_files
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main

# Deploy to LocalStack
awslocal lambda create-function \
  --function-name delete-files \
  --runtime go1.x --handler main --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --environment "Variables={AWS_ENDPOINT=http://localhost:4566,AWS_REGION=eu-west-1}"

# Test
awslocal lambda invoke \
  --function-name delete-files \
  --payload '{"assetId":"test123","files":[{"bucket":"raw-storage","key":"test123/test.mp4"}]}' \
  response.json
```