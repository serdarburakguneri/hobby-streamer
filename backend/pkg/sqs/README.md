# SQS Package

A shared library for working with AWS SQS. Includes producer and consumer utilities, along with a registry system for managing multiple consumers in one place.

## Features

- **Producer** – Send structured messages to SQS queues
- **Consumer** – Process messages with handler functions
- **Consumer Registry** – Register and run multiple consumers centrally
- **Message Format** – Consistent message schema with `type` and `payload` fields
- **Context-Aware** – All operations support cancellation and shutdown via `context.Context`

---

## Components

### Producer

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

producer, err := sqs.NewProducer(ctx, "queue-url")
if err != nil {
    return err
}

err = producer.SendMessage(ctx, "message-type", payload)
```

---

### Consumer

```go
consumer, err := sqs.NewConsumer(ctx, "queue-url")
if err != nil {
    return err
}

consumer.Start(ctx, func(msg sqs.Message) error {
    // Handle the message
    return nil
})
```

---

### Consumer Registry

```go
consumerRegistry := sqs.NewConsumerRegistry()

consumerRegistry.Register("queue-url-1", handler1)
consumerRegistry.Register("queue-url-2", handler2)

consumerRegistry.Start(context.Background())

// On shutdown
consumerRegistry.Stop()
```

---

## Custom Message Handlers

```go
type MyHandler struct {
    logger *logger.Logger
}

func (h *MyHandler) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
    h.logger.Info("Handling message", "type", msgType)

    switch msgType {
    case "my-event":
        return h.handleMyEvent(ctx, payload)
    default:
        return nil
    }
}

func (h *MyHandler) handleMyEvent(ctx context.Context, payload map[string]interface{}) error {
    // Do something with the payload
    return nil
}

// Register the handler
myHandler := &MyHandler{logger: logger.Get().WithService("my-handler")}
consumerRegistry.Register("my-queue-url", myHandler.HandleMessage)
```

---

## Message Format

All messages follow the same structure:

```json
{
  "type": "message-type",
  "payload": {
    "key": "value"
  }
}
```

---

## Error Handling

- Consumer continues running even if message processing fails
- Failed messages are logged (SQS handles retries)

---

