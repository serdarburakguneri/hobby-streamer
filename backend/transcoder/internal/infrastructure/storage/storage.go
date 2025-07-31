package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

type Storage struct {
	client *s3.Client
}

func NewStorage(client *s3.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) Download(ctx context.Context, input string) (string, error) {
	if strings.HasPrefix(input, "s3://") {
		return s.client.Download(ctx, input)
	}
	return input, nil
}

func (s *Storage) CreateDir(path string) error {
	return os.MkdirAll(path, 0750)
}

func (s *Storage) Remove(path string) error {
	return os.Remove(path)
}

func (s *Storage) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (s *Storage) Upload(ctx context.Context, localDir, s3Path string) error {
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
			return s.client.Upload(ctx, path, bucket, s3Key)
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
