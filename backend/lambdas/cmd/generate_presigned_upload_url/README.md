# Generate Presigned Upload URL Lambda

Lambda function that generates presigned S3 URLs for uploading video files. Works with three S3 buckets:

| Bucket Name           | Purpose                        |
|----------------------|--------------------------------|
| `raw-storage`        | Direct uploads (presigned URL) |


## Local Development with LocalStack

### 1. Start LocalStack and Create Resources
Use the project's build script:

```bash
./build.sh
```

This script will:
- Start LocalStack via docker-compose
- Create all required S3 buckets (`raw-storage`, `transcoded-storage`, `thumbnails-storage`)
- Build and deploy the Lambda function
- Set up DynamoDB tables and SQS queues
- Start other services

### 2. Manual Setup (Alternative)

#### Start LocalStack
```bash
docker-compose up -d
```

#### Create S3 Buckets
```bash
awslocal s3 mb s3://raw-storage
awslocal s3 mb s3://transcoded-storage
awslocal s3 mb s3://thumbnails-storage
```

#### Build and Deploy Lambda
```bash
cd backend/storage/cmd/generate_presigned_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main

awslocal lambda create-function \
  --function-name generate-presigned-url \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role
```

#### Set Environment Variables
```bash
awslocal lambda update-function-configuration \
  --function-name generate-presigned-url \
  --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

### 3. Test the Lambda Function

#### Generate a Presigned URL
```bash
awslocal lambda invoke \
  --function-name generate-presigned-url \
  --region us-east-1 \
  --payload '{"body": "{\"fileName\": \"test_video.mp4\"}"}' \
  --cli-binary-format raw-in-base64-out \
  response.json

cat response.json | jq '.body' | jq -r | jq '.url'
```

#### Upload a File Using the Presigned URL
```bash
curl -X PUT -T test_video.mp4 "<PASTE_URL_HERE>"
```

## Environment Variables
- `BUCKET_NAME`: Name of the S3 bucket for uploads (default: `raw-storage`)
- `BUCKET_REGION`: AWS region (default: `us-east-1`)
- `AWS_ENDPOINT`: Custom endpoint for S3 (default: `http://localhost:4566` for LocalStack)


## Example: Generate a Test Video
```bash
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 -pix_fmt yuv420p test_video.mp4
```


