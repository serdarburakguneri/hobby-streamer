# Kafka Architecture

Event streaming architecture using Apache Kafka with CloudEvents format.

## Architecture Overview

```
┌─────────────────┐                           ┌─────────────────┐
│   Upload Lambda │                           │   Asset Manager │
│   (S3 Trigger)  │                           │   (Producer)    │
└─────────┬───────┘                           └─────────┬───────┘
          │                                             │
          ▼                                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Apache Kafka                             │
└─────────────────────────────────────────────────────────────────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Asset Manager │    │    Transcoder   │    │ Content Analyzer│
│   (Consumer)    │    │   (Consumer)    │    │   (Consumer)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Event Flow

### Video Upload Flow
1. **Upload Lambda** (S3 trigger) → `raw-video-uploaded`
2. **Asset Manager** (consumer) → adds video to asset → `analyze.job.requested`
3. **Transcoder** (consumer) → analyzes video → `analyze.job.completed`
4. **Asset Manager** (consumer) → saves video metadata

### HLS Transcoding Flow
1. **CMS UI** → **Asset Manager** → `hls.job.requested`
2. **Asset Manager** (consumer) → saves HLS video to asset
3. **Transcoder** (consumer) → transcodes to HLS → `hls.job.completed`
4. **Asset Manager** (consumer) → updates HLS video details


### DASH Transcoding Flow
1. **CMS UI** → **Asset Manager** → `dash.job.requested`
2. **Asset Manager** (consumer) → saves DASH video to asset
3. **Transcoder** (consumer) → transcodes to DASH → `dash.job.completed`
4. **Asset Manager** (consumer) → updates DASH video details


## Topics

- **`raw-video-uploaded`** - Raw video upload notifications
- **`analyze.job.requested`** - Video analysis job requests
- **`hls.job.requested`** - HLS transcoding job requests
- **`dash.job.requested`** - DASH transcoding job requests
- **`analyze.job.completed`** - Video analysis job completions
- **`hls.job.completed`** - HLS transcoding job completions
- **`dash.job.completed`** - DASH transcoding job completions


## Consumer Groups

- **`asset-manager-group`** - Processes uploads and job completions
- **`transcoder-group`** - Processes video analysis and transcoding jobs

## CloudEvents Format

All events follow CloudEvents 1.0 specification:

```json
{
  "specversion": "1.0",
  "id": "uuid-v4",
  "source": "upload-lambda",
  "type": "raw-video-uploaded",
  "datacontenttype": "application/json",
  "time": "2024-01-01T12:00:00Z",
  "data": {
    "assetId": "asset-123",
    "videoId": "video-456",
    "storageLocation": "s3://bucket/key",
    "filename": "video.mp4",
    "size": 1048576,
    "contentType": "video/mp4"
  }
}
```

## Quick Start

### View Kafka UI
```bash
open http://localhost:8086
```

### View Kibana
```bash
open http://localhost:5601
```

### Topic Management
```bash
# List topics
docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Describe topic
docker exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic raw-video-uploaded

# Consumer group status
docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list
```
