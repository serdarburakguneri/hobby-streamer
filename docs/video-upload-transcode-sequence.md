# Video Upload and Transcoding Sequence

How video upload and transcoding works in hobby-streamer.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant User as User
    participant CMS as CMS Frontend
    participant API as API Gateway
    participant Lambda as Lambda Functions
    participant S3 as S3 Storage
    participant AssetMgr as Asset Manager
    participant Neo4j as Neo4j Database
    participant SQS as SQS Queue
    participant Transcoder as Transcoder Worker

    Note over User,Transcoder: Upload Flow
    User->>CMS: Select video file
    CMS->>API: Get presigned URL
    API->>Lambda: Generate presigned URL
    Lambda-->>CMS: Return presigned URL
    CMS->>S3: Upload video directly
    CMS->>AssetMgr: Save video metadata
    AssetMgr->>Neo4j: Create asset record
    AssetMgr->>SQS: Trigger analysis job

    Note over User,Transcoder: Analysis Flow
    SQS->>Transcoder: Process analysis job
    Transcoder->>S3: Download video
    Transcoder->>SQS: Send analysis complete
    SQS->>AssetMgr: Update video status to ready

    Note over User,Transcoder: Transcoding Flow
    User->>CMS: Trigger HLS/DASH transcoding
    CMS->>AssetMgr: addVideo GraphQL mutation
    AssetMgr->>Neo4j: Create new video variant
    AssetMgr->>SQS: Send transcode job
    SQS->>Transcoder: Process transcode job
    Transcoder->>S3: Download raw video
    Transcoder->>S3: Upload transcoded files
    Transcoder->>SQS: Send transcode complete
    SQS->>AssetMgr: Update video status to ready
```

## Storage Structure
S3: content-east/{assetId}/source/video.mp4, content-east/{assetId}/hls/1080p/playlist.m3u8, content-east/{assetId}/dash/1080p/manifest.mpd

## Video Model
```json
{
  "id": "video-123",
  "type": "MAIN",
  "format": "raw|hls|dash",
  "storageLocation": {
    "bucket": "content-east",
    "key": "asset-456/source/video.mp4"
  },
  "width": 1920,
  "height": 1080,
  "duration": 120.5,
  "status": "pending|analyzing|transcoding|ready|failed",
  "streamInfo": {
    "cdnPrefix": "http://localhost:8083/cdn",
    "url": "http://localhost:8083/cdn/asset-456/hls/1080p/playlist.m3u8"
  }
}
```

## Status Flow
1. Upload: status "ready"
2. Transcode: user triggers via GraphQL, new video record with status "transcoding" → "ready"
3. Multiple formats: raw, HLS, DASH variants
4. Retry: transcoder retries failed jobs
5. Error handling: validation errors discarded, others retried





