# Hobby Streamer

Scalable video upload, processing, and streaming platform using AWS Free Tier. Users upload videos processed with AWS Elemental MediaConvert, stored in S3, and streamed via CloudFront. Includes optional FFmpeg for custom transcoding.

## Key Components

S3: Store raw/processed videos.

MediaConvert: Transcode videos.

DynamoDB: Manage video metadata.

CloudFront: Stream videos globally.

Lambda: Automate workflows.

IAM: Secure access.

## Workflow Overview

Video Upload: Videos are uploaded to S3, triggering Lambda.

Processing: Lambda calls MediaConvert, outputs saved to S3.

Metadata: Stored in DynamoDB.

Streaming: CloudFront delivers videos with adaptive bitrate.

Cost Optimization

Use AWS Free Tier (S3, CloudFront, DynamoDB).

Automate storage cleanup with lifecycle rules.

Leverage MediaConvert "Basic" tier for low-cost transcoding.

Optimize caching to minimize S3 requests.
