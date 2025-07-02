# Hobby Streamer

Hobby Streamer is a lightweight content management system (CMS) and streaming platform designed for experimenting end-to-end video workflows. The project enables you to:

- **Upload and manage video assets** through a simple API.
- **Process and transcode videos** for adaptive streaming (HLS/DASH) using FFMPEG.
- **Organize content** into buckets for easy management.
- **Deliver video content** in a way that mimics modern streaming platforms.

The goal is to provide a hands-on, cost-free environment for learning and prototyping solutions—without relying on any paid cloud infrastructure. All services (asset manager, transcoder, storage) run on your local machine, and AWS services (DynamoDB, SQS) are emulated using [LocalStack](https://github.com/localstack/localstack).

## Tech Stack
- LocalStack (DynamoDB, SQS, S3. Lambda) – Local AWS service emulation
- Go – Backend code for all services
- FFMPEG – For the transcoder service 

## 🗂️ Architecture

![Architecture Diagram](docs/hobby-streamer.drawio.svg)

## TODO

- Centralized logging
- A search mechanism for the asset manager service

## 🧪 Local Testing

### Prerequisites
- [Docker](https://www.docker.com/products/docker-desktop/) installed
- [Go](https://go.dev/doc/install) installed

### Local Environment Setup

To set up your entire local AWS-like environment (LocalStack, S3 buckets, DynamoDB tables, SQS queues) and start the core services, simply run:

```sh
./build.sh
```

This script will:
- Start LocalStack (via Docker Compose) if not already running
- Wait for LocalStack to be ready
- Create the required S3 buckets: `raw-storage`, `transcoded-storage`, `thumbnails-storage`
- Create the required DynamoDB tables: `asset`, `bucket`
- Create the required SQS queue: `transcoder-jobs`
- Start the Asset Manager service (on port 8080)
- Start the Transcoder service (connected to the local SQS queue)

Logs for these services are written to `asset-manager.log` and `transcoder.log` in the project root.

No manual setup is needed—just run the script and you're ready to go!

### Running the Services

The Asset Manager and Transcoder services are started automatically by `build.sh`.

### Testing the API
`tbd`

## Notes
tbd


