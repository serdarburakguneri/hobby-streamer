# Asset Manager Service

The Asset Manager service provides RESTful APIs to manage media assets and buckets, including asset creation, metadata management, image/video association, and bucket organization.

## Features
- CRUD operations for assets and buckets
- Associate images and video variants with assets
- Organize assets into buckets
- DynamoDB-backed storage

## Requirements
- Go 1.22+
- LocalStack (for local AWS emulation)
- Docker (optional, for containerized runs)

## Environment Variables
- `PORT`: HTTP port to listen on (default: 8080)
- AWS credentials and region (for DynamoDB access, handled by LocalStack or your AWS profile)

## Running Locally

### 1. Start LocalStack and create DynamoDB tables (see project root `build.sh`)

### 2. Run the service:
```sh
cd services/asset-manager
PORT=8080 go run ./cmd/main.go
```

### 3. Or build and run with Docker:
```sh
docker build -t asset-manager .
docker run -p 8080:8080 --env PORT=8080 asset-manager
```

## API Endpoints

### Assets
- `GET    /assets` — List assets (supports `limit` and `nextKey` query params)
- `GET    /assets/{id}` — Get asset by ID
- `POST   /assets` — Create asset
- `PATCH  /assets/{id}` — Patch asset fields
- `GET    /assets/{id}/publishRule` — Get publish rule
- `PATCH  /assets/{id}/publishRule` — Patch publish rule
- `POST   /assets/{id}/videos/{label}` — Add/update video variant
- `DELETE /assets/{id}/videos/{label}` — Delete video variant
- `POST   /assets/{id}/images` — Add image
- `DELETE /assets/{id}/images/{filename}` — Delete image

### Buckets
- `GET    /buckets` — List buckets (supports `limit` and `nextKey` query params)
- `GET    /buckets/{id}` — Get bucket by ID
- `POST   /buckets` — Create bucket
- `PATCH  /buckets/{id}` — Patch bucket fields
- `POST   /buckets/{id}/assets` — Add asset to bucket
- `DELETE /buckets/{id}/assets/{assetId}` — Remove asset from bucket

## Notes
- Uses DynamoDB tables `asset` and `bucket` (see `build.sh` for setup)
- Designed to work with other services in the hobby-streamer project
- See code for request/response schemas 