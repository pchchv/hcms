package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pchchv/hcms/models"
)

// CreateNews inserts a new news item and returns its ID.
func CreateNews(db *sql.DB, n *models.News) (int64, error) {
	log.Printf("Creating news: title=%q, date=%v", n.Title, n.Date)

	res, err := db.Exec(
		`INSERT INTO news (date, title, image, announce, description)
		 VALUES (?, ?, ?, ?, ?)`,
		n.Date.UTC(), n.Title, n.Image, n.Announce, n.Description,
	)
	if err != nil {
		return 0, fmt.Errorf("create news: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get news id: %w", err)
	}

	log.Printf("News created with ID=%d", id)
	return id, nil
}
