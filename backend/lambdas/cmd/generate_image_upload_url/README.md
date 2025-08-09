# Generate Image Upload URL Lambda

Lambda for presigned S3 URLs to upload images to `images-storage`.

## Features
Presigned S3 upload URLs, supports common image types, LocalStack support, easy local testing.

## Env Vars
BUCKET_NAME (default: images-storage), BUCKET_REGION (default: us-east-1), AWS_ENDPOINT (default: http://localhost:4566)

## Quick Setup
```bash
./build.sh
```

## Notes
Use with LocalStack in local environments. See `local/build.sh`.

## Storage Pattern
`images-storage/<assetId>/<imageType>/<filename>`
Examples: `images-storage/123/poster/poster.jpg`, `images-storage/123/thumbnail/thumb.jpg` 