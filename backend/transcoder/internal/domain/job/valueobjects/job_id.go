package valueobjects

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type JobID struct {
	value string
}

func NewJobID(value string) (*JobID, error) {
	if value == "" {
		return nil, fmt.Errorf("job ID cannot be empty")
	}
	return &JobID{value: value}, nil
}

func (j JobID) Value() string {
	return j.value
}

func (j JobID) String() string {
	return j.value
}

func GenerateJobID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
