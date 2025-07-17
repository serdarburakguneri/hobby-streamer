# S3 Package

A shared Go library for working with S3-compatible storage. Simplifies common file operations and supports LocalStack for local development.

## Features

- Download S3 objects to local temporary files
- Upload local files to S3
- Upload entire directories to S3 with prefix support
- Context-aware operations with cancellation support
- Built-in LocalStack compatibility
- Structured logging and consistent error handling

---

## Usage

### Create Client

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"

ctx := context.Background()
client, err := s3.NewClient(ctx)
if err != nil {
    // handle error
}
```

---

### Download a File from S3

```go
localPath, err := client.Download(ctx, "s3://bucket-name/path/to/file.mp4")
if err != nil {
    // handle error
}
defer os.Remove(localPath) // Clean up after use
```

---

### Upload a File to S3

```go
err := client.Upload(ctx, "/local/path/file.mp4", "bucket-name", "path/to/file.mp4")
if err != nil {
    // handle error
}
```

---

### Upload a Directory to S3

```go
err := client.UploadDirectory(ctx, "/local/dir", "bucket-name", "path/prefix")
if err != nil {
    // handle error
}
```

---

## Environment Variables

| Variable               | Description                            | Default                  |
|------------------------|----------------------------------------|--------------------------|
| `AWS_ENDPOINT`         | Custom AWS endpoint (for LocalStack)   | `http://localstack:4566` |
| `AWS_REGION`           | AWS region                             | `us-east-1`              |
| `AWS_ACCESS_KEY_ID`    | AWS access key                         | `test`                   |
| `AWS_SECRET_ACCESS_KEY`| AWS secret key                         | `test`                   |

---

## Notes

- Downloaded files are saved to the system’s temporary directory
- The package handles AWS session creation internally
- All operations support `context.Context` for cancellation
