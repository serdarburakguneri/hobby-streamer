# Hobby Streamer

This hobby project leverages the AWS Free Tier to build a lightweight CMS and a streaming platform. It allows users to upload videos, manage content, and prepare files for streaming. The goal is to create an end-to-end workflow that covers video ingestion, processing, and delivery with minimal infrastructure cost.

## Tech Stack
	AWS S3 – Video storage
	AWS Lambda (Go) – Serverless backend logic
    AWS SQS - For Async communication between internal services
	Amazon API Gateway – API endpoint management
	AWS Fargate - A transcoder service powered by FFMPEG – Video processing and transcoding
	DynamoDB – Metadata and CMS data storage
	CloudFront – Content delivery (CDN)
	Terraform – Infrastructure as code

## 🗂️ Architecture

![Architecture Diagram](docs/hobby-streamer.drawio.svg)

## Status

Currently, I've been designing the architecture and meanwhile terraforming the backend parts.

## TODO

- Centralized logging
- A search mechanism for the asset manager service

## How to run
### Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed
- [AWS CLI](https://aws.amazon.com/cli/) installed and configured with your AWS credentials
- [Go](https://go.dev/doc/install) installed

### Steps

- Clone the repository:
- in the root directory, run build.sh to compile the Go services.
- in terraform directory, have your variables set up in `terraform.tfvars` file
- Run `terraform init` to initialize the Terraform configuration.
- Run `terraform plan` to see the changes that will be applied.
- Run `terraform apply` to create the infrastructure.


