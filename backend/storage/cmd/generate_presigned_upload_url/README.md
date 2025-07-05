# Storage Service

This Lambda function generates a presigned S3 URL for uploading video files. It is designed to work with three S3 buckets in the hobby-streamer project:

| Bucket Name           | Purpose                        |
|----------------------|--------------------------------|
| `raw-storage`        | Direct uploads (presigned URL) |
| `transcoded-storage` | Processed video files (HLS/DASH) |
| `thumbnails-storage` | Video thumbnails and images    |

## Local Development with LocalStack

### 1. Start LocalStack and Create Resources
The easiest way to set up the entire environment is to use the project's build script:

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

If you prefer to set up manually:

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

## Bucket Usage

### raw-storage
- **Purpose**: Initial video uploads from clients
- **Content**: Raw video files uploaded via presigned URLs
- **Access**: Direct upload via presigned URLs

### transcoded-storage
- **Purpose**: Store processed video files
- **Content**: HLS playlists, DASH manifests, and video segments
- **Access**: Read-only for video streaming

### thumbnails-storage
- **Purpose**: Store video thumbnails and images
- **Content**: Generated thumbnails, poster images, and metadata images
- **Access**: Read-only for UI display

## Integration with Other Services

This storage service works with:
- **Asset Manager Service**: Manages metadata and references to stored files
- **Transcoder Service**: Processes raw videos and stores results in `transcoded-storage`
- **Frontend Applications**: Upload raw videos and stream processed content

## Notes
- The Lambda currently supports uploading to the `raw-storage` bucket only
- CORS is enabled for all origins by default (adjust for production)
- For production deployment, remove the `AWS_ENDPOINT` variable
- The service is designed to work with LocalStack for local development

## Example: Generate a Test Video
```bash
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 -pix_fmt yuv420p test_video.mp4
```

## Troubleshooting

### Common Issues
1. **Lambda not found**: Ensure you've built and deployed the function
2. **Bucket not found**: Check that LocalStack is running and buckets are created
3. **Permission errors**: Verify the Lambda role has S3 permissions
4. **CORS issues**: Check that the presigned URL is being used correctly

### Debug Commands
```bash
# Check if LocalStack is running
curl http://localhost:4566/health

# List S3 buckets
awslocal s3 ls

# Check Lambda function
awslocal lambda get-function --function-name generate-presigned-url
```




