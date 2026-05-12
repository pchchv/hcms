package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// csrfSecret is generated once at startup.
var csrfSecret []byte

func init() {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic("csrf: failed to generate secret: " + err.Error())
	}

	csrfSecret = secret
}

// GenerateToken returns an HMAC-SHA256 of the sessionID using the CSRF secret.
func GenerateToken(sessionID string) string {
	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(sessionID))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks that the provided token matches the expected HMAC for sessionID.
// Uses constant-time comparison to prevent timing attacks.
func Verify(sessionID, token string) bool {
	expected := GenerateToken(sessionID)
	eBytes, err := hex.DecodeString(expected)
	if err != nil {
		return false
	}

	tBytes, err := hex.DecodeString(token)
	if err != nil {
		return false
	}

	return hmac.Equal(eBytes, tBytes)
}
