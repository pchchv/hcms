package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID returns a cryptographically random 32-byte hex string (64 characters).
func GenerateID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CreateSession inserts a new session for the given adminID, expiring in 24 hours.
// Returns the new session ID.
func CreateSession(db *sql.DB, adminID int) (string, error) {
	id, err := GenerateID()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(24 * time.Hour)
	_, err = db.Exec(
		`INSERT INTO sessions (id, admin_id, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		id, adminID, now, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return id, nil
}
