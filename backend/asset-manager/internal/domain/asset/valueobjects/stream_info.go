package valueobjects

import (
	"errors"
)

type StreamInfo struct {
	downloadURL *string
	cdnPrefix   *string
	url         *string
}

func NewStreamInfo(downloadURL, cdnPrefix, url *string) (*StreamInfo, error) {
	if downloadURL == nil && cdnPrefix == nil && url == nil {
		return nil, errors.New("at least one URL must be provided")
	}

	return &StreamInfo{
		downloadURL: downloadURL,
		cdnPrefix:   cdnPrefix,
		url:         url,
	}, nil
}

func (si StreamInfo) DownloadURL() *string {
	return si.downloadURL
}

func (si StreamInfo) CDNPrefix() *string {
	return si.cdnPrefix
}

func (si StreamInfo) URL() *string {
	return si.url
}

func (si StreamInfo) HasDownloadURL() bool {
	return si.downloadURL != nil && *si.downloadURL != ""
}

func (si StreamInfo) HasCDNPrefix() bool {
	return si.cdnPrefix != nil && *si.cdnPrefix != ""
}

func (si StreamInfo) HasURL() bool {
	return si.url != nil && *si.url != ""
}
