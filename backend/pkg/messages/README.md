# Messages Package

Central definitions for messages passed between services over SQS. Keeps job and completion message types consistent across services — especially for things like transcoding, analysis, and async workflows.

---

## Message Types

###  Job Messages (→ Transcoder)

| Type             | Payload Type     | Purpose                 |
|------------------|------------------|--------------------------|
| `analyze`        | `AnalyzePayload` | Start video analysis     |
| `transcode-hls`  | `TranscodePayload` | Transcode to HLS format |
| `transcode-dash` | `TranscodePayload` | Transcode to DASH format |

---

###  Completion Messages (← Transcoder)

| Type                      | Payload Type                | Purpose                        |
|---------------------------|-----------------------------|--------------------------------|
| `analyze-completed`       | `AnalyzeCompletionPayload`  | Video analysis done            |
| `transcode-hls-completed` | `TranscodeCompletionPayload`| HLS transcoding finished       |
| `transcode-dash-completed`| `TranscodeCompletionPayload`| DASH transcoding finished      |

---

## Usage Example

### Sending a Job Message

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"

payload := messages.AnalyzePayload{
    Input:     "s3://bucket/video.mp4",
    AssetID:   "asset-123",
    VideoType: "main",
}

err := producer.SendMessage(ctx, messages.MessageTypeAnalyze, payload)
```

---

### Handling a Completion Message

```go
func handleAnalyzeCompletion(payload messages.AnalyzeCompletionPayload) {
    if payload.Success {
        // Do something with analysis results
    } else {
        log.Error("Analysis failed", "error", payload.Error)
    }
}
```

---

> 📨 Keeping message formats centralized helps reduce bugs when services talk to each other — and avoids mismatches during refactors.