# Asset Manager Service

GraphQL API for managing media assets and buckets in the Hobby Streamer project. It stores asset metadata in Neo4j and supports parent-child relationships (e.g. series → season → episode). Also handles integration with video workflows and authentication via Keycloak.

---

## Features

- GraphQL API for asset and bucket management
- Hierarchical asset support (Movie, Series, Episode, etc.)
- SQS integration for triggering transcode jobs and receiving updates
- JWT-based authentication (via Keycloak)
- Typed error handling for more informative responses

---

## Running Locally

```bash
cd backend/asset-manager
go run ./cmd/main.go
```

---

## Docker

```bash
docker build -t asset-manager .
docker run -p 8080:8080 asset-manager
```

---

## API Endpoints

- `POST /graphql` — GraphQL query endpoint  
- `GET /playground` — Interactive GraphQL Playground (dev only)

---

## GraphQL Schema

### Asset Types

- `Movie`
- `Series`, `Season`, `Episode`
- `Documentary`
- `Music`
- `Podcast`
- `Trailer`
- `BehindTheScenes`
- `Interview`

### Asset Status

- `draft` — Not yet published  
- `scheduled` — Set to be published later  
- `published` — Publicly available  
- `expired` — No longer available

### Video Status

- `pending` — Awaiting processing  
- `analyzing` — Format/resolution check  
- `transcoding` — Actively being processed  
- `ready` — Available for playback  
- `failed` — Processing failed

---

## Example Query

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

---

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

---

> ⚠️ This service is still evolving — primarily used for local development and exploring asset modeling patterns.