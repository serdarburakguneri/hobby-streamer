# Hobby Streamer

This hobby project leverages the AWS Free Tier to build a lightweight CMS for a streaming platform. It allows users to upload videos, manage content, and prepare files for streaming. The goal is to create an end-to-end workflow that covers video ingestion, processing, and delivery with minimal infrastructure cost.

## Tech Stack
	AWS S3 – Video storage
	AWS Lambda (Go) – Serverless backend logic
	Amazon API Gateway – API endpoint management
	AWS Elastic Transcoder / MediaConvert – Video processing and transcoding
	DynamoDB – Metadata and CMS data storage
	CloudFront – Content delivery (CDN)
	Terraform – Infrastructure as code