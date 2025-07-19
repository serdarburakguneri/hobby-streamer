# Generate Image Upload URL Lambda

Lambda function that generates presigned S3 URLs for uploading image files to `images-storage` bucket.

## Image Types

THUMBNAIL, POSTER, BANNER, HERO, LOGO, SCREENSHOT, BEHIND_THE_SCENES, INTERVIEW

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BUCKET_NAME` | S3 bucket for images | `images-storage` |
| `BUCKET_REGION` | AWS region | `us-east-1` |
| `AWS_ENDPOINT` | Custom endpoint (LocalStack) | `http://localhost:4566` |

## Local Development

### Quick Setup
```bash
./build.sh  # Starts LocalStack, creates buckets, deploys Lambda
```

### Manual Setup
```bash
# Start LocalStack
docker-compose up -d

# Create bucket
awslocal s3 mb s3://images-storage

# Build and deploy
cd backend/lambdas/cmd/generate_image_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main

awslocal lambda create-function \
  --function-name generate-image-upload-url \
  --runtime go1.x --handler main --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role

# Set environment
awslocal lambda update-function-configuration \
  --function-name generate-image-upload-url \
  --environment "Variables={BUCKET_NAME=images-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

## Testing

```bash
# Generate presigned URL
awslocal lambda invoke \
  --function-name generate-image-upload-url \
  --payload '{"body": "{\"fileName\": \"poster.jpg\", \"assetId\": \"123\", \"imageType\": \"poster\"}"}' \
  response.json

# Upload image
curl -X PUT -T poster.jpg "$(cat response.json | jq -r '.body' | jq -r '.url')"
```

## Storage Pattern

`images-storage/<assetId>/<imageType>/<filename>`

**Examples:** `images-storage/123/poster/poster.jpg`, `images-storage/123/thumbnail/thumb.jpg` 