# Transcoder Service

A background worker that processes video jobs from an SQS queue. Uses FFmpeg to perform video analysis and generate HLS and DASH outputs for playback.

---

## Features

- Listens to SQS for incoming video job messages
- Supports:
  - Metadata analysis (duration, format, streams)
  - HLS transcoding
  - DASH transcoding
- Built with the shared SQS client package
- Includes retry logic and structured logging for failed jobs

---

## Requirements

- Go 1.22+
- FFmpeg installed and accessible via `$PATH`
- LocalStack running with the correct SQS queue configured

---

## Running Locally

### 1. Ensure queues are available

Use the root-level `./local/build.sh` script to spin up LocalStack and initialize required queues.

### 2. Start the worker

```bash
cd backend/transcoder
TRANSCODER_QUEUE_URL=http://localhost:4566/000000000000/transcoder-jobs \
go run ./cmd/worker/main.go
```

---

## Supported Job Types

### `analyze`

Runs FFmpeg to extract video metadata (streams, duration, etc.).

**Example:**

```json
{
  "input": "path/to/video.mp4",
  "assetId": "asset-id",
  "videoType": "main"
}
```

---

### `transcode-hls`

Transcodes a video to HLS format (`.m3u8` + segments).

**Example:**

```json
{
  "input": "path/to/video.mp4",
  "output": "output/path/playlist.m3u8"
}
```

---

### `transcode-dash`

Transcodes a video to MPEG-DASH format (`.mpd` + segments).

**Example:**

```json
{
  "input": "path/to/video.mp4",
  "output": "output/path/manifest.mpd"
}
```

---

> ⚠️ This service is intended for local testing and development only. FFmpeg settings and job handling may evolve as use cases expand.