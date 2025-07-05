# Asset Manager Service

The Asset Manager service provides RESTful APIs to manage media assets and buckets, including asset creation, metadata management, image/video association, and bucket organization. The service integrates with Keycloak for authentication and implements role-based authorization.

## Features
- CRUD operations for assets and buckets
- Associate images and video variants with assets
- Organize assets into buckets
- DynamoDB-backed storage
- **JWT-based authentication via Keycloak**
- **Role-based authorization (admin/user roles)**

## Authentication & Authorization

This service uses JWT tokens from Keycloak for authentication and implements role-based access control:

### Roles
- **`admin`**: Full access to all endpoints (create, read, update, delete)
- **`user`**: Read-only access to assets and buckets

### Authentication
- All endpoints require a valid JWT token in the `Authorization` header
- Format: `Authorization: Bearer <jwt_token>`
- Tokens are validated against Keycloak's JWKS endpoint

### Authorization Matrix

| Endpoint | Method | Admin | User | Description |
|----------|--------|-------|------|-------------|
| `/assets` | GET | ✅ | ✅ | List assets |
| `/assets/{id}` | GET | ✅ | ✅ | Get asset by ID |
| `/assets` | POST | ✅ | ❌ | Create asset |
| `/assets/{id}` | PATCH | ✅ | ❌ | Update asset |
| `/assets/{id}/publishRule` | GET | ✅ | ✅ | Get publish rule |
| `/assets/{id}/publishRule` | PATCH | ✅ | ❌ | Update publish rule |
| `/assets/{id}/videos/{label}` | POST | ✅ | ❌ | Add/update video variant |
| `/assets/{id}/videos/{label}` | DELETE | ✅ | ❌ | Delete video variant |
| `/assets/{id}/images` | POST | ✅ | ❌ | Add image |
| `/assets/{id}/images/{filename}` | DELETE | ✅ | ❌ | Delete image |
| `/buckets` | GET | ✅ | ✅ | List buckets |
| `/buckets/{id}` | GET | ✅ | ✅ | Get bucket by ID |
| `/buckets` | POST | ✅ | ❌ | Create bucket |
| `/buckets/{id}` | PATCH | ✅ | ❌ | Update bucket |
| `/buckets/{id}/assets` | POST | ✅ | ❌ | Add asset to bucket |
| `/buckets/{id}/assets/{assetId}` | DELETE | ✅ | ❌ | Remove asset from bucket |

## Requirements
- Go 1.22+
- LocalStack (for local AWS emulation)
- Docker (optional, for containerized runs)
- **Keycloak server** (for authentication)

## Environment Variables
- `PORT`: HTTP port to listen on (default: 8080)
- `KEYCLOAK_URL`: Keycloak server URL (e.g., `http://localhost:8080`)
- `KEYCLOAK_REALM`: Keycloak realm name (e.g., `hobby-streamer`)
- `KEYCLOAK_CLIENT_ID`: Keycloak client ID (e.g., `asset-manager`)
- AWS credentials and region (for DynamoDB access, handled by LocalStack or your AWS profile)

## Running Locally

### 1. Start LocalStack and create DynamoDB tables (see project root `build.sh`)

### 2. Start Keycloak (see project root `docker-compose.yml`)

### 3. Run the service:
```sh
cd backend/asset-manager
PORT=8080 \
KEYCLOAK_URL=http://localhost:8080 \
KEYCLOAK_REALM=hobby-streamer \
KEYCLOAK_CLIENT_ID=asset-manager \
go run ./cmd/main.go
```

### 4. Or build and run with Docker:
```sh
docker build -t asset-manager .
docker run -p 8080:8080 \
  --env PORT=8080 \
  --env KEYCLOAK_URL=http://localhost:8080 \
  --env KEYCLOAK_REALM=hobby-streamer \
  --env KEYCLOAK_CLIENT_ID=asset-manager \
  asset-manager
```

## API Endpoints

### Assets
- `GET    /assets` — List assets (supports `limit` and `nextKey` query params)
- `GET    /assets/{id}` — Get asset by ID
- `POST   /assets` — Create asset *(admin only)*
- `PATCH  /assets/{id}` — Patch asset fields *(admin only)*
- `GET    /assets/{id}/publishRule` — Get publish rule
- `PATCH  /assets/{id}/publishRule` — Patch publish rule *(admin only)*
- `POST   /assets/{id}/videos/{label}` — Add/update video variant *(admin only)*
- `DELETE /assets/{id}/videos/{label}` — Delete video variant *(admin only)*
- `POST   /assets/{id}/images` — Add image *(admin only)*
- `DELETE /assets/{id}/images/{filename}` — Delete image *(admin only)*

### Buckets
- `GET    /buckets` — List buckets (supports `limit` and `nextKey` query params)
- `GET    /buckets/{id}` — Get bucket by ID
- `POST   /buckets` — Create bucket *(admin only)*
- `PATCH  /buckets/{id}` — Patch bucket fields *(admin only)*
- `POST   /buckets/{id}/assets` — Add asset to bucket *(admin only)*
- `DELETE /buckets/{id}/assets/{assetId}` — Remove asset from bucket *(admin only)*

## Authentication Examples

### Getting a JWT Token
```bash
# Using curl to get a token from Keycloak
curl -X POST http://localhost:8080/realms/hobby-streamer/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=asset-manager" \
  -d "username=adminuser" \
  -d "password=password"
```

### Using the Token
```bash
# Example API call with JWT token
curl -H "Authorization: Bearer <your_jwt_token>" \
  http://localhost:8080/assets
```

## Error Responses

The service returns standardized error responses:

```json
{
  "error": "Insufficient permissions"
}
```

Common error scenarios:
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User lacks required role for the endpoint
- `400 Bad Request`: Invalid request body or parameters

## Notes
- Uses DynamoDB tables `asset` and `bucket` (see `build.sh` for setup)
- Designed to work with other services in the hobby-streamer project
- JWT tokens are validated against Keycloak's JWKS endpoint with caching
- See code for request/response schemas 