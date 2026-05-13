package validators

// Error represents a field-level validation error.
type Error struct {
	Field   string
	Message string
}

// LeadInput holds the raw input for a lead form submission.
type LeadInput struct {
	Name    string
	Phone   string
	Email   string
	Comment string
}
