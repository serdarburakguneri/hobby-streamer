# Domain Events

Domain events decouple services and enable async processing. Events flow through SQS queues for loose, scalable integration.

## Overview
When domain state changes (asset created, video added, etc.), events are published to SQS. Other services consume and react (trigger analysis, update caches, etc).

## Event Types
Asset: asset.created, asset.updated, asset.deleted, asset.published. Video: video.added, video.removed, video.status.updated. Image: image.added, image.removed. Bucket: bucket.created, bucket.updated, bucket.deleted, bucket.asset.added, bucket.asset.removed. Job: job (trigger), job-completed (done).

## Event Flow
1. Domain change (asset-manager updates state)
2. Event published (to SQS)
3. Service consumes (other services pick up)
4. Reaction (trigger jobs, update caches, etc)

## Event Envelope
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
asset-events (from asset-manager), job-queue (for transcoder), completion-queue (job done).

