# Trigger Transcode Job Lambda

Lambda function that triggers video transcoding by sending messages to an SQS queue. Accepts input via HTTP POST and pushes jobs to the transcoder.

## Features

- Accepts requests with `assetId`, `videoType`, and `format`
- Sends structured SQS messages to the transcoder job queue
- Supports HLS and DASH job types
- Returns success or error response

---

## Request Format

```json
{
  "assetId": "asset-123",
  "videoType": "main",
  "format": "hls"
}
```

---

## Response Format

```json
{
  "message": "Transcoding job triggered successfully",
  "jobType": "transcode-hls"
}
```

---

## Environment Variables

| Variable              | Description                            | Default         |
|-----------------------|----------------------------------------|-----------------|
| `TRANSCODER_QUEUE_URL`| SQS queue URL for transcoding jobs     | (required)      |
| `AWS_REGION`          | AWS region                             | `us-east-1`     |
| `AWS_ENDPOINT`        | Custom AWS endpoint (e.g. LocalStack)  | (optional)      |

---

## Build and Deploy (Manual)

```bash
cd backend/lambdas/cmd/trigger_transcode_job
go mod tidy
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -j function.zip main
```

You can then deploy the zipped function to AWS or LocalStack using the CLI:

```bash
awslocal lambda create-function \
  --function-name trigger-transcode-job \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role
```

---
