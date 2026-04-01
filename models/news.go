package models

import "time"

// News represents a news article.
type News struct {
	ID          int
	Title       string
	Image       string
	Description string
	Announce    string
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
