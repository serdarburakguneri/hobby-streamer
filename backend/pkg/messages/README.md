# Messages Package

Centralized message definitions for backend services, used for SQS inter-service communication.

## Features
Typed message structures, job/completion payloads, consistent envelope, easy integration with SQS, simple handler patterns.

## Message Types
- `job`: Job trigger (analyze, transcode)
- `job-completed`: Job completion notification

## Usage Example
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
jobPayload := messages.JobPayload{JobType: "analyze", Input: "s3://bucket/video.mp4", AssetID: "asset-123", VideoID: "video-456"}
err := producer.SendMessage(ctx, messages.MessageTypeJob, jobPayload)
```

## Queue Architecture
Job queue: asset-manager → transcoder, Completion queue: transcoder → asset-manager.

## Best Practices
Validate message types, use structured error handling, include relevant metadata, use consistent naming, handle unknown types gracefully.