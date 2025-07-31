# S3 Package

Shared Go library for S3-compatible storage. Handles upload/download, LocalStack support, context-aware, structured logging.

## Features
Download S3 objects to temp files, upload files/directories, prefix support, context-aware, LocalStack compatible, consistent error handling, logging.

## Quick Usage

### Create Client
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
ctx := context.Background()
client, err := s3.NewClient(ctx)
```

### Download
```go
localPath, err := client.Download(ctx, "s3://bucket-name/path/to/file.mp4")
defer os.Remove(localPath)
```

### Upload
```go
err := client.Upload(ctx, "/local/path/file.mp4", "bucket-name", "path/to/file.mp4")
```

### Upload Directory
```go
err := client.UploadDirectory(ctx, "/local/dir", "bucket-name", "path/prefix")
```

## Env Vars
AWS_ENDPOINT (default: http://localstack:4566), AWS_REGION (default: us-east-1), AWS_ACCESS_KEY_ID (default: test), AWS_SECRET_ACCESS_KEY (default: test)

## Notes
Temp files in system temp dir, AWS sessions/config handled internally, all operations respect context.Context.

