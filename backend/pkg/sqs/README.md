# SQS Package

This package provides SQS producer functionality for sending messages to SQS queues in the hobby-streamer project.

## Features

- SQS message producer for sending jobs to queues
- LocalStack support for local development
- Structured logging integration
- JSON message serialization

## Usage

### Basic Setup

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"

func main() {
    ctx := context.Background()
    queueURL := "http://localhost:4566/000000000000/transcoder-jobs"
    
    producer, err := sqs.NewProducer(ctx, queueURL)
    if err != nil {
        log.Fatal(err)
    }
    
    // Send a message
    payload := map[string]interface{}{
        "input": "s3://raw-storage/video.mp4",
    }
    
    err = producer.SendMessage(ctx, "analyze", payload)
    if err != nil {
        log.Error(err)
    }
}
```

### Message Format

Messages are automatically wrapped in the expected format:

```json
{
  "type": "analyze",
  "payload": {
    "input": "s3://raw-storage/video.mp4"
  }
}
```

## Environment Variables

- `AWS_ENDPOINT`: Custom endpoint for SQS (default: `http://localhost:4566` for LocalStack)
- `AWS_REGION`: AWS region (default: `us-east-1`)

## Integration

This package is designed to work with:
- **Asset Manager Service**: Sends analyze jobs when videos are uploaded
- **Transcoder Service**: Consumes jobs from the same queue
- **LocalStack**: For local development and testing 