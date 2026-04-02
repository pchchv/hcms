package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

// UpdateNews updates all fields of a news item and sets updated_at to now.
func UpdateNews(db *sql.DB, n *models.News) error {
	_, err := db.Exec(
		`UPDATE news SET date = ?, title = ?, image = ?, announce = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		n.Date.UTC(), n.Title, n.Image, n.Announce, n.Description, time.Now().UTC(), n.ID,
	)
	if err != nil {
		return fmt.Errorf("update news: %w", err)
	}
	return nil
}

// CountNews returns the total number of news items.
func CountNews(db *sql.DB) (count int, err error) {
	if err = db.QueryRow(`SELECT COUNT(*) FROM news`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count news: %w", err)
	}
	return
}

// scanNews scans a news row.
func scanNews(scan func(...any) error) (*models.News, error) {
	var n models.News
	if err := scan(&n.ID, &n.Date, &n.Title, &n.Image, &n.Announce, &n.Description,
		&n.CreatedAt, &n.UpdatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}
