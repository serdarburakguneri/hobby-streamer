# Events

CloudEvents 1.0 producer/consumer helpers for Kafka with correlation and simple patterns.

## Features
CloudEvents 1.0, Kafka producer/consumer, correlation/causation IDs, error handling, DLQ-friendly.

## Quick usage
```go
producer, _ := events.NewProducer(ctx, "localhost:9092")
defer producer.Close()

evt := events.NewEvent("asset.created", map[string]any{"assetId": "a1"})
_ = producer.SendEvent(ctx, "asset-events", evt)
```

See `backend/pkg/events/example/` for a fuller example.