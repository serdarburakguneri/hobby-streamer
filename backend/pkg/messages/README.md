# Messages Package

Shared message definitions used for SQS communication between backend services. Provides a central place for defining job and completion message types to ensure consistency across services.

## Message Types

### Job Messages (To Transcoder)

| Type               | Payload Type           | Purpose                  |
|--------------------|------------------------|--------------------------|
| `analyze`          | `AnalyzePayload`       | Trigger video analysis   |
| `transcode-hls`    | `TranscodePayload`     | Trigger HLS transcoding  |
| `transcode-dash`   | `TranscodePayload`     | Trigger DASH transcoding |

### Completion Messages (From Transcoder)

| Type                      | Payload Type               | Purpose                          |
|---------------------------|----------------------------|----------------------------------|
| `analyze-completed`       | `AnalyzeCompletionPayload` | Notify that analysis is complete |
| `transcode-hls-completed` | `TranscodeCompletionPayload` | Notify HLS transcoding complete |
| `transcode-dash-completed`| `TranscodeCompletionPayload` | Notify DASH transcoding complete |

---

## Usage Example

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"

// Create an analyze job payload
payload := messages.AnalyzePayload{
    Input:     "s3://bucket/video.mp4",
    AssetID:   "asset-123",
    VideoType: "main",
}

// Send to SQS using a producer
err := producer.SendMessage(ctx, messages.MessageTypeAnalyze, payload)
```

### Handling Completion

```go
func handleAnalyzeCompletion(payload messages.AnalyzeCompletionPayload) {
    if payload.Success {
        // Process successful analysis
    } else {
        // Handle failure
        log.Error("Analysis failed", "error", payload.Error)
    }
}
```

---
