package bucket

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type BucketID struct {
	value string
}

func NewBucketID(value string) (*BucketID, error) {
	if value == "" {
		return nil, ErrInvalidBucketID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidBucketID
	}

	return &BucketID{value: value}, nil
}

func (id BucketID) Value() string {
	return id.value
}

func (id BucketID) Equals(other BucketID) bool {
	return id.value == other.value
}

type BucketKey struct {
	value string
}

func NewBucketKey(value string) (*BucketKey, error) {
	if value == "" {
		return nil, ErrInvalidBucketKey
	}

	if len(value) < 3 || len(value) > 50 {
		return nil, ErrInvalidBucketKey
	}

	keyRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !keyRegex.MatchString(value) {
		return nil, ErrInvalidBucketKey
	}

	return &BucketKey{value: strings.ToLower(value)}, nil
}

func (k BucketKey) Value() string {
	return k.value
}

func (k BucketKey) Equals(other BucketKey) bool {
	return k.value == other.value
}

type BucketName struct {
	value string
}

func NewBucketName(value string) (*BucketName, error) {
	if value == "" {
		return nil, ErrInvalidBucketName
	}

	if len(value) > 100 {
		return nil, ErrInvalidBucketName
	}

	return &BucketName{value: strings.TrimSpace(value)}, nil
}

func (n BucketName) Value() string {
	return n.value
}

func (n BucketName) Equals(other BucketName) bool {
	return n.value == other.value
}

type BucketDescription struct {
	value string
}

func NewBucketDescription(value string) (*BucketDescription, error) {
	if len(value) > 500 {
		return nil, ErrInvalidBucketDescription
	}

	return &BucketDescription{value: strings.TrimSpace(value)}, nil
}

func (d BucketDescription) Value() string {
	return d.value
}

func (d BucketDescription) Equals(other BucketDescription) bool {
	return d.value == other.value
}

type BucketType struct {
	value string
}

func NewBucketType(value string) (*BucketType, error) {
	if value == "" {
		return nil, ErrInvalidBucketType
	}

	validTypes := map[string]bool{
		"collection": true,
		"playlist":   true,
		"category":   true,
		"featured":   true,
		"trending":   true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidBucketType
	}

	return &BucketType{value: value}, nil
}

func (t BucketType) Value() string {
	return t.value
}

func (t BucketType) Equals(other BucketType) bool {
	return t.value == other.value
}

type BucketStatus struct {
	value string
}

func NewBucketStatus(value string) (*BucketStatus, error) {
	if value == "" {
		return nil, ErrInvalidBucketStatus
	}

	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"draft":    true,
		"archived": true,
	}

	if !validStatuses[value] {
		return nil, ErrInvalidBucketStatus
	}

	return &BucketStatus{value: value}, nil
}

func (s BucketStatus) Value() string {
	return s.value
}

func (s BucketStatus) Equals(other BucketStatus) bool {
	return s.value == other.value
}

type AssetIDs struct {
	values []string
}

func NewAssetIDs(assetIDs []string) (*AssetIDs, error) {
	if len(assetIDs) > 1000 {
		return nil, ErrTooManyAssets
	}

	validatedIDs := make([]string, 0, len(assetIDs))
	for _, id := range assetIDs {
		if len(id) > 100 {
			return nil, ErrInvalidAssetID
		}

		idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !idRegex.MatchString(id) {
			return nil, ErrInvalidAssetID
		}

		validatedIDs = append(validatedIDs, id)
	}

	return &AssetIDs{values: validatedIDs}, nil
}

func (a AssetIDs) Values() []string {
	return a.values
}

func (a AssetIDs) Contains(assetID string) bool {
	for _, id := range a.values {
		if id == assetID {
			return true
		}
	}
	return false
}

func (a AssetIDs) Add(assetID string) *AssetIDs {
	if !a.Contains(assetID) {
		newIDs := make([]string, len(a.values)+1)
		copy(newIDs, a.values)
		newIDs[len(a.values)] = assetID
		return &AssetIDs{values: newIDs}
	}
	return &AssetIDs{values: a.values}
}

func (a AssetIDs) Remove(assetID string) *AssetIDs {
	newIDs := make([]string, 0, len(a.values))
	for _, id := range a.values {
		if id != assetID {
			newIDs = append(newIDs, id)
		}
	}
	return &AssetIDs{values: newIDs}
}

type CreatedAt struct {
	value time.Time
}

func NewCreatedAt(value time.Time) *CreatedAt {
	return &CreatedAt{value: value}
}

func (c CreatedAt) Value() time.Time {
	return c.value
}

type UpdatedAt struct {
	value time.Time
}

func NewUpdatedAt(value time.Time) *UpdatedAt {
	return &UpdatedAt{value: value}
}

func (u UpdatedAt) Value() time.Time {
	return u.value
}

var (
	ErrInvalidBucketID          = errors.New("invalid bucket ID")
	ErrInvalidBucketKey         = errors.New("invalid bucket key")
	ErrInvalidBucketName        = errors.New("invalid bucket name")
	ErrInvalidBucketDescription = errors.New("invalid bucket description")
	ErrInvalidBucketType        = errors.New("invalid bucket type")
	ErrInvalidBucketStatus      = errors.New("invalid bucket status")
	ErrInvalidAssetID           = errors.New("invalid asset ID")
	ErrTooManyAssets            = errors.New("too many assets")
)
