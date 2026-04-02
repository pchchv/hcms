package database

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateID returns a cryptographically random 32-byte hex string (64 characters).
func GenerateID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}
	return hex.EncodeToString(b), nil
}
