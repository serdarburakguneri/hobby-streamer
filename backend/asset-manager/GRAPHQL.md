# GraphQL API Implementation

This document describes the GraphQL API implementation for the asset-manager service.

## Overview

The asset-manager service now provides both REST and GraphQL APIs. The GraphQL API is built using gqlgen and provides a type-safe, efficient way to query and mutate assets and buckets.

## Features

- **Type-safe GraphQL API** using gqlgen
- **Neo4j integration** for graph-based queries
- **Hierarchical asset relationships** (series → seasons → episodes)
- **Flexible querying** with pagination support
- **GraphQL Playground** for development and testing

## API Endpoints

- **GraphQL Endpoint**: `POST /graphql`
- **GraphQL Playground**: `GET /playground` (development only)

## Schema

The GraphQL schema includes:

### Types

- `Asset` - Media assets with hierarchical relationships
- `Bucket` - Collections of assets
- `AssetPage` / `BucketPage` - Paginated results
- `AssetInput` / `BucketInput` - Input types for mutations

### Enums

- `AssetShape` - video, image, audio
- `AssetType` - movie, series, season, episode, etc.
- `BucketType` - playlist, collection, etc.

### Queries

- `asset(id: ID!)` - Get single asset
- `assets(limit: Int, nextKey: String)` - List assets with pagination
- `bucket(id: ID!)` - Get single bucket
- `buckets(limit: Int, nextKey: String)` - List buckets with pagination

### Mutations

- `createAsset(input: AssetInput!)` - Create new asset
- `updateAsset(id: ID!, input: AssetInput!)` - Update asset
- `deleteAsset(id: ID!)` - Delete asset (soft delete)
- `createBucket(input: BucketInput!)` - Create new bucket
- `updateBucket(id: ID!, input: BucketInput!)` - Update bucket
- `deleteBucket(id: ID!)` - Delete bucket (soft delete)

## Example Queries

### Get Assets with Pagination
```graphql
query {
  assets(limit: 10) {
    items {
      id
      title
      description
      shape
      type
      genre
      tags
      status
      createdAt
      updatedAt
    }
    nextKey
  }
}
```

### Get Asset with Relationships
```graphql
query {
  asset(id: "asset-123") {
    id
    title
    description
    parent {
      id
      title
    }
    children {
      id
      title
    }
    buckets {
      id
      name
    }
  }
}
```

### Create Asset
```graphql
mutation {
  createAsset(input: {
    title: "Sample Movie"
    description: "A sample movie description"
    shape: VIDEO
    type: MOVIE
    genre: "Action"
    tags: ["action", "adventure"]
    status: "draft"
  }) {
    id
    title
    createdAt
  }
}
```

## Neo4j Integration

The GraphQL API leverages Neo4j's graph capabilities for:

- **Hierarchical queries** - Finding parent/child relationships
- **Type-based queries** - Filtering by asset type (series, season, episode)
- **Genre-based queries** - Finding assets by genre
- **Tag-based queries** - Finding assets by tags
- **Shape-based queries** - Filtering by content shape (video, image, audio)

## Repository Layer

The implementation includes comprehensive repository interfaces:

### AssetRepository
- Basic CRUD operations
- Graph-specific methods for hierarchical queries
- Type and genre-based filtering
- Tag-based searching

### BucketRepository
- Basic CRUD operations
- Type-based filtering
- Asset-bucket relationship queries

## Service Layer

The service layer provides business logic and orchestrates repository operations:

- **AssetService** - Handles asset business logic
- **BucketService** - Handles bucket business logic

## Authentication

The GraphQL API uses the same authentication middleware as the REST API:

- JWT token validation via Keycloak
- Role-based access control
- Admin and user role support

## Development

### Running the Service

1. Start Neo4j:
   ```bash
   docker-compose up neo4j
   ```

2. Start the service:
   ```bash
   go run cmd/main.go
   ```

3. Access GraphQL Playground:
   ```
   http://localhost:8080/playground
   ```

### Code Generation

The GraphQL code is generated using gqlgen:

```bash
# Generate code from schema
go run github.com/99designs/gqlgen generate

# Generate code with specific config
go run github.com/99designs/gqlgen generate --config gqlgen.yml
```

### Schema Changes

To modify the GraphQL schema:

1. Edit `graph/schema.graphqls`
2. Run code generation: `go run github.com/99designs/gqlgen generate`
3. Implement new resolvers in `graph/schema.resolvers.go`
4. Update mappers in `graph/mappers.go` if needed

## Testing

The GraphQL API can be tested using:

- **GraphQL Playground** - Interactive testing interface
- **Postman** - REST client with GraphQL support
- **curl** - Command line testing
- **Integration tests** - Automated testing

## Performance Considerations

- **Pagination** - All list queries support pagination
- **Field selection** - Clients can select only needed fields
- **Neo4j optimization** - Queries are optimized for graph traversal
- **Caching** - Consider implementing Redis caching for frequently accessed data

## Future Enhancements

- **Subscriptions** - Real-time updates for asset changes
- **File uploads** - Direct file upload via GraphQL
- **Advanced filtering** - Complex filter combinations
- **Bulk operations** - Batch create/update operations
- **Search** - Full-text search integration 