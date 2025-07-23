package job

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

func downloadFromS3(ctx context.Context, s3Client *s3.Client, input string) (string, error) {
	if strings.HasPrefix(input, "s3://") {
		var localPath string
		var err error

		retryFunc := func(ctx context.Context) error {
			localPath, err = s3Client.Download(ctx, input)
			return err
		}

		retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 3)
		if retryErr != nil {
			logger.Get().WithError(retryErr).Error("Failed to download from S3 after retries", "input", input)
			return "", pkgerrors.NewExternalError("failed to download from S3", retryErr)
		}

		return localPath, nil
	}
	return input, nil
}

func uploadToS3(ctx context.Context, s3Client *s3.Client, localDir, s3Path string) error {
	if !strings.HasPrefix(s3Path, "s3://") {
		return pkgerrors.NewValidationError("output must be an S3 path", nil)
	}

	parts := strings.SplitN(s3Path[5:], "/", 2)
	if len(parts) != 2 {
		return pkgerrors.NewValidationError("invalid S3 path: "+s3Path, nil)
	}

	bucket := parts[0]
	keyPrefix := parts[1]
	manifestName := filepath.Base(keyPrefix)
	manifestDir := filepath.Dir(keyPrefix)

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		var s3Key string
		if relPath == manifestName {
			s3Key = keyPrefix
		} else {
			s3Key = manifestDir
			if s3Key != "" && !strings.HasSuffix(s3Key, "/") {
				s3Key += "/"
			}
			s3Key += relPath
		}

		retryFunc := func(ctx context.Context) error {
			return s3Client.Upload(ctx, bucket, s3Key, path)
		}

		retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 3)
		if retryErr != nil {
			logger.Get().WithError(retryErr).Error("Failed to upload file to S3 after retries", "local_file", path, "s3_key", s3Key)
			return pkgerrors.NewExternalError("failed to upload file to S3", retryErr)
		}

		logger.Get().Info("Successfully uploaded file to S3", "local_file", path, "s3_key", s3Key)
		return nil
	})
	if err != nil {
		return pkgerrors.NewInternalError("failed to walk output directory", err)
	}
	return nil
}
