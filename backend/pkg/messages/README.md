# Messages Package

A centralized message definitions package for the Hobby Streamer backend services.

## Overview

This package defines the message structures and types used for inter-service communication via SQS queues.

## Message Types

### Job Messages

| Message Type | Payload Type | Description |
|--------------|--------------|-------------|
| `job` | `JobPayload` | Job trigger (analyze, transcode) |
| `job-completed` | `JobCompletionPayload` | Job completion notification |

## Payload Structures

### JobPayload

Used for triggering jobs (analyze, transcode):

```go
type JobPayload struct {
    JobType      string `json:"jobType"`      // "analyze" or "transcode"
    Input        string `json:"input"`        // Input file path/URL
    AssetID      string `json:"assetId"`      // Asset identifier
    VideoID      string `json:"videoId"`      // Video identifier
    Format       string `json:"format,omitempty"`       // Format for transcode jobs
    Quality      string `json:"quality,omitempty"`      // Quality setting
    OutputBucket string `json:"outputBucket,omitempty"` // S3 bucket for output
    OutputKey    string `json:"outputKey,omitempty"`    // S3 key for output
}
```

### JobCompletionPayload

Used for job completion notifications:

```go
type JobCompletionPayload struct {
    JobType      string  `json:"jobType"`      // "analyze" or "transcode"
    AssetID      string  `json:"assetId"`      // Asset identifier
    VideoID      string  `json:"videoId"`      // Video identifier
    Format       string  `json:"format,omitempty"`       // Format for transcode jobs
    Success      bool    `json:"success"`      // Job success status
    Error        string  `json:"error,omitempty"`        // Error message if failed
    
    // Video metadata (for analyze jobs)
    Width        int     `json:"width,omitempty"`
    Height       int     `json:"height,omitempty"`
    Duration     float64 `json:"duration,omitempty"`
    Bitrate      int     `json:"bitrate,omitempty"`
    Codec        string  `json:"codec,omitempty"`
    Size         int64   `json:"size,omitempty"`
    ContentType  string  `json:"contentType,omitempty"`
    
    // Transcode metadata (for transcode jobs)
    Bucket       string  `json:"bucket,omitempty"`
    Key          string  `json:"key,omitempty"`
    URL          string  `json:"url,omitempty"`
}
```

## Usage

### Publishing Messages

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"

// Create a job payload
jobPayload := messages.JobPayload{
    JobType: "analyze",
    Input:   "s3://bucket/video.mp4",
    AssetID: "asset-123",
    VideoID: "video-456",
}

// Send the message
err := producer.SendMessage(ctx, messages.MessageTypeJob, jobPayload)
```

### Consuming Messages

```go
// Handle job messages
if msgType == messages.MessageTypeJob {
    var payload messages.JobPayload
    if err := json.Unmarshal(payloadBytes, &payload); err != nil {
        return err
    }
    
    switch payload.JobType {
    case "analyze":
        return handleAnalyzeJob(ctx, payload)
    case "transcode":
        return handleTranscodeJob(ctx, payload)
    }
}

// Handle completion messages
if msgType == messages.MessageTypeJobCompleted {
    var payload messages.JobCompletionPayload
    if err := json.Unmarshal(payloadBytes, &payload); err != nil {
        return err
    }
    
    switch payload.JobType {
    case "analyze":
        return handleAnalyzeCompletion(ctx, payload)
    case "transcode":
        return handleTranscodeCompletion(ctx, payload)
    }
}
```

## Queue Architecture

The system uses a simplified two-queue architecture:

1. **Job Queue** (`job-queue`): 
   - Asset-manager publishes job triggers
   - Transcoder consumes and processes jobs

2. **Completion Queue** (`completion-queue`):
   - Transcoder publishes job completions
   - Asset-manager consumes and updates asset status

## Message Flow

### Analyze Job Flow

1. Asset-manager publishes `JobPayload` with `jobType: "analyze"` to job-queue
2. Transcoder consumes and processes the analyze job
3. Transcoder publishes `JobCompletionPayload` with `jobType: "analyze"` to completion-queue
4. Asset-manager consumes and updates video metadata

### Transcode Job Flow

1. Asset-manager publishes `JobPayload` with `jobType: "transcode"` to job-queue
2. Transcoder consumes and processes the transcode job
3. Transcoder publishes `JobCompletionPayload` with `jobType: "transcode"` to completion-queue
4. Asset-manager consumes and updates transcoding status

## Best Practices

1. Always validate message types before processing
2. Use structured error handling for failed jobs
3. Include all relevant metadata in completion payloads
4. Use consistent naming conventions for job types
5. Handle unknown job types gracefully