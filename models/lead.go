package models

import "time"

const (
	StatusNew   LeadStatus = "new"
	StatusSent  LeadStatus = "sent"
	StatusError LeadStatus = "error"
)

// LeadStatus type alias for string.
type LeadStatus = string

// Lead represents a contact form submission.
type Lead struct {
	ID             int
	Name           string
	Phone          string
	Email          string
	Comment        string
	BitrixResponse string
	BitrixSentAt   *time.Time
	CreatedAt      time.Time
	Status         LeadStatus
}
