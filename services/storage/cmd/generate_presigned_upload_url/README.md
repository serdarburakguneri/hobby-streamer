# Presigned Upload URL Lambda

This Lambda function generates a presigned S3 URL for uploading video files. It is designed to work with three S3 buckets in your domain:

| Bucket Name           | Purpose                        |
|----------------------|--------------------------------|
| my-raw-bucket        | Direct uploads (presigned URL) |
| my-processing-bucket | Intermediate processing        |
| my-public-bucket     | Final, public assets           |

## Local Development with LocalStack

### 1. Start LocalStack
```
localstack start
```

### 2. Create S3 Buckets
```
awslocal s3 mb s3://my-raw-bucket
awslocal s3 mb s3://my-processing-bucket
awslocal s3 mb s3://my-public-bucket
```

### 3. Build and Deploy Lambda
```
GOOS=linux GOARCH=amd64 go build -o main main.go
zip function.zip main
awslocal lambda create-function \
  --function-name generate-presigned-url \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role
```

### 4. Set Environment Variables
```
awslocal lambda update-function-configuration \
  --function-name generate-presigned-url \
  --environment "Variables={BUCKET_NAME=my-raw-bucket,BUCKET_REGION=us-east-1,AWS_ENDPOINT=http://localhost:4566}"
```

### 5. Invoke Lambda to Get a Presigned URL
```
awslocal lambda invoke \
  --function-name generate-presigned-url \
  --region us-east-1 \
  --payload '{"body": "{\"fileName\": \"test_video.mp4\"}"}' \
  --cli-binary-format raw-in-base64-out \
  response.json
cat response.json | jq '.body' | jq -r | jq '.url'
```

### 6. Upload a File Using the Presigned URL
```
curl -X PUT -T test_video.mp4 "<PASTE_URL_HERE>"
```

## Environment Variables
- `BUCKET_NAME`: Name of the S3 bucket for uploads (e.g., `my-raw-bucket`)
- `BUCKET_REGION`: AWS region (default: `eu-north-1`)
- `AWS_ENDPOINT`: (Optional) Custom endpoint for S3, e.g., `http://localhost:4566` for LocalStack

## Notes
- The Lambda currently only supports uploading to the raw bucket. Support for multiple buckets (processing, public) can be added as needed.
- CORS is enabled for all origins by default. Adjust as needed for production.
- For production, deploy to AWS and remove the `AWS_ENDPOINT` variable.

## Example: Generate a Test Video
```
ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 -pix_fmt yuv420p test_video.mp4
```




