# Generate Presigned Upload URL Lambda

Lambda function that generates presigned S3 URLs for uploading video files. Used for direct browser or client uploads to the `raw-storage` S3 bucket.

## Bucket Usage

| Bucket Name        | Purpose                        |
|--------------------|--------------------------------|
| `raw-storage`      | Direct uploads (presigned URL) |

---

## Environment Variables

| Variable         | Description                                 | Default                   |
|------------------|---------------------------------------------|---------------------------|
| `BUCKET_NAME`    | S3 bucket for uploads                       | `raw-storage`             |
| `BUCKET_REGION`  | AWS region for the bucket                   | `us-east-1`               |
| `AWS_ENDPOINT`   | Custom AWS endpoint (for LocalStack)        | `http://localhost:4566`   |

---

## Local Development with LocalStack

### Option 1: Use Build Script

Run the project’s setup script:

```bash
./build.sh
```

This will:
- Start LocalStack via Docker Compose
- Create required S3 buckets: `raw-storage`, `transcoded-storage`, `thumbnails-storage`
- Deploy the Lambda function
- Set up related services (DynamoDB, SQS, etc.)

---

### Option 2: Manual Setup

#### 1. Start LocalStack

```bash
docker-compose up -d
```

#### 2. Create S3 Buckets

```bash
awslocal s3 mb s3://raw-storage
awslocal s3 mb s3://transcoded-storage
awslocal s3 mb s3://thumbnails-storage
```

#### 3. Build and Deploy Lambda

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

#### 4. Set Environment Variables

```bash
awslocal lambda update-function-configuration \
  --function-name generate-presigned-url \
  --environment "Variables={BUCKET_NAME=raw-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

---

## Testing

### 1. Generate a Presigned URL

```bash
awslocal lambda invoke \
  --function-name generate-presigned-url \
  --region us-east-1 \
  --payload '{"body": "{\"fileName\": \"test_video.mp4\"}"}' \
  --cli-binary-format raw-in-base64-out \
  response.json

cat response.json | jq '.body' | jq -r | jq '.url'
```

### 2. Upload a File with the URL

```bash
curl -X PUT -T test_video.mp4 "<PASTE_URL_HERE>"
```

---

## Utility: Generate a Test Video with FFmpeg

```bash
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 -pix_fmt yuv420p test_video.mp4
```