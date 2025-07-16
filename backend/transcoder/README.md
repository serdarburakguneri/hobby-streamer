# Transcoder Service

Background worker that processes video jobs from an SQS queue, including analysis and transcoding to HLS and DASH formats.

## Features
- Consumes jobs from an SQS queue using the shared SQS package
- Runs ffmpeg-based analysis and transcoding jobs (HLS, DASH)
- Direct job handler registration without dispatcher abstraction

## Requirements
- Go 1.22+
- ffmpeg (must be installed and available in PATH)
- LocalStack (for local AWS emulation)

## Environment Variables
- `TRANSCODER_QUEUE_URL`: The SQS queue URL to consume jobs from (required)
- `ANALYZE_QUEUE_URL`: Optional SQS queue URL for analyze completion messages
- `AWS_ENDPOINT`: Custom endpoint for AWS services (default: `http://localstack:4566` for LocalStack)
- `AWS_REGION`: AWS region (default: `us-east-1`)
- `AWS_ACCESS_KEY_ID`: AWS access key (default: `test` for LocalStack)
- `AWS_SECRET_ACCESS_KEY`: AWS secret key (default: `test` for LocalStack)
- `LOG_FORMAT`: Log format (default: `text`)

## Running Locally

### 1. Start LocalStack and create the SQS queue (see project root `build.sh`)

### 2. Run the worker:
```sh
cd backend/transcoder
TRANSCODER_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs go run ./cmd/worker/main.go
```

## Job Types
- analyze: Runs ffmpeg to analyze a video file (input: `{ "input": "path/to/file", "assetId": "asset-id", "videoType": "type" }`)
- transcode-hls: Transcodes a video to HLS format (input: `{ "input": "path/to/file", "output": "path/to/output.m3u8" }`)
- transcode-dash: Transcodes a video to DASH format (input: `{ "input": "path/to/file", "output": "path/to/output.mpd" }`)
