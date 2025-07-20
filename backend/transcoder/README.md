# Transcoder Service

Background worker that processes video jobs from SQS queue. Uses FFmpeg for video analysis and HLS/DASH transcoding.

## Features

SQS job processing, metadata analysis, HLS/DASH transcoding, retry logic, structured logging.

## Requirements

Go 1.22+, FFmpeg in `$PATH`, LocalStack with SQS queues.

## Running

```bash
# Start queues first
./local/build.sh

# Start worker
cd backend/transcoder
TRANSCODER_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs \
go run ./cmd/worker/main.go
```

## Job Types

### `analyze`
Extract video metadata (streams, duration).  
**Input:** `{"input": "path/to/video.mp4", "assetId": "asset-id", "videoType": "main"}`

### `transcode-hls`
Transcode to HLS format (`.m3u8` + segments).  
**Input:** `{"input": "path/to/video.mp4", "output": "output/path/playlist.m3u8"}`

### `transcode-dash`
Transcode to MPEG-DASH format (`.mpd` + segments).  
**Input:** `{"input": "path/to/video.mp4", "output": "output/path/manifest.mpd"}`

> ⚠️ Local testing and development only. FFmpeg settings may evolve.