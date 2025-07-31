# Asset Manager Service

GraphQL API for managing media assets and buckets, using DDD, Neo4j for metadata, event-driven workflows, and SQS for async processing.

## Features
DDD architecture, aggregate roots (asset, bucket), event-driven, GraphQL API, Neo4j storage, SQS integration, type-safe schema, parent-child relationships, asset-bucket organization.

## Architecture
- Domain: `internal/domain/` (asset, bucket)
- Application: `internal/application/` (asset, bucket services)
- Infrastructure: `internal/infrastructure/` (neo4j, sqs)
- Interfaces: `internal/interfaces/graphql/` (resolvers)

## API
- `POST /query` (GraphQL)
- `GET /` (playground)
- `GET /health`

## GraphQL Usage Examples

### Asset
```graphql
mutation { createAsset(input: {slug: "my-video", title: "My Video", type: "movie"}) { id slug title status } }
query { assets(limit: 10) { id slug title status } }
query { asset(id: "asset-id") { id slug title videos { id label status } } }
```

### Bucket
```graphql
mutation { createBucket(input: {name: "My Collection", slug: "my-collection", description: "A collection"}) { id name slug assetCount } }
query { buckets(limit: 10) { items { id name slug assetCount } hasMore } }
query { bucket(id: "bucket-id") { id name assets { id slug title status } } }
mutation { addAssetToBucket(input: {bucketId: "bucket-id", assetId: "asset-id", ownerId: "user-id"}) }
mutation { removeAssetFromBucket(input: {bucketId: "bucket-id", assetId: "asset-id", ownerId: "user-id"}) }
```

## Running
```bash
cd backend/asset-manager && go run ./cmd/main.go
```

## Development
```bash
go run github.com/99designs/gqlgen generate
go build ./... && go test ./...
```
