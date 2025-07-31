package valueobjects

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(idType string, length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
