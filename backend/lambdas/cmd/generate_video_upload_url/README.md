# Generate Video Upload URL Lambda

Lambda function that generates presigned S3 URLs for uploading video files to `raw-storage` bucket.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BUCKET_NAME` | S3 bucket for uploads | `raw-storage` |
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

# Create buckets
awslocal s3 mb s3://raw-storage s3://transcoded-storage s3://images-storage

# Build and deploy
cd backend/lambdas/cmd/generate_video_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main

awslocal lambda create-function \
  --function-name generate-video-upload-url \
  --runtime go1.x --handler main --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role

# Set environment
awslocal lambda update-function-configuration \
  --function-name generate-video-upload-url \
  --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

## Testing

```bash
# Generate presigned URL
awslocal lambda invoke \
  --function-name generate-video-upload-url \
  --payload '{"body": "{\"fileName\": \"test.mp4\", \"assetId\": \"123\", \"videoType\": \"main\"}"}' \
  response.json

# Upload file
curl -X PUT -T test.mp4 "$(cat response.json | jq -r '.body' | jq -r '.url')"

# Generate test video
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 test.mp4
```