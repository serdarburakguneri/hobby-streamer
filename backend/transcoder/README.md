# Transcoder Service

Worker for video analysis and HLS/DASH transcoding with FFmpeg.

## Features
Analyze metadata, HLS/DASH transcode, retries, structured logs.

## Run
```bash
./local/build.sh
cd backend/transcoder && go run ./cmd/worker/main.go
```