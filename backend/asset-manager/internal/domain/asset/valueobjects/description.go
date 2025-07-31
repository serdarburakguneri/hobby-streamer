package valueobjects

type Description struct {
	value string
}

func NewDescription(value string) (*Description, error) {
	validatedString, err := NewValidatedString(value, 2000, "description")
	if err != nil {
		return nil, err
	}
	return &Description{value: validatedString.Value()}, nil
}

func (d Description) Value() string {
	return d.value
}

func (d Description) Equals(other Description) bool {
	return d.value == other.value
}
