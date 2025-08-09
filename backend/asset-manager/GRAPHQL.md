# GraphQL API

Asset Manager exposes a GraphQL API (gqlgen) alongside REST.

## Endpoints
`POST /graphql`, `GET /playground` (dev).

## Schema
Types include `Asset`, `Bucket`, and pagination wrappers. Inputs for create/update.

## Development
```bash
go run github.com/99designs/gqlgen generate
```

## Notes
Use the playground to explore schema and queries. Neo4j backs hierarchical queries.