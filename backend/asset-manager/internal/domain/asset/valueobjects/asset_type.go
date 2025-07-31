package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type AssetType struct {
	value string
}

func NewAssetType(value string) (*AssetType, error) {
	if value == "" {
		return nil, errors.New("asset type cannot be empty")
	}

	if !constants.IsValidAssetType(value) {
		return nil, errors.New("invalid asset type")
	}

	return &AssetType{value: value}, nil
}

func (at AssetType) Value() string {
	return at.value
}

func (at AssetType) Equals(other AssetType) bool {
	return at.value == other.value
}
