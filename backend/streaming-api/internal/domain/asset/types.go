package asset

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
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

type PublishRuleValue struct {
	publishAt   *time.Time
	unpublishAt *time.Time
	regions     []string
	ageRating   *string
}

func NewPublishRuleValue(publishAt, unpublishAt *time.Time, regions []string, ageRating *string) (*PublishRuleValue, error) {
	if publishAt != nil && unpublishAt != nil {
		if publishAt.After(*unpublishAt) {
			return nil, ErrInvalidPublishDates
		}
	}

	if len(regions) > 50 {
		return nil, ErrTooManyRegions
	}

	validatedRegions := make([]string, 0, len(regions))
	for _, region := range regions {
		if len(region) > 10 {
			return nil, ErrInvalidRegion
		}

		regionRegex := regexp.MustCompile(`^[A-Z]{2,3}$`)
		if !regionRegex.MatchString(region) {
			return nil, ErrInvalidRegion
		}

		validatedRegions = append(validatedRegions, strings.ToUpper(region))
	}

	if ageRating != nil {
		validRatings := map[string]bool{
			constants.AgeRatingG: true, constants.AgeRatingPG: true, constants.AgeRatingPG13: true, constants.AgeRatingR: true, constants.AgeRatingNC17: true,
			constants.AgeRatingTVY: true, constants.AgeRatingTVY7: true, constants.AgeRatingTVG: true, constants.AgeRatingTVPG: true, constants.AgeRatingTV14: true, constants.AgeRatingTVMA: true,
		}

		if !validRatings[*ageRating] {
			return nil, ErrInvalidAgeRating
		}
	}

	return &PublishRuleValue{
		publishAt:   publishAt,
		unpublishAt: unpublishAt,
		regions:     validatedRegions,
		ageRating:   ageRating,
	}, nil
}

func (p PublishRuleValue) PublishAt() *time.Time {
	return p.publishAt
}

func (p PublishRuleValue) UnpublishAt() *time.Time {
	return p.unpublishAt
}

func (p PublishRuleValue) Regions() []string {
	return p.regions
}

func (p PublishRuleValue) AgeRating() *string {
	return p.ageRating
}

func (p PublishRuleValue) Equals(other PublishRuleValue) bool {
	return p.publishAt == other.publishAt &&
		p.unpublishAt == other.unpublishAt &&
		p.ageRating == other.ageRating
}

type VideoID struct {
	value string
}

func NewVideoID(value string) (*VideoID, error) {
	if value == "" {
		return nil, ErrInvalidVideoID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidVideoID
	}

	return &VideoID{value: value}, nil
}

func (v VideoID) Value() string {
	return v.value
}

func (v VideoID) Equals(other VideoID) bool {
	return v.value == other.value
}

type VideoType struct {
	value string
}

func NewVideoType(value string) (*VideoType, error) {
	if value == "" {
		return nil, ErrInvalidVideoType
	}

	validTypes := map[string]bool{
		constants.VideoTypeMain:      true,
		constants.VideoTypeTrailer:   true,
		constants.VideoTypeBehind:    true,
		constants.VideoTypeInterview: true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidVideoType
	}

	return &VideoType{value: value}, nil
}

func (v VideoType) Value() string {
	return v.value
}

func (v VideoType) Equals(other VideoType) bool {
	return v.value == other.value
}

type VideoFormat struct {
	value string
}

func NewVideoFormat(value string) (*VideoFormat, error) {
	if value == "" {
		return nil, ErrInvalidVideoFormat
	}

	validFormats := map[string]bool{
		"mp4":  true,
		"webm": true,
		"avi":  true,
		"mov":  true,
		"mkv":  true,
		"hls":  true,
		"dash": true,
	}

	if !validFormats[value] {
		return nil, ErrInvalidVideoFormat
	}

	return &VideoFormat{value: value}, nil
}

func (v VideoFormat) Value() string {
	return v.value
}

func (v VideoFormat) Equals(other VideoFormat) bool {
	return v.value == other.value
}

type ImageID struct {
	value string
}

func NewImageID(value string) (*ImageID, error) {
	if value == "" {
		return nil, ErrInvalidImageID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidImageID
	}

	return &ImageID{value: value}, nil
}

func (i ImageID) Value() string {
	return i.value
}

func (i ImageID) Equals(other ImageID) bool {
	return i.value == other.value
}

type FileName struct {
	value string
}

func NewFileName(value string) (*FileName, error) {
	if value == "" {
		return nil, ErrInvalidFileName
	}

	if len(value) > 255 {
		return nil, ErrInvalidFileName
	}

	fileNameRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !fileNameRegex.MatchString(value) {
		return nil, ErrInvalidFileName
	}

	return &FileName{value: value}, nil
}

func (f FileName) Value() string {
	return f.value
}

func (f FileName) Equals(other FileName) bool {
	return f.value == other.value
}

type ImageType struct {
	value string
}

func NewImageType(value string) (*ImageType, error) {
	if value == "" {
		return nil, ErrInvalidImageType
	}

	validTypes := map[string]bool{
		constants.ImageTypePoster:    true,
		constants.ImageTypeThumbnail: true,
		constants.ImageTypeBanner:    true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidImageType
	}

	return &ImageType{value: value}, nil
}

func (i ImageType) Value() string {
	return i.value
}

func (i ImageType) Equals(other ImageType) bool {
	return i.value == other.value
}

var (
	ErrInvalidS3Bucket     = errors.New("invalid S3 bucket name")
	ErrInvalidS3Key        = errors.New("invalid S3 key")
	ErrInvalidS3URL        = errors.New("invalid S3 URL")
	ErrInvalidDownloadURL  = errors.New("invalid download URL")
	ErrInvalidCDNPrefix    = errors.New("invalid CDN prefix")
	ErrInvalidStreamURL    = errors.New("invalid stream URL")
	ErrInvalidPublishDates = errors.New("invalid publish dates")
	ErrTooManyRegions      = errors.New("too many regions")
	ErrInvalidRegion       = errors.New("invalid region")
	ErrInvalidAgeRating    = errors.New("invalid age rating")
	ErrInvalidVideoID      = errors.New("invalid video ID")
	ErrInvalidVideoType    = errors.New("invalid video type")
	ErrInvalidVideoFormat  = errors.New("invalid video format")
	ErrInvalidImageID      = errors.New("invalid image ID")
	ErrInvalidFileName     = errors.New("invalid file name")
	ErrInvalidImageType    = errors.New("invalid image type")
)
