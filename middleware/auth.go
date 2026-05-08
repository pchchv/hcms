package middleware

// SessionKey is the context key for storing the authenticated session.
const SessionKey ctxKey = 0

// ctxKey is a private key type for context values to avoid collisions.
type ctxKey int
