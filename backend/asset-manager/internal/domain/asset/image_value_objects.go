package asset

import (
	"encoding/json"
	"regexp"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

var (
	ErrInvalidImageID       = pkgerrors.NewValidationError("invalid image id", nil)
	ErrInvalidImageFileName = pkgerrors.NewValidationError("invalid image file name", nil)
	ErrInvalidImageURL      = pkgerrors.NewValidationError("invalid image URL", nil)
)

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

func (id ImageID) Value() string {
	return id.value
}

func (id ImageID) Equals(other ImageID) bool {
	return id.value == other.value
}

func (id ImageID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value)
}

func (id *ImageID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	id.value = s
	return nil
}

type ImageFileName struct {
	value string
}

func NewImageFileName(value string) (*ImageFileName, error) {
	if value == "" {
		return nil, ErrInvalidImageFileName
	}

	if len(value) > 255 {
		return nil, ErrInvalidImageFileName
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, ErrInvalidImageFileName
	}

	return &ImageFileName{value: trimmed}, nil
}

func (f ImageFileName) Value() string {
	return f.value
}

func (f ImageFileName) Equals(other ImageFileName) bool {
	return f.value == other.value
}

func (f ImageFileName) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

func (f *ImageFileName) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	f.value = s
	return nil
}

type ImageURL struct {
	value string
}

func NewImageURL(value string) (*ImageURL, error) {
	if value == "" {
		return nil, ErrInvalidImageURL
	}

	if len(value) > 2048 {
		return nil, ErrInvalidImageURL
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, ErrInvalidImageURL
	}

	return &ImageURL{value: trimmed}, nil
}

func (u ImageURL) Value() string {
	return u.value
}

func (u ImageURL) Equals(other ImageURL) bool {
	return u.value == other.value
}

func (u ImageURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.value)
}

func (u *ImageURL) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	u.value = s
	return nil
}

type ImageContentType struct {
	value string
}

func NewImageContentType(value string) (*ImageContentType, error) {
	if value == "" {
		return &ImageContentType{value: ""}, nil
	}
	if len(value) > 100 {
		return nil, ErrInvalidImageContentType
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return &ImageContentType{value: ""}, nil
	}
	return &ImageContentType{value: trimmed}, nil
}

func (c ImageContentType) Value() string {
	return c.value
}

func (c ImageContentType) Equals(other ImageContentType) bool {
	return c.value == other.value
}

func (c ImageContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c *ImageContentType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	c.value = s
	return nil
}
