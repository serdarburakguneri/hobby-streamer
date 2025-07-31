package valueobjects

const (
	imageIDType   = "image"
	imageIDLength = 32
)

func GenerateImageID() (*ID, error) {
	id, err := GenerateID(imageIDType, imageIDLength/2)
	if err != nil {
		return nil, err
	}
	return NewID(id, imageIDType, imageIDLength)
}
