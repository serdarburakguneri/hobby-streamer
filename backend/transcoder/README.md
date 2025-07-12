# Transcoder Service

The Transcoder service processes video jobs from an SQS queue, including analysis and transcoding to HLS and DASH formats. It is designed to work as a background worker in the hobby-streamer project.

## Features
- Consumes jobs from an SQS queue
- Runs ffmpeg-based analysis and transcoding jobs (HLS, DASH)
- Extensible job registry for new job types

## Requirements
- Go 1.22+
- ffmpeg (must be installed and available in PATH)
- LocalStack (for local AWS emulation)

## Environment Variables
- `SQS_QUEUE_URL`: The SQS queue URL to consume jobs from (required)

## Running Locally

### 1. Start LocalStack and create the SQS queue (see project root `build.sh`)

### 2. Run the worker:
```sh
cd backend/transcoder
SQS_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs go run ./cmd/worker/main.go
```

## Job Types
- analyze: Runs ffmpeg to analyze a video file (input: `{ "input": "path/to/file" }`)
- transcode-hls: Transcodes a video to HLS format (input: `{ "input": "path/to/file", "output": "path/to/output.m3u8" }`)
- transcode-dash: Transcodes a video to DASH format (input: `{ "input": "path/to/file", "output": "path/to/output.mpd" }`)

## Notes
- The service will exit if `SQS_QUEUE_URL` is not set.
- See the `internal/job/` directory for job runner implementations and payload formats.
- Designed to work with other services in the hobby-streamer project. 