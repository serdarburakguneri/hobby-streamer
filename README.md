# Hobby streamer

Build a scalable, cost-effective video upload, processing, and streaming platform leveraging AWS Free Tier services. Users will upload videos, which will be processed using AWS Elemental MediaConvert, stored securely, and made available for streaming through AWS CloudFront. Additionally, custom video transcoding using tools like FFmpeg will be implemented as an enhancement.

## Key Components

### 1. Infrastructure Overview

AWS S3: For storing raw and processed video files.

AWS Elemental MediaConvert: For video processing and format conversion.

AWS DynamoDB: To manage metadata for videos (e.g., file name, format, status).

AWS CloudFront: For distributing processed videos globally.

AWS Lambda: For triggering workflows (e.g., on video upload).

AWS IAM: For secure role-based access control.

### 2. System Workflow

Video Upload:

Users upload videos to an S3 bucket.

S3 triggers a Lambda function to initiate processing.

Video Processing:

Lambda invokes MediaConvert to transcode the uploaded video.

Transcoded videos are saved back to S3.

Metadata Storage:

Video details (e.g., file path, resolution, format) are saved to DynamoDB.

Video Streaming:

CloudFront serves the processed video files from S3 to end-users.

Adaptive bitrate streaming is configured for optimal performance.

Cost Optimization

AWS Free Tier Services:

Utilize the free tier limits for S3, CloudFront, and DynamoDB.

Use small-size EC2 instances if needed for additional processing (e.g., FFmpeg testing).

Storage Management:

Automatically delete raw video files after processing.

Apply lifecycle rules to S3 buckets to archive older files to Glacier.

MediaConvert:

Use the MediaConvert "Basic" tier for cost-effective transcoding.

CloudFront:

Optimize caching and minimize requests to S3 to reduce costs.
