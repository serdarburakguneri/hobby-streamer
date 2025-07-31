package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type AssetType struct {
	value string
}

var allowedAssetTypes = constants.AllowedAssetTypes

func NewAssetType(value string) (*AssetType, error) {
	if value == "" {
		return nil, ErrInvalidAssetType
	}
	if _, ok := allowedAssetTypes[value]; !ok {
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

var ErrInvalidAssetType = errors.New("invalid asset type")
