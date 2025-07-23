package asset

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type ImageType string

const (
	ImageTypePoster    ImageType = "poster"
	ImageTypeBackdrop  ImageType = "backdrop"
	ImageTypeThumbnail ImageType = "thumbnail"
	ImageTypeLogo      ImageType = "logo"
)

func NewImageType(value string) (ImageType, error) {
	validTypes := map[string]bool{
		"poster":    true,
		"backdrop":  true,
		"thumbnail": true,
		"logo":      true,
	}

	if !validTypes[value] {
		return "", ErrInvalidImageType
	}

	return ImageType(value), nil
}

type S3Object struct {
	bucket string
	key    string
	url    string
}

func NewS3Object(bucket, key, urlStr string) (*S3Object, error) {
	if bucket == "" {
		return nil, ErrInvalidS3Bucket
	}

	if key == "" {
		return nil, ErrInvalidS3Key
	}

	if urlStr == "" {
		return nil, ErrInvalidS3URL
	}

	if _, err := url.Parse(urlStr); err != nil {
		return nil, ErrInvalidS3URL
	}

	bucketRegex := regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	if !bucketRegex.MatchString(bucket) {
		return nil, ErrInvalidS3Bucket
	}

	if len(bucket) < 3 || len(bucket) > 63 {
		return nil, ErrInvalidS3Bucket
	}

	if len(key) > 1024 {
		return nil, ErrInvalidS3Key
	}

	return &S3Object{
		bucket: bucket,
		key:    key,
		url:    urlStr,
	}, nil
}

func (s S3Object) Bucket() string {
	return s.bucket
}

func (s S3Object) Key() string {
	return s.key
}

func (s S3Object) URL() string {
	return s.url
}

func (s S3Object) Equals(other S3Object) bool {
	return s.bucket == other.bucket &&
		s.key == other.key &&
		s.url == other.url
}

func (s S3Object) BuildS3URL() string {
	return fmt.Sprintf("s3://%s/%s", s.bucket, s.key)
}

func BuildS3URL(bucket, key string) string {
	return fmt.Sprintf("s3://%s/%s", bucket, key)
}

func BuildHLSOutputKey(assetID, quality string) string {
	return fmt.Sprintf("%s/hls/%s/playlist.m3u8", assetID, quality)
}

func BuildDASHOutputKey(assetID, quality string) string {
	return fmt.Sprintf("%s/dash/%s/manifest.mpd", assetID, quality)
}

type StreamInfo struct {
	downloadURL *string
	cdnPrefix   *string
	url         *string
}

func NewStreamInfo(downloadURL, cdnPrefix, urlStr *string) (*StreamInfo, error) {
	if downloadURL != nil {
		if _, err := url.Parse(*downloadURL); err != nil {
			return nil, ErrInvalidStreamInfoURL
		}
	}

	if cdnPrefix != nil {
		if _, err := url.Parse(*cdnPrefix); err != nil {
			return nil, ErrInvalidStreamInfoURL
		}
	}

	if urlStr != nil {
		if _, err := url.Parse(*urlStr); err != nil {
			return nil, ErrInvalidStreamInfoURL
		}
	}

	return &StreamInfo{
		downloadURL: downloadURL,
		cdnPrefix:   cdnPrefix,
		url:         urlStr,
	}, nil
}

func (s StreamInfo) DownloadURL() *string {
	return s.downloadURL
}

func (s StreamInfo) CDNPrefix() *string {
	return s.cdnPrefix
}

func (s StreamInfo) URL() *string {
	return s.url
}

func (s StreamInfo) Equals(other StreamInfo) bool {
	return s.downloadURL == other.downloadURL &&
		s.cdnPrefix == other.cdnPrefix &&
		s.url == other.url
}

type TranscodingInfo struct {
	jobID       string
	progress    float64
	outputURL   string
	error       *string
	completedAt *time.Time
}

func NewTranscodingInfo(jobID string, progress float64, outputURL string, errorMsg *string, completedAt *time.Time) (*TranscodingInfo, error) {
	if jobID == "" {
		return nil, ErrInvalidTranscodingJobID
	}

	if progress < 0 || progress > 100 {
		return nil, ErrInvalidTranscodingProgress
	}

	if outputURL != "" {
		if _, err := url.Parse(outputURL); err != nil {
			return nil, ErrInvalidTranscodingOutputURL
		}
	}

	return &TranscodingInfo{
		jobID:       jobID,
		progress:    progress,
		outputURL:   outputURL,
		error:       errorMsg,
		completedAt: completedAt,
	}, nil
}

func (t TranscodingInfo) JobID() string {
	return t.jobID
}

func (t TranscodingInfo) Progress() float64 {
	return t.progress
}

func (t TranscodingInfo) OutputURL() string {
	return t.outputURL
}

func (t TranscodingInfo) Error() *string {
	return t.error
}

func (t TranscodingInfo) CompletedAt() *time.Time {
	return t.completedAt
}

func (t TranscodingInfo) Equals(other TranscodingInfo) bool {
	return t.jobID == other.jobID &&
		t.progress == other.progress &&
		t.outputURL == other.outputURL &&
		t.error == other.error &&
		t.completedAt == other.completedAt
}

type Credit struct {
	role     string
	name     string
	personID *string
}

func NewCredit(role, name string, personID *string) (*Credit, error) {
	if role == "" {
		return nil, ErrInvalidCreditRole
	}

	if name == "" {
		return nil, ErrInvalidCreditName
	}

	if len(role) > 50 {
		return nil, ErrInvalidCreditRole
	}

	if len(name) > 100 {
		return nil, ErrInvalidCreditName
	}

	if personID != nil {
		if len(*personID) > 100 {
			return nil, ErrInvalidCreditPersonID
		}

		personIDRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !personIDRegex.MatchString(*personID) {
			return nil, ErrInvalidCreditPersonID
		}
	}

	return &Credit{
		role:     strings.TrimSpace(role),
		name:     strings.TrimSpace(name),
		personID: personID,
	}, nil
}

func (c Credit) Role() string {
	return c.role
}

func (c Credit) Name() string {
	return c.name
}

func (c Credit) PersonID() *string {
	return c.personID
}

func (c Credit) Equals(other Credit) bool {
	return c.role == other.role &&
		c.name == other.name &&
		c.personID == other.personID
}

type PublishRule struct {
	publishAt   *time.Time
	unpublishAt *time.Time
	regions     []string
	ageRating   *string
}

func NewPublishRule(publishAt, unpublishAt *time.Time, regions []string, ageRating *string) (*PublishRule, error) {
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
			"G": true, "PG": true, "PG-13": true, "R": true, "NC-17": true,
			"TV-Y": true, "TV-Y7": true, "TV-G": true, "TV-PG": true, "TV-14": true, "TV-MA": true,
		}

		if !validRatings[*ageRating] {
			return nil, ErrInvalidAgeRating
		}
	}

	return &PublishRule{
		publishAt:   publishAt,
		unpublishAt: unpublishAt,
		regions:     validatedRegions,
		ageRating:   ageRating,
	}, nil
}

func (p PublishRule) PublishAt() *time.Time {
	return p.publishAt
}

func (p PublishRule) UnpublishAt() *time.Time {
	return p.unpublishAt
}

func (p PublishRule) Regions() []string {
	return p.regions
}

func (p PublishRule) AgeRating() *string {
	return p.ageRating
}

func (p PublishRule) Equals(other PublishRule) bool {
	return p.publishAt == other.publishAt &&
		p.unpublishAt == other.unpublishAt &&
		p.ageRating == other.ageRating
}

type AssetPage struct {
	Items   []*Asset               `json:"items"`
	LastKey map[string]interface{} `json:"lastKey,omitempty"`
	HasMore bool                   `json:"hasMore"`
	Total   int                    `json:"total"`
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func isValidSlug(slug string) bool {
	if len(slug) < 3 || len(slug) > 50 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return matched
}

var (
	ErrInvalidImageType            = pkgerrors.NewValidationError("invalid image type", nil)
	ErrInvalidS3Bucket             = pkgerrors.NewValidationError("invalid S3 bucket", nil)
	ErrInvalidS3Key                = pkgerrors.NewValidationError("invalid S3 key", nil)
	ErrInvalidS3URL                = pkgerrors.NewValidationError("invalid S3 URL", nil)
	ErrInvalidStreamInfoURL        = pkgerrors.NewValidationError("invalid stream info URL", nil)
	ErrInvalidPublishDates         = pkgerrors.NewValidationError("invalid publish dates", nil)
	ErrTooManyRegions              = pkgerrors.NewValidationError("too many regions", nil)
	ErrInvalidRegion               = pkgerrors.NewValidationError("invalid region", nil)
	ErrInvalidAgeRating            = pkgerrors.NewValidationError("invalid age rating", nil)
	ErrInvalidTranscodingJobID     = pkgerrors.NewValidationError("invalid transcoding job ID", nil)
	ErrInvalidTranscodingProgress  = pkgerrors.NewValidationError("invalid transcoding progress", nil)
	ErrInvalidTranscodingOutputURL = pkgerrors.NewValidationError("invalid transcoding output URL", nil)
	ErrInvalidCreditRole           = pkgerrors.NewValidationError("invalid credit role", nil)
	ErrInvalidCreditName           = pkgerrors.NewValidationError("invalid credit name", nil)
	ErrInvalidCreditPersonID       = pkgerrors.NewValidationError("invalid credit person ID", nil)
)
