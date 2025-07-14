# SQS Package

This package provides AWS SQS client functionality including producers, consumers, and a consumer registry system for managing multiple message consumers.

## Overview

The SQS package provides:
- **Producer**: Send messages to SQS queues
- **Consumer**: Receive and process messages from SQS queues
- **Consumer Registry**: Manage multiple consumers with a unified interface

## Components

### Producer
The `Producer` sends messages to SQS queues with structured payloads.

### Consumer
The `Consumer` receives messages from SQS queues and processes them.

### Consumer Registry
The `ConsumerRegistry` manages multiple consumers and provides a unified interface for starting and stopping them.

### MessageHandler
A function type that processes messages: `func(ctx context.Context, msgType string, payload map[string]interface{}) error`

## Usage

### Producer

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

producer, err := sqs.NewProducer(ctx, "queue-url")
if err != nil {
    return err
}

err = producer.SendMessage(ctx, "message-type", payload)
```

### Consumer

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

consumer, err := sqs.NewConsumer(ctx, "queue-url")
if err != nil {
    return err
}

consumer.Start(ctx, func(msg sqs.Message) error {
    // Process message
    return nil
})
```

### Consumer Registry

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

// Create a registry
consumerRegistry := sqs.NewConsumerRegistry()

// Register consumers
consumerRegistry.Register("queue-url-1", handler1)
consumerRegistry.Register("queue-url-2", handler2)

// Start all consumers
consumerRegistry.Start(context.Background())

// Stop all consumers (on shutdown)
consumerRegistry.Stop()
```

### Creating Custom Message Handlers

```go
type MyHandler struct {
    logger *logger.Logger
}

func (h *MyHandler) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
    h.logger.Info("Processing message", "message_type", msgType)
    
    switch msgType {
    case "my-event":
        return h.handleMyEvent(ctx, payload)
    default:
        return nil
    }
}

func (h *MyHandler) handleMyEvent(ctx context.Context, payload map[string]interface{}) error {
    // Handle your specific event
    return nil
}

// Register in your service
myHandler := &MyHandler{logger: logger.Get().WithService("my-handler")}
consumerRegistry.Register("my-queue-url", myHandler.HandleMessage)
```

## Message Format

Messages are sent and received in a structured format:

```json
{
  "type": "message-type",
  "payload": {
    "key": "value"
  }
}
```

## Error Handling

- Consumers continue running even if individual message processing fails
- Failed messages are logged but not retried (SQS handles retries)
- If a consumer fails to initialize, it's logged and skipped
- Empty queue URLs are skipped with a warning

## Dependencies

This package depends on:
- `github.com/aws/aws-sdk-go-v2` - AWS SDK for Go v2
- `github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger` - For structured logging 