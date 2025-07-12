# Delete Asset Files Lambda

This Lambda function handles deletion of asset files from S3 buckets. It is designed to work with the hobby-streamer project's storage buckets.

## Purpose

When an asset is deleted from the asset-manager, this function is called to clean up all associated S3 files including:
- Raw video files
- Transcoded video files (HLS/DASH)
- Thumbnail images

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
cd backend/storage/cmd/delete_asset_files
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main
```

### Deploy to LocalStack
```bash
awslocal lambda create-function \
  --function-name delete-asset-files \
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
  --function-name delete-asset-files \
  --payload '{"assetId":"test123","files":[{"bucket":"raw-storage","key":"test123/test.mp4"}]}' \
  response.json
```

## Integration

This function should be called by the asset-manager service when an asset is deleted, after the database record is removed but before the response is sent to the client. 