# CDN Proposal: Single Bucket Storage with Cross-Region Replication

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   Primary       │    │   Secondary     │
│   Zone          │    │   Zone          │
│   (AZ1)         │    │   (AZ2)         │
├─────────────────┤    ├─────────────────┤
│ content-east    │    │ content-west    │
│                 │    │                 │
│ {assetId}/      │    │ {assetId}/      │
│ ├── source/     │    │ ├── source/     │
│ ├── hls/        │    │ ├── hls/        │
│ ├── dash/       │    │ ├── dash/       │
│ └── images/     │    │ └── images/     │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────┬───────────────┘
                 │
         ┌─────────────────┐
         │   Nginx CDN     │
         │   (Port 8083)   │
         │                 │
         │ Primary → Failover
         │ Automatic routing
         └─────────────────┘
                 │
         ┌─────────────────┐
         │   Content       │
         │   Viewers       │
         └─────────────────┘
```

## Storage Structure

### Bucket Organization
```
content-east/
├── asset-123/
│   ├── source/
│   │   └── original.mp4
│   ├── hls/
│   │   ├── 1080p/
│   │   │   ├── playlist.m3u8
│   │   │   └── segment_001.ts
│   │   └── 720p/
│   │       ├── playlist.m3u8
│   │       └── segment_001.ts
│   ├── dash/
│   │   ├── 1080p/
│   │   │   ├── manifest.mpd
│   │   │   └── segment_001.m4s
│   │   └── 720p/
│   │       ├── manifest.mpd
│   │       └── segment_001.m4s
│   └── images/
│       ├── poster/
│       │   └── poster.jpg
│       ├── thumbnail/
│       │   └── thumbnail.jpg
│       └── banner/
│           └── banner.jpg
└── asset-456/
    └── ...
```
Types: source (raw video), hls (HLS files), dash (DASH files), images (posters, thumbnails, banners).

## Cross-Region Replication
Setup: `./setup-s3-buckets.sh` enables versioning and replication from content-east (us-east-1) to content-west (us-west-2). Storage class: STANDARD, delete marker replication: enabled.

## CDN with Automatic Failover
Nginx (port 8083) acts as CDN with failover. Primary: http://localhost:8083/cdn/{assetId}/{type}/{quality}/{filename}. If primary fails, routes to secondary. Header: X-CDN-Failover: true on failover.

## StreamInfo Integration
Video: streamInfo.cdnPrefix, streamInfo.url. Image: streamInfo.cdnPrefix, streamInfo.url.

## Configuration
Asset manager: cdn.prefix. Lambdas: video upload uses content-east/{assetId}/source/{filename}, image upload uses content-east/{assetId}/images/{type}/{filename}, transcoder outputs to content-east/{assetId}/{videoId}/{format}/{filename}.

## Frontend Integration
Video player and image display use streamInfo.url if available, fallback to direct S3 URL.

## Setup Commands
- `./local/build.sh` (all setup)
- `./local/setup-s3-buckets.sh` (S3 only)
Manual: create buckets, enable versioning, configure replication.

## Testing Failover
Upload to primary, access via CDN, simulate failure, check X-CDN-Failover header.

## Monitoring
- CDN health: curl http://localhost:8083/health
- Replication status: aws s3api head-object 