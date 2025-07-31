package valueobjects

import (
	"errors"
	"regexp"
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

var ErrInvalidAssetID = errors.New("invalid asset ID")
