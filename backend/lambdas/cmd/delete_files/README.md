# Delete Files Lambda

Lambda function that handles deletion of files from S3 buckets. Used to clean up S3 files when assets are deleted.

## Features

- Deletes files from multiple S3 buckets
- Handles raw video files, transcoded files (HLS/DASH), and thumbnails
- Returns detailed response with success/error information

## API

### Request Format
```json
{
  "assetId": "asset123",
  "files": [
    {
      "bucket": "raw-storage",
      "key": "asset123/main_1234567890.mp4"
    },
    {
      "bucket": "transcoded-storage", 
      "key": "asset123/main_1234567890.m3u8"
    },
    {
      "bucket": "thumbnails-storage",
      "key": "asset123/thumbnail_1234567890.jpg"
    }
  ]
}
```

### Response Format
```json
{
  "message": "Deleted 3 files for asset asset123",
  "deleted": [
    {
      "bucket": "raw-storage",
      "key": "asset123/main_1234567890.mp4"
    }
  ],
  "errors": [
    "Failed to delete asset123/thumbnail_1234567890.jpg from thumbnails-storage: NoSuchKey"
  ]
}
```

## Local Development

### Build
```bash
cd backend/storage/cmd/delete_files
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main
```

### Deploy to LocalStack
```bash
awslocal lambda create-function \
  --function-name delete-files \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --environment "Variables={AWS_ENDPOINT=http://localhost:4566,AWS_REGION=eu-west-1}" \
  --region eu-west-1
```

### Test
```bash
awslocal lambda invoke \
  --function-name delete-files \
  --payload '{"assetId":"test123","files":[{"bucket":"raw-storage","key":"test123/test.mp4"}]}' \
  response.json
```
