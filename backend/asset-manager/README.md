# Asset Manager Service

A GraphQL API for managing media assets and buckets. Stores asset metadata in Neo4j and supports parent-child relationships between assets. Includes basic integration with video processing workflows and authentication via Keycloak.

## Features

- GraphQL API for querying and managing assets and buckets
- Support for hierarchical asset types (e.g. series → season → episode)
- SQS integration for publishing transcode jobs and receiving status updates
- JWT authentication using Keycloak
- Typed error handling for better debugging


### Run Locally

```bash
cd backend/asset-manager
go run ./cmd/main.go
```

### Run with Docker

```bash
docker build -t asset-manager .
docker run -p 8080:8080 asset-manager
```

## API

### Endpoints

- `POST /graphql` – Main GraphQL endpoint
- `GET /playground` – GraphQL Playground UI (enabled in development)

## GraphQL Schema Overview

### Asset Types

- Movie
- Series / Season / Episode
- Documentary
- Music
- Podcast
- Trailer
- BehindTheScenes
- Interview

### Asset Status

- `draft` – Unpublished content
- `scheduled` – Scheduled to be published later
- `published` – Publicly available
- `expired` – Previously published, now expired

### Video Status

- `pending` – Waiting to be processed
- `analyzing` – Being analyzed for format, resolution, etc.
- `transcoding` – In the process of being transcoded
- `ready` – Fully processed and available
- `failed` – Processing failed

## Example Query

### List Assets

```graphql
query {
  assets(limit: 10) {
    items {
      id
      title
      status
      videos {
        label
        status
      }
    }
  }
}
```

## Development

### Generate GraphQL Code

```bash
go run github.com/99designs/gqlgen generate
```

### Build

```bash
go build ./...
```

### Run Tests

```bash
go test ./...
```