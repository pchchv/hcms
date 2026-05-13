package validators

// Error represents a field-level validation error.
type Error struct {
	Field   string
	Message string
}
