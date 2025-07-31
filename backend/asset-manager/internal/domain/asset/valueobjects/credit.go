package valueobjects

import (
	"errors"
)

type Credit struct {
	personID  *string
	role      string
	name      string
	order     int
	biography *string
	photoURL  *string
}

func NewCredit(role, name string, order int) (*Credit, error) {
	if role == "" {
		return nil, errors.New("role cannot be empty")
	}

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if len(role) > 100 {
		return nil, errors.New("role too long")
	}

	if len(name) > 200 {
		return nil, errors.New("name too long")
	}

	if order < 0 {
		return nil, errors.New("order cannot be negative")
	}

	return &Credit{
		role:  role,
		name:  name,
		order: order,
	}, nil
}

func (c Credit) PersonID() *string {
	return c.personID
}

func (c Credit) Role() string {
	return c.role
}

func (c Credit) Name() string {
	return c.name
}

func (c Credit) Order() int {
	return c.order
}

func (c Credit) Biography() *string {
	return c.biography
}

func (c Credit) PhotoURL() *string {
	return c.photoURL
}

func (c Credit) SetPersonID(personID string) {
	c.personID = &personID
}

func (c Credit) SetBiography(biography string) {
	c.biography = &biography
}

func (c Credit) SetPhotoURL(photoURL string) {
	c.photoURL = &photoURL
}
