package valueobjects

import (
	"time"
)

type CreatedAt struct {
	value time.Time
}

func NewCreatedAt(value time.Time) CreatedAt {
	return CreatedAt{value: value}
}

func (c CreatedAt) Value() time.Time {
	return c.value
}

type UpdatedAt struct {
	value time.Time
}

func NewUpdatedAt(value time.Time) UpdatedAt {
	return UpdatedAt{value: value}
}

func (u UpdatedAt) Value() time.Time {
	return u.value
}
