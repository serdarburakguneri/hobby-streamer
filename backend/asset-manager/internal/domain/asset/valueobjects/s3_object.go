package valueobjects

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type S3Object struct {
	bucket string
	key    string
	url    string
}

func NewS3Object(bucket, key, url string) (*S3Object, error) {
	if bucket == "" {
		return nil, errors.New("bucket cannot be empty")
	}

	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if len(bucket) > 63 {
		return nil, errors.New("bucket name too long")
	}

	if len(key) > 1024 {
		return nil, errors.New("key too long")
	}

	matched, err := regexp.MatchString(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`, bucket)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("invalid bucket name format")
	}

	return &S3Object{
		bucket: bucket,
		key:    key,
		url:    url,
	}, nil
}

func NewS3ObjectFromURL(url string) (*S3Object, error) {
	if url == "" {
		return nil, errors.New("URL cannot be empty")
	}

	bucket, key, err := parseS3URL(url)
	if err != nil {
		return nil, err
	}

	return NewS3Object(bucket, key, url)
}

func (s3 S3Object) Bucket() string {
	return s3.bucket
}

func (s3 S3Object) Key() string {
	return s3.key
}

func (s3 S3Object) URL() string {
	return s3.url
}

func (s3 S3Object) BuildS3URL() string {
	if s3.url != "" {
		return s3.url
	}
	return fmt.Sprintf("s3://%s/%s", s3.bucket, s3.key)
}

func (s3 S3Object) Equals(other S3Object) bool {
	return s3.bucket == other.bucket && s3.key == other.key
}

func parseS3URL(url string) (string, string, error) {
	if strings.HasPrefix(url, "s3://") {
		parts := strings.SplitN(url[5:], "/", 2)
		if len(parts) != 2 {
			return "", "", errors.New("invalid S3 URL format")
		}
		return parts[0], parts[1], nil
	}

	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(url, "/")
		if len(parts) < 4 {
			return "", "", errors.New("invalid S3 HTTPS URL format")
		}
		bucket := parts[2]
		key := strings.Join(parts[3:], "/")
		return bucket, key, nil
	}

	return "", "", errors.New("unsupported URL format")
}
