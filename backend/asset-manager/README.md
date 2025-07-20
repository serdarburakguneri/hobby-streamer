# Asset Manager Service

GraphQL API for managing media assets and buckets with Domain-Driven Design (DDD) architecture. Stores metadata in Neo4j with parent-child relationships and handles video workflows with event-driven processing.

## Architecture

**Domain Layer** (`internal/domain/`)
- `asset/` - Asset aggregate root with encapsulated videos
- `bucket/` - Bucket aggregate root for organizing assets

**Application Layer** (`internal/application/`)
- `asset/` - Asset application service orchestrating use cases
- `bucket/` - Bucket application service orchestrating use cases

**Infrastructure Layer** (`internal/infrastructure/`)
- `neo4j/` - Neo4j implementations for both asset and bucket repositories
- `sqs/` - SQS event publishing for domain events

**Interfaces Layer** (`internal/interfaces/graphql/`)
- GraphQL resolvers using application services

## Features

- **DDD Architecture**: Clean separation of concerns with domain-driven design
- **Aggregate Pattern**: Asset and Bucket as aggregate roots
- **Event-Driven**: Domain events for loose coupling
- **GraphQL API**: Type-safe GraphQL interface
- **Neo4j Storage**: Graph database for hierarchical relationships
- **SQS Integration**: Asynchronous event processing

## Running

```bash
cd backend/asset-manager && go run ./cmd/main.go
```

## API

- `POST /query` - GraphQL endpoint
- `GET /` - Interactive playground
- `GET /health` - Health check endpoint

## Domain Model

**Asset Aggregate Root:**
- Controls access to videos through aggregate methods
- Enforces business rules and consistency
- Publishes domain events for state changes

**Bucket Aggregate Root:**
- Organizes assets into logical groups
- Manages asset relationships within buckets
- Enforces ownership and access control

**Video Entity:**
- Part of Asset aggregate
- Cannot be modified outside asset context
- Supports multiple formats and statuses

## GraphQL API

### Asset Operations

```graphql
mutation CreateAsset {
  createAsset(input: {
    slug: "my-video"
    title: "My Video"
    type: "movie"
  }) {
    id
    slug
    title
    status
  }
}

query GetAssets {
  assets(limit: 10) {
    id
    slug
    title
    status
  }
}

query GetAsset {
  asset(id: "asset-id") {
    id
    slug
    title
    videos {
      id
      label
      status
    }
  }
}
```

### Bucket Operations

```graphql
mutation CreateBucket {
  createBucket(input: {
    name: "My Collection"
    slug: "my-collection"
    description: "A collection of my favorite videos"
  }) {
    id
    name
    slug
    assetCount
  }
}

query GetBuckets {
  buckets(limit: 10) {
    items {
      id
      name
      slug
      assetCount
    }
    hasMore
  }
}

query GetBucket {
  bucket(id: "bucket-id") {
    id
    name
    assets {
      id
      slug
      title
      status
    }
  }
}

mutation AddAssetToBucket {
  addAssetToBucket(input: {
    bucketId: "bucket-id"
    assetId: "asset-id"
    ownerId: "user-id"
  })
}

mutation RemoveAssetFromBucket {
  removeAssetFromBucket(input: {
    bucketId: "bucket-id"
    assetId: "asset-id"
    ownerId: "user-id"
  })
}
```

## Development

```bash
go run github.com/99designs/gqlgen generate
go build ./... && go test ./...
```
