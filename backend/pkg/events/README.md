# Events Package

CloudEvents 1.0 compliant event streaming library for Kafka integration, providing producer/consumer patterns with correlation tracking and schema evolution.

## Features
CloudEvents 1.0 format, Kafka producer/consumer, event correlation, causation tracking, schema versioning, error handling, dead letter topics, monitoring integration.

## CloudEvents Format
```json
{
  "specversion": "1.0",
  "id": "uuid-v4",
  "source": "asset-manager",
  "type": "com.hobbystreamer.asset.created",
  "datacontenttype": "application/json",
  "time": "2024-01-01T12:00:00Z",
  "data": {
    "assetId": "asset-123",
    "slug": "my-video",
    "title": "My Video"
  },
  "correlationid": "req-456",
  "causationid": "event-789"
}
```

## Usage

### Producer
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"

producer, err := events.NewProducer(ctx, "localhost:9092")
defer producer.Close()

event := events.NewEvent("asset.created", map[string]interface{}{
    "assetId": "asset-123",
    "slug": "my-video",
})

err = producer.SendEvent(ctx, "asset-events", event)
```

### Consumer
```go
consumer, err := events.NewConsumer(ctx, "localhost:9092", "asset-manager-group")
defer consumer.Close()

consumer.Subscribe("asset-events", func(event *events.Event) error {
    // Handle event
    return nil
})

consumer.Start(ctx)
```

## Event Types
- Asset: `asset.created`, `asset.updated`, `asset.deleted`, `asset.published`
- Video: `video.added`, `video.removed`, `video.status.updated`
- Bucket: `bucket.created`, `bucket.updated`, `bucket.deleted`
- Job: `job.analyze.requested`, `job.transcode.requested`, `job.completed`
- Content: `content.analysis.requested`, `content.analysis.completed`

## Best Practices
Use correlation IDs for tracing, implement proper error handling, monitor consumer lag, version event schemas, use dead letter topics for failed events. 