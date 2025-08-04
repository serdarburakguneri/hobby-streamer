# DASH Job Requested Lambda

This lambda function handles DASH transcoding job requests from the CMS UI. It receives HTTP requests and publishes DASH job events to Kafka for processing by the transcoder service.

## Features

- HTTP API endpoint for triggering DASH transcoding jobs
- CORS support for CMS UI integration
- Kafka event publishing for job processing
- Input validation and error handling
- Structured logging

## Usage

The lambda is exposed via API Gateway at the `/dash-job-requested` endpoint and accepts POST requests with the following JSON payload:

```json
{
  "assetId": "asset-123",
  "videoId": "video-456", 
  "input": "s3://bucket-name/path/to/video.mp4"
}
```

## Response

Success response (200):
```json
{
  "message": "DASH job requested successfully",
  "assetId": "asset-123",
  "videoId": "video-456"
}
```

Error response (400/500):
```json
{
  "error": "Error description"
}
```

## Environment Variables

- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (default: kafka:29092)

## Dependencies

- AWS Lambda Go runtime
- Kafka producer for event publishing
- Structured logging 