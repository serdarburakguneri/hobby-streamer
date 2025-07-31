package valueobjects

const (
	videoIDType   = "video"
	videoIDLength = 32
)

func NewVideoID(value string) (*ID, error) {
	return NewID(value, videoIDType, videoIDLength)
}

func GenerateVideoID() (*ID, error) {
	id, err := GenerateID(videoIDType, videoIDLength/2)
	if err != nil {
		return nil, err
	}
	return NewVideoID(id)
}
