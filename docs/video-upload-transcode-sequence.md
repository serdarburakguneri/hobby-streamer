# Video Upload and Transcoding Sequence Diagram

This document provides a detailed UML sequence diagram showing the complete flow of video uploading and transcoding in the hobby-streamer project.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant User as User (CMS UI)
    participant CMS as CMS Frontend
    participant API as API Gateway
    participant Lambda as Presigned URL Lambda
    participant S3 as S3 Storage
    participant AssetMgr as Asset Manager
    participant Neo4j as Neo4j Database
    participant TranscodeLambda as Transcode Lambda
    participant SQS as SQS Queue
    participant Transcoder as Transcoder Worker
    participant FFmpeg as FFmpeg
    participant Consumer as Completion Consumers

    Note over User,Consumer: Video Upload Flow
    User->>CMS: Select video file
    CMS->>API: POST /upload (get presigned URL)
    API->>Lambda: Invoke generate-presigned-url
    Lambda->>S3: Generate presigned URL
    S3-->>Lambda: Return presigned URL
    Lambda-->>API: Return presigned URL
    API-->>CMS: Return presigned URL
    CMS->>S3: PUT video file (direct upload)
    S3-->>CMS: Upload success
    CMS->>AssetMgr: GraphQL addVideo mutation
    AssetMgr->>Neo4j: Save asset with video metadata
    Neo4j-->>AssetMgr: Asset saved
    AssetMgr->>SQS: Send analyze job message
    AssetMgr-->>CMS: Return updated asset
    CMS-->>User: Show upload success

    Note over User,Consumer: Video Analysis Flow
    SQS->>Transcoder: Consume analyze job
    Transcoder->>S3: Download video file
    S3-->>Transcoder: Return video file
    Transcoder->>FFmpeg: Analyze video (ffmpeg -i)
    FFmpeg-->>Transcoder: Analysis results
    Transcoder->>SQS: Send analyze-completed message
    SQS->>Consumer: Consume analyze-completed
    Consumer->>AssetMgr: HandleAnalyzeCompletion
    AssetMgr->>Neo4j: Update video status to "ready"
    Neo4j-->>AssetMgr: Status updated
    AssetMgr-->>Consumer: Completion handled

    Note over User,Consumer: Manual Transcoding Flow
    User->>CMS: Click "Transcode to HLS/DASH"
    CMS->>API: POST /transcode
    API->>TranscodeLambda: Invoke trigger-transcode-job
    TranscodeLambda->>SQS: Send transcode-hls/dash message
    TranscodeLambda-->>API: Return success
    API-->>CMS: Return success
    CMS-->>User: Show transcoding started

    Note over User,Consumer: Transcoding Processing Flow
    SQS->>Transcoder: Consume transcode job
    Transcoder->>S3: Download raw video
    S3-->>Transcoder: Return video file
    Transcoder->>FFmpeg: Transcode to HLS/DASH
    FFmpeg-->>Transcoder: Transcoded files
    Transcoder->>S3: Upload HLS/DASH files
    S3-->>Transcoder: Upload success
    Transcoder->>SQS: Send transcode-completed message
    SQS->>Consumer: Consume transcode-completed
    Consumer->>AssetMgr: HandleTranscodeCompletion
    AssetMgr->>Neo4j: Update HLS/DASH video variant
    Neo4j-->>AssetMgr: Variant updated
    AssetMgr-->>Consumer: Completion handled

    Note over User,Consumer: Status Update Flow
    User->>CMS: Refresh asset list
    CMS->>AssetMgr: GraphQL query asset
    AssetMgr->>Neo4j: Get asset with video status
    Neo4j-->>AssetMgr: Return asset data
    AssetMgr-->>CMS: Return asset with video status
    CMS-->>User: Show updated video status
```

## Key Components and Their Roles

### Frontend Components
- **CMS UI**: React Native application for managing assets

### Backend Services
- **Asset Manager**: GraphQL API for asset management and metadata
- **Transcoder Worker**: Background service for video processing
- **Neo4j**: Graph database for asset relationships and metadata

### AWS Services (LocalStack)
- **S3 Storage**: File storage with different buckets for raw, HLS, and DASH content
- **SQS**: Message queue for job coordination
- **Lambda Functions**: Serverless functions for presigned URLs and job triggering
- **API Gateway**: HTTP endpoints for Lambda functions

### Video Processing
- **FFmpeg**: Video analysis and transcoding engine
- **Video Variants**: Raw, HLS, and DASH formats with different storage locations

## Message Flow Details

### 1. Upload Flow
1. User selects video file in CMS
2. CMS requests presigned URL from Lambda via API Gateway
3. Lambda generates S3 presigned URL for direct upload
4. CMS uploads video directly to S3 using presigned URL
5. CMS calls Asset Manager to save video metadata
6. Asset Manager automatically triggers analysis job

### 2. Analysis Flow
1. Transcoder worker consumes analyze job from SQS
2. Downloads video from S3
3. Runs FFmpeg analysis to validate video
4. Sends completion message back to SQS
5. Asset Manager updates video status to "ready"

### 3. Transcoding Flow
1. User manually triggers HLS/DASH transcoding
2. Transcode Lambda sends job to SQS
3. Transcoder worker processes the job
4. Downloads raw video, transcodes with FFmpeg
5. Uploads transcoded files to appropriate S3 bucket
6. Sends completion message with new file locations
7. Asset Manager updates video variants with transcoded content

## Storage Structure

```
S3 Buckets:
├── raw-storage/
│   └── {assetId}/
│       └── {videoType}/
│           └── {timestamp}_{filename}
├── hls-storage/
│   └── {assetId}/
│       └── {videoType}/
│           └── playlist.m3u8 + segments
└── dash-storage/
    └── {assetId}/
        └── {videoType}/
            └── manifest.mpd + segments
```

## Status Transitions

### Video Status Flow
```
pending → analyzing → ready (for raw videos)
pending → transcoding → ready (for HLS/DASH variants)
```



