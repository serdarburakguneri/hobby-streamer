package valueobjects

const (
	assetIDType   = "asset"
	assetIDLength = 32
)

func NewAssetID(value string) (*ID, error) {
	return NewID(value, assetIDType, assetIDLength)
}

func GenerateAssetID() (*ID, error) {
	id, err := GenerateID(assetIDType, assetIDLength/2)
	if err != nil {
		return nil, err
	}
	return NewAssetID(id)
}
