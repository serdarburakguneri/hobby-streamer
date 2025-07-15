# S3 Package

A shared Go package for S3 operations in the hobby-streamer project. Provides a simple interface for uploading and downloading files to/from S3.

## Features
- Download files from S3 URLs to local temporary files
- Upload local files to S3
- Upload entire directories to S3
- Automatic AWS session management with LocalStack support
- Consistent logging and error handling

## Usage

### Creating a Client
```go
import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"

ctx := context.Background()
client, err := s3.NewClient(ctx)
if err != nil {
    // handle error
}
```

### Downloading from S3
```go
localPath, err := client.Download(ctx, "s3://bucket-name/path/to/file.mp4")
if err != nil {
    // handle error
}
defer os.Remove(localPath) // Clean up temporary file
```

### Uploading to S3
```go
err := client.Upload(ctx, "/local/path/file.mp4", "bucket-name", "path/to/file.mp4")
if err != nil {
    // handle error
}
```

### Uploading a Directory
```go
err := client.UploadDirectory(ctx, "/local/directory", "bucket-name", "path/prefix")
if err != nil {
    // handle error
}
```

## Environment Variables
- `AWS_ENDPOINT`: Custom endpoint for AWS services (default: `http://localstack:4566` for LocalStack)
- `AWS_REGION`: AWS region (default: `us-east-1`)
- `AWS_ACCESS_KEY_ID`: AWS access key (default: `test` for LocalStack)
- `AWS_SECRET_ACCESS_KEY`: AWS secret key (default: `test` for LocalStack)

## Notes
- Downloaded files are stored in the system's temporary directory
- The package automatically handles AWS session creation with LocalStack support
- All operations are context-aware for proper cancellation
- Consistent error handling and logging throughout 