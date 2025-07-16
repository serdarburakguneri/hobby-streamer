# Messages Package

Common SQS message payload structures used across services to ensure consistency.

## Message Types

### Job Messages (sent to transcoder)

#### Analyze Job
- **Type**: `analyze`
- **Payload**: `AnalyzePayload`
- **Purpose**: Triggers video analysis

#### Transcode HLS Job
- **Type**: `transcode-hls`
- **Payload**: `TranscodePayload`
- **Purpose**: Triggers HLS transcoding

#### Transcode DASH Job
- **Type**: `transcode-dash`
- **Payload**: `TranscodePayload`
- **Purpose**: Triggers DASH transcoding

### Completion Messages (sent from transcoder)

#### Analyze Completion
- **Type**: `analyze-completed`
- **Payload**: `AnalyzeCompletionPayload`
- **Purpose**: Notifies analyze job completion

#### Transcode HLS Completion
- **Type**: `transcode-hls-completed`
- **Payload**: `TranscodeCompletionPayload`
- **Purpose**: Notifies HLS transcoding completion

#### Transcode DASH Completion
- **Type**: `transcode-dash-completed`
- **Payload**: `TranscodeCompletionPayload`
- **Purpose**: Notifies DASH transcoding completion

## Usage

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"

// Create analyze payload
analyzePayload := messages.AnalyzePayload{
    Input:     "s3://bucket/key.mp4",
    AssetID:   "asset-123",
    VideoType: "main",
}

// Send message using SQS producer
err := producer.SendMessage(ctx, messages.MessageTypeAnalyze, analyzePayload)

// Handle completion message
func handleAnalyzeCompletion(payload messages.AnalyzeCompletionPayload) {
    if payload.Success {
        // Handle success
    } else {
        // Handle error
        log.Error("Analysis failed", "error", payload.Error)
    }
}
```

## Benefits

- **Type Safety**: Compile-time checking of message structures
- **Consistency**: Shared structures prevent message format mismatches
- **Documentation**: Clear definition of all message types and payloads
- **Maintainability**: Single source of truth for message formats 