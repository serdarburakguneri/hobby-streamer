# Generate Image Upload URL Lambda

Lambda function that generates presigned S3 URLs for uploading image files. Used for direct browser or client uploads to the `images-storage` S3 bucket.

## Bucket Usage

| Bucket Name        | Purpose                        |
|--------------------|--------------------------------|
| `images-storage`   | Direct image uploads (presigned URL) |

## Image Types

The following image types are supported:

- **THUMBNAIL** - Small preview image (e.g., 320x180)
- **POSTER** - Main promotional image (e.g., 1920x1080)
- **BANNER** - Wide banner image (e.g., 1920x400)
- **HERO** - Large hero image for featured content (e.g., 2560x1440)
- **LOGO** - Brand/studio logo
- **SCREENSHOT** - Screenshots from the content
- **BEHIND_THE_SCENES** - Behind the scenes photos
- **INTERVIEW** - Interview photos

## Environment Variables

| Variable         | Description                                 | Default                   |
|------------------|---------------------------------------------|---------------------------|
| `BUCKET_NAME`    | S3 bucket for image uploads                 | `images-storage`          |
| `BUCKET_REGION`  | AWS region for the bucket                   | `us-east-1`               |
| `AWS_ENDPOINT`   | Custom AWS endpoint (for LocalStack)        | `http://localhost:4566`   |

## Local Development with LocalStack

### Option 1: Use Build Script

Run the project's setup script:

```bash
./build.sh
```

This will:
- Start LocalStack via Docker Compose
- Create required S3 buckets including `images-storage`
- Deploy the Lambda function
- Set up related services

### Option 2: Manual Setup

#### 1. Start LocalStack

```bash
docker-compose up -d
```

#### 2. Create S3 Buckets

```bash
awslocal s3 mb s3://images-storage
```

#### 3. Build and Deploy Lambda

```bash
cd backend/lambdas/cmd/generate_image_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main

awslocal lambda create-function \
  --function-name generate-image-upload-url \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role
```

#### 4. Set Environment Variables

```bash
awslocal lambda update-function-configuration \
  --function-name generate-image-upload-url \
  --environment "Variables={BUCKET_NAME=images-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

## Testing

### 1. Generate a Presigned URL

```bash
awslocal lambda invoke \
  --function-name generate-image-upload-url \
  --region us-east-1 \
  --payload '{"body": "{\"fileName\": \"poster.jpg\", \"assetId\": \"123\", \"imageType\": \"poster\"}"}' \
  --cli-binary-format raw-in-base64-out \
  response.json

cat response.json | jq '.body' | jq -r | jq '.url'
```

### 2. Upload an Image with the URL

```bash
curl -X PUT -T poster.jpg "<PASTE_URL_HERE>"
```

## API Usage

### Request

```json
{
  "fileName": "poster.jpg",
  "assetId": "123",
  "imageType": "poster"
}
```

### Response

```json
{
  "url": "https://s3.amazonaws.com/images-storage/123/poster/poster.jpg?..."
}
```

## Storage Path Pattern

Images are stored using the following pattern:
```
images-storage/<assetId>/<imageType>/<filename>
```

Examples:
- `images-storage/123/poster/poster.jpg`
- `images-storage/123/thumbnail/thumb.jpg`
- `images-storage/123/banner/banner.png` 