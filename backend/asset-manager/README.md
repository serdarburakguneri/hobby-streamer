# Asset Manager Service

A simple GraphQL API service for managing media assets and buckets in the hobby streaming project. Uses Neo4j for data storage and supports hierarchical asset relationships.

## What it does

- GraphQL API for assets and buckets
- Asset status management (draft/published)
- Video processing status tracking
- Hierarchical asset relationships
- Neo4j graph database backend
- JWT authentication via Keycloak

## Quick Start

### Prerequisites
- Go 1.23+
- Neo4j database
- Keycloak server

### Environment Variables
```
PORT=8080
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=password
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=hobby-realm
KEYCLOAK_CLIENT_ID=asset-manager
ENV=development
```

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

## API Endpoints

- `POST /graphql` - GraphQL endpoint
- `GET /playground` - GraphQL Playground (development only)

## GraphQL Schema

### Asset Types
- Movie, Series, Season, Episode
- Documentary, Music, Podcast
- Trailer, BehindTheScenes, Interview

### Asset Shapes
- Video, Image, Audio, Document

### Asset Status
- `draft` - Default status for assets without publish rules
- `scheduled` - Assets scheduled for future publication
- `published` - Currently published assets
- `expired` - Assets that have passed their unpublish date

### Video Status
- `pending` - Default status for new videos
- `analyzing` - Video being analyzed
- `transcoding` - Video being transcoded
- `ready` - Video processing complete
- `failed` - Video processing failed

## Example Queries

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

### Create Asset
```graphql
mutation {
  createAsset(input: {
    title: "Sample Movie"
    shape: VIDEO
    type: MOVIE
    genre: "action"
  }) {
    id
    title
    status
  }
}
```

### Update Asset with JSON Patch
```graphql
mutation {
  patchAsset(id: "123", patches: [
    { op: "replace", path: "/title", value: "Updated Title" }
  ]) {
    id
    title
  }
}
```

### Update Video Status
```graphql
mutation {
  updateVideoStatus(id: "123", label: "main", status: "transcoding") {
    id
    videos {
      label
      status
    }
  }
}
```

## Project Structure

```
cmd/main.go          # Application entry point
internal/
├── asset/           # Asset domain logic
├── bucket/          # Bucket domain logic
graph/
├── schema.graphqls  # GraphQL schema
└── schema.resolvers.go # Resolvers
```

## Development

Generate GraphQL code:
```bash
go run github.com/99designs/gqlgen generate
```

Build:
```bash
go build ./...
```

Run tests:
```bash
go test ./...
``` 