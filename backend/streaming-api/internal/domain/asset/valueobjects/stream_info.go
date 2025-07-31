package valueobjects

import (
	"net/url"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type StreamInfoValue struct {
	downloadURL *string
	cdnPrefix   *string
	url         *string
}

func NewStreamInfoValue(downloadURL, cdnPrefix, urlStr *string) (*StreamInfoValue, error) {
	if downloadURL != nil {
		if _, err := url.Parse(*downloadURL); err != nil {
			return nil, ErrInvalidDownloadURL
		}
	}

	if cdnPrefix != nil {
		if _, err := url.Parse(*cdnPrefix); err != nil {
			return nil, ErrInvalidCDNPrefix
		}
	}

	if urlStr != nil {
		if _, err := url.Parse(*urlStr); err != nil {
			return nil, ErrInvalidStreamURL
		}
	}

	return &StreamInfoValue{
		downloadURL: downloadURL,
		cdnPrefix:   cdnPrefix,
		url:         urlStr,
	}, nil
}

func (s StreamInfoValue) DownloadURL() *string {
	return s.downloadURL
}

func (s StreamInfoValue) CDNPrefix() *string {
	return s.cdnPrefix
}

func (s StreamInfoValue) URL() *string {
	return s.url
}

func (s StreamInfoValue) Equals(other StreamInfoValue) bool {
	return s.downloadURL == other.downloadURL &&
		s.cdnPrefix == other.cdnPrefix &&
		s.url == other.url
}

var (
	ErrInvalidDownloadURL = pkgerrors.NewValidationError("invalid download URL", nil)
	ErrInvalidCDNPrefix   = pkgerrors.NewValidationError("invalid CDN prefix", nil)
	ErrInvalidStreamURL   = pkgerrors.NewValidationError("invalid stream URL", nil)
)
