package models

const (
	StatusNew   LeadStatus = "new"
	StatusSent  LeadStatus = "sent"
	StatusError LeadStatus = "error"
)

// LeadStatus type alias for string.
type LeadStatus = string
