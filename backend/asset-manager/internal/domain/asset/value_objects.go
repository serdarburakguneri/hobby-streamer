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

	if len(value) < 3 || len(value) > 50 {
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

	if len(value) > 200 {
		return nil, ErrInvalidTitle
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, ErrInvalidTitle
	}

	return &Title{value: trimmed}, nil
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

	trimmed := strings.TrimSpace(value)
	return &Description{value: trimmed}, nil
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
		constants.AssetTypeMovie:        true,
		constants.AssetTypeTVShow:       true,
		constants.AssetTypeSeries:       true,
		constants.AssetTypeSeason:       true,
		constants.AssetTypeEpisode:      true,
		constants.AssetTypeDocumentary:  true,
		constants.AssetTypeShort:        true,
		constants.AssetTypeTrailer:      true,
		constants.AssetTypeBonus:        true,
		constants.AssetTypeBehindScenes: true,
		constants.AssetTypeInterview:    true,
		constants.AssetTypeMusicVideo:   true,
		constants.AssetTypePodcast:      true,
		constants.AssetTypeLive:         true,
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

	genreRegex := regexp.MustCompile(`^[a-zA-Z0-9\s&-]+$`)
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

func NewGenres(genreStrings []string) (*Genres, error) {
	if len(genreStrings) > 10 {
		return nil, ErrTooManyGenres
	}

	genres := make([]Genre, 0, len(genreStrings))
	seen := make(map[string]bool)

	for _, genreStr := range genreStrings {
		if seen[genreStr] {
			continue
		}

		genre, err := NewGenre(genreStr)
		if err != nil {
			return nil, err
		}

		genres = append(genres, *genre)
		seen[genreStr] = true
	}

	return &Genres{values: genres}, nil
}

func (g Genres) Values() []Genre {
	return g.values
}

func (g Genres) Contains(genre Genre) bool {
	for _, g := range g.values {
		if g.Equals(genre) {
			return true
		}
	}
	return false
}

func (g Genres) Add(genre Genre) *Genres {
	if g.Contains(genre) {
		return &Genres{values: g.values}
	}

	newGenres := make([]Genre, len(g.values)+1)
	copy(newGenres, g.values)
	newGenres[len(g.values)] = genre

	return &Genres{values: newGenres}
}

func (g Genres) Remove(genre Genre) *Genres {
	newGenres := make([]Genre, 0, len(g.values))
	for _, g := range g.values {
		if !g.Equals(genre) {
			newGenres = append(newGenres, g)
		}
	}
	return &Genres{values: newGenres}
}

type Tag struct {
	value string
}

func NewTag(value string) (*Tag, error) {
	if value == "" {
		return nil, ErrInvalidTag
	}

	if len(value) > 30 {
		return nil, ErrInvalidTag
	}

	tagRegex := regexp.MustCompile(`^[a-zA-Z0-9\s&-]+$`)
	if !tagRegex.MatchString(value) {
		return nil, ErrInvalidTag
	}

	return &Tag{value: strings.TrimSpace(value)}, nil
}

func (t Tag) Value() string {
	return t.value
}

func (t Tag) Equals(other Tag) bool {
	return t.value == other.value
}

type Tags struct {
	values []Tag
}

func NewTags(tagStrings []string) (*Tags, error) {
	if len(tagStrings) > 20 {
		return nil, ErrTooManyTags
	}

	tags := make([]Tag, 0, len(tagStrings))
	seen := make(map[string]bool)

	for _, tagStr := range tagStrings {
		if seen[tagStr] {
			continue
		}

		tag, err := NewTag(tagStr)
		if err != nil {
			return nil, err
		}

		tags = append(tags, *tag)
		seen[tagStr] = true
	}

	return &Tags{values: tags}, nil
}

func (t Tags) Values() []Tag {
	return t.values
}

func (t Tags) Contains(tag Tag) bool {
	for _, t := range t.values {
		if t.Equals(tag) {
			return true
		}
	}
	return false
}

func (t Tags) Add(tag Tag) *Tags {
	if t.Contains(tag) {
		return &Tags{values: t.values}
	}

	newTags := make([]Tag, len(t.values)+1)
	copy(newTags, t.values)
	newTags[len(t.values)] = tag

	return &Tags{values: newTags}
}

func (t Tags) Remove(tag Tag) *Tags {
	newTags := make([]Tag, 0, len(t.values))
	for _, t := range t.values {
		if !t.Equals(tag) {
			newTags = append(newTags, t)
		}
	}
	return &Tags{values: newTags}
}

type OwnerID struct {
	value string
}

func NewOwnerID(value string) (*OwnerID, error) {
	if value == "" {
		return nil, ErrInvalidOwnerID
	}

	if len(value) > 100 {
		return nil, ErrInvalidOwnerID
	}

	ownerIDRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !ownerIDRegex.MatchString(value) {
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
	ErrInvalidSlug        = errors.New("invalid slug")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidAssetType   = errors.New("invalid asset type")
	ErrInvalidGenre       = errors.New("invalid genre")
	ErrTooManyGenres      = errors.New("too many genres")
	ErrInvalidTag         = errors.New("invalid tag")
	ErrTooManyTags        = errors.New("too many tags")
	ErrInvalidOwnerID     = errors.New("invalid owner ID")
)
