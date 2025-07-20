# Domain Events

The project uses domain events to decouple services and enable asynchronous processing. Events flow through SQS queues to keep things loose and scalable.

## Overview

When domain state changes (asset created, video added, etc.), events get published to SQS queues. Other services consume these events to react accordingly — like triggering video analysis or updating caches.

## Event Types

### Asset Events
- `asset.created` - New asset created
- `asset.updated` - Asset metadata changed  
- `asset.deleted` - Asset removed
- `asset.published` - Asset made public

### Video Events
- `video.added` - Video file uploaded
- `video.removed` - Video file deleted
- `video.status.updated` - Processing status changed

### Image Events
- `image.added` - Image uploaded
- `image.removed` - Image deleted

### Bucket Events
- `bucket.created` - New bucket created
- `bucket.updated` - Bucket metadata changed
- `bucket.deleted` - Bucket removed
- `bucket.asset.added` - Asset added to bucket
- `bucket.asset.removed` - Asset removed from bucket

## Job Events

### Job Triggers
- `job` - Triggers video analysis or transcoding

### Job Completions
- `job-completed` - Analysis or transcoding finished

## Event Flow

1. **Domain Change** - Asset manager updates domain state
2. **Event Published** - Domain event sent to SQS queue
3. **Service Consumes** - Other services pick up events
4. **Reaction** - Services react (trigger jobs, update caches, etc.)

## Implementation

Events use a simple JSON envelope:

```json
{
  "type": "event-name",
  "payload": {
    "assetId": "123",
    "videoId": "456",
    "timestamp": "2024-01-01T00:00:00Z",
    "description": "Description"
  }
}
```

## Queues

- `asset-events` - Domain events from asset-manager
- `job-queue` - Job triggers for transcoder
- `completion-queue` - Job completion notifications

