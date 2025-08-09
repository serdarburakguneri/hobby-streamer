# Asset Manager

GraphQL API for managing assets and buckets (DDD, Neo4j, event-driven, SQS).

## Endpoints
`POST /query` (GraphQL), `GET /` (playground), `GET /health`.

Schema lives under `internal/interfaces/graphql/`. Use the playground to explore.

## Running
```bash
cd backend/asset-manager && go run ./cmd/main.go
```

## Development
```bash
go run github.com/99designs/gqlgen generate
go build ./... && go test ./...
```
