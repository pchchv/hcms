package handlers

// Flash represents a flash notification message.
type Flash struct {
	Type    string `json:"type"` // success | error | warning
	Message string `json:"message"`
}
