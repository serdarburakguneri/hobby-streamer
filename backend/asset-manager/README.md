# Asset Manager Service

GraphQL API for managing media assets and buckets. Stores metadata in Neo4j with parent-child relationships (series → season → episode). Handles video workflows and Keycloak authentication.

## Features

GraphQL API, hierarchical assets, SQS integration, JWT auth, typed error handling.

## Running

```bash
# Local
cd backend/asset-manager && go run ./cmd/main.go

# Docker
docker build -t asset-manager . && docker run -p 8080:8080 asset-manager
```

## API

- `POST /graphql` — GraphQL endpoint  
- `GET /playground` — Interactive playground (dev)

## Schema

**Asset Types:** Movie, Series, Season, Episode, Documentary, Music, Podcast, Trailer, BehindTheScenes, Interview

**Asset Status:** draft, scheduled, published, expired

**Video Status:** pending, analyzing, transcoding, ready, failed

## Example Query

```graphql
query {
  assets(limit: 10) {
    items {
      id
      title
      status
      videos { label status }
    }
  }
}
```

## Development

```bash
go run github.com/99designs/gqlgen generate
go build ./... && go test ./...
```

> ⚠️ Evolving service — local development and asset modeling exploration.