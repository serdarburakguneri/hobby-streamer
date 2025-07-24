# Transcoder Service

Background worker for video jobs from SQS, using FFmpeg for analysis and HLS/DASH transcoding.

## Features
SQS job processing, metadata analysis, HLS/DASH transcoding, retry logic, structured logging.

## Requirements
Go 1.22+, FFmpeg in $PATH, LocalStack with SQS queues.

## Running
```bash
./local/build.sh # Start queues
cd backend/transcoder
TRANSCODER_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs \
go run ./cmd/worker/main.go
```

## Job Types
- `analyze`: Extract video metadata. Input: `{ "input": "path/to/video.mp4", "assetId": "asset-id", "videoType": "main" }`
- `transcode-hls`: Transcode to HLS. Input: `{ "input": "path/to/video.mp4", "output": "output/path/playlist.m3u8" }`
- `transcode-dash`: Transcode to MPEG-DASH. Input: `{ "input": "path/to/video.mp4", "output": "output/path/manifest.mpd" }`