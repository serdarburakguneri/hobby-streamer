# Generate Video Upload URL Lambda

Lambda for presigned S3 URLs to upload videos to `raw-storage`.

## Features
Presigned S3 upload URLs, LocalStack support, easy local testing, simple environment config.

## Env Vars
BUCKET_NAME (default: raw-storage), BUCKET_REGION (default: us-east-1), AWS_ENDPOINT (default: http://localhost:4566)

## Quick Setup
```bash
./build.sh
```

## Notes
Use with LocalStack in local environments. See `local/build.sh`.