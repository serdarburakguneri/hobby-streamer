package valueobjects

import "fmt"

type AssetID struct {
	value string
}

func NewAssetID(value string) (*AssetID, error) {
	if value == "" {
		return nil, fmt.Errorf("asset ID cannot be empty")
	}
	return &AssetID{value: value}, nil
}

func (a AssetID) Value() string {
	return a.value
}

func (a AssetID) String() string {
	return a.value
}
