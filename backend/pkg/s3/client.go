package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Client struct {
	client *s3.S3
	logger *logger.Logger
}

func NewClient(ctx context.Context) (*Client, error) {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localstack:4566"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(awsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	})
	if err != nil {
		return nil, errors.NewInternalError("failed to create AWS session", err)
	}

	return &Client{
		client: s3.New(sess),
		logger: logger.WithService("s3-client"),
	}, nil
}

func (c *Client) Download(ctx context.Context, s3URL string) (string, error) {
	log := c.logger.WithContext(ctx)

	if !strings.HasPrefix(s3URL, "s3://") {
		return "", errors.NewInternalError(fmt.Sprintf("invalid S3 URL: %s", s3URL), nil)
	}

	parts := strings.SplitN(s3URL[5:], "/", 2)
	if len(parts) != 2 {
		return "", errors.NewInternalError(fmt.Sprintf("invalid S3 URL format: %s", s3URL), nil)
	}

	bucket := parts[0]
	key := parts[1]

	tempDir := os.TempDir()
	filename := filepath.Base(key)
	if filename == "" {
		filename = fmt.Sprintf("file_%d", time.Now().Unix())
	}
	localPath := filepath.Join(tempDir, filename)

	log.Info("Downloading from S3", "bucket", bucket, "key", key, "local_path", localPath)

	file, err := os.Create(localPath)
	if err != nil {
		return "", errors.NewInternalError("failed to create local file", err)
	}
	defer file.Close()

	result, err := c.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if removeErr := os.Remove(localPath); removeErr != nil {
			log.WithError(removeErr).Error("Failed to remove local file after S3 download error")
		}
		return "", errors.NewInternalError("failed to get object from S3", err)
	}
	defer result.Body.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		if removeErr := os.Remove(localPath); removeErr != nil {
			log.WithError(removeErr).Error("Failed to remove local file after copy error")
		}
		return "", errors.NewInternalError("failed to write file to disk", err)
	}

	log.Info("Successfully downloaded from S3", "local_path", localPath)
	return localPath, nil
}

func (c *Client) Upload(ctx context.Context, localPath, bucket, key string) error {
	log := c.logger.WithContext(ctx)

	file, err := os.Open(localPath)
	if err != nil {
		return errors.NewInternalError("failed to open local file", err)
	}
	defer file.Close()

	log.Info("Uploading to S3", "bucket", bucket, "key", key, "local_path", localPath)

	_, err = c.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return errors.NewInternalError("failed to upload object to S3", err)
	}

	log.Info("Successfully uploaded to S3", "bucket", bucket, "key", key)
	return nil
}

func (c *Client) UploadDirectory(ctx context.Context, localDir, bucket, keyPrefix string) error {
	log := c.logger.WithContext(ctx)

	files, err := filepath.Glob(filepath.Join(localDir, "*"))
	if err != nil {
		return errors.NewInternalError("failed to find files in directory", err)
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			fileName := filepath.Base(file)
			key := keyPrefix + "/" + fileName

			err = c.Upload(ctx, file, bucket, key)
			if err != nil {
				log.WithError(err).Error("Failed to upload file", "file", file, "bucket", bucket, "key", key)
				return err
			}
		}
	}

	return nil
}
