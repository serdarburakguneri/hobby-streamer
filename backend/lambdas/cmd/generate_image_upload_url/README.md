# Generate Image Upload URL Lambda

Lambda for presigned S3 URLs to upload images to `images-storage`.

## Features
Presigned S3 upload URLs, supports image types (THUMBNAIL, POSTER, BANNER, HERO, LOGO, SCREENSHOT, BEHIND_THE_SCENES, INTERVIEW), LocalStack support, easy local testing, simple environment config.

## Env Vars
BUCKET_NAME (default: images-storage), BUCKET_REGION (default: us-east-1), AWS_ENDPOINT (default: http://localhost:4566)

## Quick Setup
```bash
./build.sh
```

## Manual Setup
```bash
docker-compose up -d
awslocal s3 mb s3://images-storage
cd backend/lambdas/cmd/generate_image_upload_url
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main
awslocal lambda create-function --function-name generate-image-upload-url --runtime go1.x --handler main --zip-file fileb://function.zip --role arn:aws:iam::000000000000:role/lambda-role
awslocal lambda update-function-configuration --function-name generate-image-upload-url --environment "Variables={BUCKET_NAME=images-storage,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

## Testing
```bash
awslocal lambda invoke --function-name generate-image-upload-url --payload '{"body": "{\"fileName\": \"poster.jpg\", \"assetId\": \"123\", \"imageType\": \"poster\"}"}' response.json
curl -X PUT -T poster.jpg "$(cat response.json | jq -r '.body' | jq -r '.url')"
```

## Storage Pattern
`images-storage/<assetId>/<imageType>/<filename>`
Examples: `images-storage/123/poster/poster.jpg`, `images-storage/123/thumbnail/thumb.jpg` 