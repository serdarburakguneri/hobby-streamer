package valueobjects

import (
	"net/url"
	"regexp"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type S3ObjectValue struct {
	bucket string
	key    string
	url    string
}

func NewS3ObjectValue(bucket, key, urlStr string) (*S3ObjectValue, error) {
	if bucket == "" {
		return nil, ErrInvalidS3Bucket
	}

	if key == "" {
		return nil, ErrInvalidS3Key
	}

	if urlStr != "" {
		if _, err := url.Parse(urlStr); err != nil {
			return nil, ErrInvalidS3URL
		}
	}

	bucketRegex := regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	if !bucketRegex.MatchString(bucket) {
		return nil, ErrInvalidS3Bucket
	}

	return &S3ObjectValue{
		bucket: bucket,
		key:    key,
		url:    urlStr,
	}, nil
}

func (s S3ObjectValue) Bucket() string {
	return s.bucket
}

func (s S3ObjectValue) Key() string {
	return s.key
}

func (s S3ObjectValue) URL() string {
	return s.url
}

func (s S3ObjectValue) Equals(other S3ObjectValue) bool {
	return s.bucket == other.bucket && s.key == other.key && s.url == other.url
}

var (
	ErrInvalidS3Bucket = pkgerrors.NewValidationError("invalid S3 bucket name", nil)
	ErrInvalidS3Key    = pkgerrors.NewValidationError("invalid S3 key", nil)
	ErrInvalidS3URL    = pkgerrors.NewValidationError("invalid S3 URL", nil)
)
