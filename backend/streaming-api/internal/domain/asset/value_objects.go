package asset

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type AssetID struct {
	value string
}

func NewAssetID(value string) (*AssetID, error) {
	if value == "" {
		return nil, ErrInvalidAssetID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidAssetID
	}

	return &AssetID{value: value}, nil
}

func (id AssetID) Value() string {
	return id.value
}

func (id AssetID) Equals(other AssetID) bool {
	return id.value == other.value
}

type Slug struct {
	value string
}

func NewSlug(value string) (*Slug, error) {
	if value == "" {
		return nil, ErrInvalidSlug
	}

	if len(value) < 3 || len(value) > 100 {
		return nil, ErrInvalidSlug
	}

	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(value) {
		return nil, ErrInvalidSlug
	}

	return &Slug{value: strings.ToLower(value)}, nil
}

func (s Slug) Value() string {
	return s.value
}

func (s Slug) Equals(other Slug) bool {
	return s.value == other.value
}

type Title struct {
	value string
}

func NewTitle(value string) (*Title, error) {
	if value == "" {
		return nil, ErrInvalidTitle
	}

	if len(value) > 255 {
		return nil, ErrInvalidTitle
	}

	return &Title{value: strings.TrimSpace(value)}, nil
}

func (t Title) Value() string {
	return t.value
}

func (t Title) Equals(other Title) bool {
	return t.value == other.value
}

type Description struct {
	value string
}

func NewDescription(value string) (*Description, error) {
	if len(value) > 1000 {
		return nil, ErrInvalidDescription
	}

	return &Description{value: strings.TrimSpace(value)}, nil
}

func (d Description) Value() string {
	return d.value
}

func (d Description) Equals(other Description) bool {
	return d.value == other.value
}

type AssetType struct {
	value string
}

func NewAssetType(value string) (*AssetType, error) {
	if value == "" {
		return nil, ErrInvalidAssetType
	}

	validTypes := map[string]bool{
		"movie":       true,
		"series":      true,
		"episode":     true,
		"documentary": true,
		"short":       true,
		"trailer":     true,
		"music":       true,
		"podcast":     true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidAssetType
	}

	return &AssetType{value: value}, nil
}

func (t AssetType) Value() string {
	return t.value
}

func (t AssetType) Equals(other AssetType) bool {
	return t.value == other.value
}

type Genre struct {
	value string
}

func NewGenre(value string) (*Genre, error) {
	if value == "" {
		return nil, ErrInvalidGenre
	}

	if len(value) > 50 {
		return nil, ErrInvalidGenre
	}

	genreRegex := regexp.MustCompile(`^[a-zA-Z0-9\s-]+$`)
	if !genreRegex.MatchString(value) {
		return nil, ErrInvalidGenre
	}

	return &Genre{value: strings.TrimSpace(value)}, nil
}

func (g Genre) Value() string {
	return g.value
}

func (g Genre) Equals(other Genre) bool {
	return g.value == other.value
}

type Genres struct {
	values []Genre
}

func NewGenres(genres []string) (*Genres, error) {
	genreObjects := make([]Genre, 0, len(genres))
	for _, genreStr := range genres {
		genre, err := NewGenre(genreStr)
		if err != nil {
			return nil, err
		}
		genreObjects = append(genreObjects, *genre)
	}
	return &Genres{values: genreObjects}, nil
}

func (g Genres) Values() []Genre {
	return g.values
}

func (g Genres) Contains(genre Genre) bool {
	for _, existingGenre := range g.values {
		if existingGenre.Equals(genre) {
			return true
		}
	}
	return false
}

type Tags struct {
	values []string
}

func NewTags(tags []string) (*Tags, error) {
	if len(tags) > 20 {
		return nil, ErrTooManyTags
	}

	validatedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag) > 30 {
			return nil, ErrInvalidTag
		}

		tagRegex := regexp.MustCompile(`^[a-zA-Z0-9\s-]+$`)
		if !tagRegex.MatchString(tag) {
			return nil, ErrInvalidTag
		}

		validatedTags = append(validatedTags, strings.TrimSpace(tag))
	}

	return &Tags{values: validatedTags}, nil
}

func (t Tags) Values() []string {
	return t.values
}

func (t Tags) Contains(tag string) bool {
	for _, existingTag := range t.values {
		if existingTag == tag {
			return true
		}
	}
	return false
}

type Status struct {
	value string
}

func NewStatus(value string) (*Status, error) {
	if value == "" {
		return nil, ErrInvalidStatus
	}

	validStatuses := map[string]bool{
		constants.AssetStatusDraft:     true,
		"processing":                   true,
		constants.VideoStatusReady:     true,
		constants.AssetStatusPublished: true,
		"archived":                     true,
		"deleted":                      true,
	}

	if !validStatuses[value] {
		return nil, ErrInvalidStatus
	}

	return &Status{value: value}, nil
}

func (s Status) Value() string {
	return s.value
}

func (s Status) Equals(other Status) bool {
	return s.value == other.value
}

type OwnerID struct {
	value string
}

func NewOwnerID(value string) (*OwnerID, error) {
	if value == "" {
		return nil, ErrInvalidOwnerID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidOwnerID
	}

	return &OwnerID{value: value}, nil
}

func (o OwnerID) Value() string {
	return o.value
}

func (o OwnerID) Equals(other OwnerID) bool {
	return o.value == other.value
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
	ErrInvalidAssetID     = errors.New("invalid asset ID")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidGenre       = errors.New("invalid genre")
	ErrInvalidTag         = errors.New("invalid tag")
	ErrTooManyTags        = errors.New("too many tags")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidOwnerID     = errors.New("invalid owner ID")
)

var validVideoStatuses = map[string]bool{
	constants.VideoStatusReady: true,
}
