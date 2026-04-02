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

// GetNews returns a news item by its ID.
func GetNews(db *sql.DB, id int) (*models.News, error) {
	row := db.QueryRow(
		`SELECT id, date, title, image, announce, description, created_at, updated_at
		 FROM news WHERE id = ?`,
		id,
	)
	n, err := scanNews(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get news %d: %w", id, err)
	}

	return n, nil
}

// DeleteNews removes a news item by ID and returns the deleted item (for image cleanup).
func DeleteNews(db *sql.DB, id int) (*models.News, error) {
	// Fetch first so we can return info for cleanup.
	n, err := GetNews(db, id)
	if err != nil {
		return nil, err
	} else if n == nil {
		return nil, nil
	}

	if _, err := db.Exec(`DELETE FROM news WHERE id = ?`, id); err != nil {
		return nil, fmt.Errorf("delete news: %w", err)
	}

	return n, nil
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
