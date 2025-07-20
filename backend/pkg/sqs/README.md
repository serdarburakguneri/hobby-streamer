# SQS Package

A shared library for working with AWS SQS in Go. Wraps producer/consumer logic and adds a simple registry to manage multiple consumers in one place. Built for structured messages and clean handler patterns.

---

## Features

-  **Producer** – Send typed messages to SQS queues  
-  **Consumer** – Receive and process messages via handler functions  
-  **Registry** – Register and run multiple consumers together  
-  **Consistent message format** – Every message has a `type` and a `payload`  
-  **Context-aware** – Supports cancellation and clean shutdown

---

## Core Components

### Producer

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

producer, err := sqs.NewProducer(ctx, "queue-url")
if err != nil {
    return err
}

err = producer.SendMessage(ctx, "transcode-hls", payload)
```

---

### Consumer

```go
consumer, err := sqs.NewConsumer(ctx, "queue-url")
if err != nil {
    return err
}

consumer.Start(ctx, func(msg sqs.Message) error {
    // Handle message
    return nil
})
```

---

### Consumer Registry

Use the registry to run multiple consumers from one place:

```go
registry := sqs.NewConsumerRegistry()

registry.Register("queue-url-1", handler1)
registry.Register("queue-url-2", handler2)

registry.Start(context.Background())

// On shutdown
registry.Stop()
```

---

## Custom Message Handlers

You can define a single handler for multiple message types using `msgType`:

```go
type MyHandler struct {
    logger *logger.Logger
}

func (h *MyHandler) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
    h.logger.Info("Handling", "type", msgType)

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
registry.Register("my-queue-url", myHandler.HandleMessage)
```

---

## Message Format

All messages use a consistent envelope:

```json
{
  "type": "event-name",
  "payload": {
    "key": "value"
  }
}
```

This makes it easier to share messages across services without version drift.

---

## Error Handling

- Failed message processing is logged (SQS handles retries via visibility timeout)
- Consumers keep running even if a message fails
- You can plug in custom logging or DLQ handling later

---

