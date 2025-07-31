package valueobjects

import (
	"errors"
	"regexp"
)

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

var (
	ErrInvalidAssetID = errors.New("invalid asset ID")
	ErrTooManyAssets  = errors.New("too many assets")
)
