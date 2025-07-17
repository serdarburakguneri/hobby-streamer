# Transcoder Service

A background worker that consumes video processing jobs from an SQS queue. Handles video analysis and transcoding using FFmpeg for HLS and DASH outputs.

## Features

- Listens to an SQS queue for job messages
- Runs FFmpeg commands for:
  - Video analysis (metadata, duration, streams)
  - HLS and DASH transcoding
- Uses the shared SQS package for job consumption
- Handles failures with retries and logging

## Requirements

- Go 1.22+
- FFmpeg installed and available in PATH
- LocalStack (for SQS emulation)


## Running Locally

### 1. Start LocalStack and ensure required queues exist

You can use the project’s `./local/build.sh` script to initialize queues.

### 2. Start the worker

```bash
cd backend/transcoder
TRANSCODER_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs \
go run ./cmd/worker/main.go
```

## Supported Job Types

### analyze

Runs FFmpeg to extract metadata and stream info.

**Example Payload**
```json
{
  "input": "path/to/video.mp4",
  "assetId": "asset-id",
  "videoType": "main"
}
```

### transcode-hls

Transcodes a video to HLS format.

**Example Payload**
```json
{
  "input": "path/to/video.mp4",
  "output": "output/path/playlist.m3u8"
}
```

### transcode-dash

Transcodes a video to DASH format.

**Example Payload**
```json
{
  "input": "path/to/video.mp4",
  "output": "output/path/manifest.mpd"
}
```