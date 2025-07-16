# Trigger Transcode Job Lambda

Lambda function that triggers transcoding jobs by sending SQS messages to the transcoder queue.

## Features

- Receives POST requests with assetId, videoType, and format
- Sends SQS messages to the transcoder queue
- Returns success/failure response

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
