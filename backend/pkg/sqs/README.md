# SQS Package

Go library for AWS SQS: producer/consumer logic, registry for multiple consumers, structured messages.

## Features
Producer (send typed messages), consumer (process with handler), registry (run multiple consumers), consistent message format, context support.

## Usage
### Producer
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
producer, err := sqs.NewProducer(ctx, "queue-url")
err = producer.SendMessage(ctx, "transcode-hls", payload)
```

### Consumer
```go
consumer, err := sqs.NewConsumer(ctx, "queue-url")
consumer.Start(ctx, func(msg sqs.Message) error { return nil })
```

### Registry
```go
registry := sqs.NewConsumerRegistry()
registry.Register("queue-url-1", handler1)
registry.Register("queue-url-2", handler2)
registry.Start(context.Background())
registry.Stop()
```

## Message Format
```json
{ "type": "event-name", "payload": { "key": "value" } }
```

## Error Handling
Failed messages are logged, consumers keep running, custom logging/DLQ can be added.

