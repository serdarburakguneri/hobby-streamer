# DASH Job Requested Lambda

Receives DASH transcode requests and publishes job events to Kafka for the transcoder.

## Features
HTTP endpoint, CORS, Kafka publishing, validation, structured logging.

## Setup
```bash
./local/build.sh
```

Env: `KAFKA_BOOTSTRAP_SERVERS`.