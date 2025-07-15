# Trigger Transcode Job Lambda

This Lambda function triggers transcoding jobs by sending SQS messages to the transcoder queue.

## Purpose

When a user wants to transcode a video to HLS or DASH format, this Lambda function:
1. Receives a POST request with assetId, videoType, and format
2. Sends an SQS message to the transcoder queue
3. Returns success/failure response

## Request Format

```json
{
  "assetId": "asset-123",
  "videoType": "main",
  "format": "hls"
}
```

## Response Format

```json
{
  "message": "Transcoding job triggered successfully",
  "jobType": "transcode-hls"
}
```

## Environment Variables

- `TRANSCODER_QUEUE_URL`: SQS queue URL for transcoder jobs
- `AWS_REGION`: AWS region (default: us-east-1)
- `AWS_ENDPOINT`: Custom endpoint for LocalStack (optional)

## Build and Deploy

```bash
cd backend/lambdas/cmd/trigger_transcode_job
go mod tidy
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main
```

## API Gateway Integration

This Lambda should be integrated with API Gateway to provide HTTP endpoints for triggering transcoding jobs. 