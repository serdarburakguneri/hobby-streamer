# Generate Video Upload URL Lambda

Lambda for presigned S3 URLs to upload videos to `raw-storage`.

## Features
Presigned S3 upload URLs, LocalStack support, easy local testing, simple environment config.

## Env Vars
BUCKET_NAME (default: raw-storage), BUCKET_REGION (default: us-east-1), AWS_ENDPOINT (default: http://localhost:4566)

## Quick Setup
```bash
./build.sh
```

## Manual Setup
```bash
docker-compose up -d
awslocal s3 mb s3://raw-storage s3://transcoded-storage s3://images-storage
cd backend/lambdas/cmd/generate_video_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main
awslocal lambda create-function --function-name generate-video-upload-url --runtime go1.x --handler main --zip-file fileb://function.zip --role arn:aws:iam::000000000000:role/lambda-role
awslocal lambda update-function-configuration --function-name generate-video-upload-url --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

## Testing
```bash
awslocal lambda invoke --function-name generate-video-upload-url --payload '{"body": "{\"fileName\": \"test.mp4\", \"assetId\": \"123\", \"videoType\": \"main\"}"}' response.json
curl -X PUT -T test.mp4 "$(cat response.json | jq -r '.body' | jq -r '.url')"
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 test.mp4
```