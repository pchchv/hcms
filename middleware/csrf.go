package middleware

import "crypto/rand"

// csrfSecret is generated once at startup.
var csrfSecret []byte

func init() {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic("csrf: failed to generate secret: " + err.Error())
	}

	csrfSecret = secret
}
