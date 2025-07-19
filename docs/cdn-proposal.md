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

### Content Types
- **source/** - Raw video files (MP4, MOV, etc.)
- **hls/** - HLS streaming files (M3U8 + TS segments)
- **dash/** - DASH streaming files (MPD + M4S segments)
- **images/** - Images (posters, thumbnails, banners)

## Cross-Region Replication

### Setup
Cross-region replication is automatically configured when running `./setup-s3-buckets.sh`:

```bash
# Enable versioning on both buckets
aws s3api put-bucket-versioning --bucket content-east --versioning-configuration Status=Enabled

# Configure replication from primary to secondary
aws s3api put-bucket-replication --bucket content-east --replication-configuration file://replication-config.json
```

### Replication Rules
- **Source**: `content-east` (us-east-1)
- **Destination**: `content-west` (us-west-2)
- **Storage Class**: STANDARD
- **Delete Marker Replication**: Enabled
- **Status**: Enabled

## CDN with Automatic Failover

### Nginx Configuration
The nginx server (port 8083) acts as a CDN with automatic failover:

```nginx
# Primary upstream
upstream s3_primary {
    server localstack:4566;
}

# Secondary upstream (failover)
upstream s3_secondary {
    server localstack:4566;
}

# CDN endpoint with failover
location /cdn/ {
    proxy_pass http://s3_primary/content-east/;
    proxy_intercept_errors on;
    error_page 502 503 504 500 = @s3_failover;
}

# Failover location
location @s3_failover {
    proxy_pass http://s3_secondary/content-west/;
    add_header X-CDN-Failover "true" always;
}
```

### URL Structure
- **Primary**: `http://localhost:8083/cdn/{assetId}/{type}/{quality}/{filename}`
- **Automatic Failover**: If primary fails, nginx automatically routes to secondary
- **Headers**: `X-CDN-Failover: true` indicates when failover occurred

## StreamInfo Integration

### Video StreamInfo
```json
{
  "streamInfo": {
    "cdnPrefix": "http://localhost:8083/cdn",
    "url": "http://localhost:8083/cdn/asset-123/hls/1080p/playlist.m3u8"
  }
}
```

### Image StreamInfo
```json
{
  "streamInfo": {
    "cdnPrefix": "http://localhost:8083/cdn",
    "url": "http://localhost:8083/cdn/asset-123/images/poster/poster.jpg"
  }
}
```

## Configuration

### Asset Manager Config
```yaml
cdn:
  prefix: "http://localhost:8083/cdn"
```

### Lambda Functions
- **Video Upload**: Uses `content-east` bucket with `{assetId}/source/{filename}` path
- **Image Upload**: Uses `content-east` bucket with `{assetId}/images/{type}/{filename}` path
- **Transcoder**: Outputs to `content-east` bucket with `{assetId}/{videoId}/{format}/{filename}` path

## Frontend Integration

### Video Player
```typescript
const getVideoUrl = (video: Video) => {
  if (video.streamInfo?.url) {
    return video.streamInfo.url; // CDN URL with failover
  }
  return video.storageLocation.url; // Fallback to direct S3 URL
};
```

### Image Display
```typescript
const getImageUrl = (image: Image) => {
  if (image.streamInfo?.url) {
    return image.streamInfo.url; // CDN URL with failover
  }
  return image.url; // Fallback to direct S3 URL
};
```

## Setup Commands

### Initial Setup
```bash
# Run the complete setup
./local/build.sh

# Or run S3 setup separately
./local/setup-s3-buckets.sh
```

### Manual Setup
```bash
# Create buckets
aws --endpoint-url=http://localhost:4566 s3api create-bucket --bucket content-east --region us-east-1
aws --endpoint-url=http://localhost:4566 s3api create-bucket --bucket content-west --region us-west-2

# Enable versioning
aws --endpoint-url=http://localhost:4566 s3api put-bucket-versioning --bucket content-east --versioning-configuration Status=Enabled

# Configure replication
aws --endpoint-url=http://localhost:4566 s3api put-bucket-replication --bucket content-east --replication-configuration file://replication-config.json
```

## Testing Failover

### Test Primary Region
```bash
# Upload to primary
aws --endpoint-url=http://localhost:4566 s3 cp test.mp4 s3://content-east/asset-123/source/test.mp4

# Access via CDN
curl http://localhost:8083/cdn/asset-123/source/test.mp4
```

### Test Failover
```bash
# Simulate primary failure (stop localstack temporarily)
docker-compose stop localstack

# Access via CDN (should failover to secondary)
curl http://localhost:8083/cdn/asset-123/source/test.mp4
# Response should include: X-CDN-Failover: true
```

## Monitoring

### Health Check
```bash
# Check CDN health
curl http://localhost:8083/health
# Response: "healthy"
```

### Replication Status
```bash
# Check replication status
aws --endpoint-url=http://localhost:4566 s3api head-object --bucket content-east --key asset-123/source/test.mp4
# Look for ReplicationStatus in response
``` 