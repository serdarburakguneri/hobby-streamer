# S3 Package

A shared helper library for working with S3-compatible storage in Go. Wraps common upload/download operations with built-in support for LocalStack, structured logging, and context handling.

---

## Features

- Download S3 objects to local temp files  
- Upload single files or entire directories  
- Prefix support for organizing uploads  
- Context-aware operations (with cancel support)  
- Fully compatible with LocalStack  
- Consistent error handling and logging

---

## Quick Usage

### Create an S3 Client

```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"

ctx := context.Background()
client, err := s3.NewClient(ctx)
if err != nil {
    // Handle init error
}
```

---

### Download a File

```go
localPath, err := client.Download(ctx, "s3://bucket-name/path/to/file.mp4")
if err != nil {
    // Handle download error
}
defer os.Remove(localPath) // Clean up temp file when done
```

---

### Upload a File

```go
err := client.Upload(ctx, "/local/path/file.mp4", "bucket-name", "path/to/file.mp4")
if err != nil {
    // Handle upload error
}
```

---

### Upload a Directory

```go
err := client.UploadDirectory(ctx, "/local/dir", "bucket-name", "path/prefix")
if err != nil {
    // Handle upload error
}
```

---

## Environment Variables

| Variable                 | Purpose                                | Default                  |
|--------------------------|----------------------------------------|--------------------------|
| `AWS_ENDPOINT`           | Custom endpoint (e.g. LocalStack)      | `http://localstack:4566` |
| `AWS_REGION`             | AWS region                             | `us-east-1`              |
| `AWS_ACCESS_KEY_ID`      | S3 access key                          | `test`                   |
| `AWS_SECRET_ACCESS_KEY`  | S3 secret key                          | `test`                   |

---

## Notes

- Temp files are stored in the system’s default temporary directory
- AWS sessions and config are handled internally — no setup needed
- All operations respect `context.Context`, so cancellation and timeouts work as expected

---

