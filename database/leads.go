package database

import (
	"database/sql"
	"fmt"

	"github.com/pchchv/hcms/models"
)

// CreateLead inserts a new lead and returns its ID.
func CreateLead(db *sql.DB, lead *models.Lead) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO leads (name, phone, email, comment, status, bitrix_response)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		lead.Name, lead.Phone, lead.Email, lead.Comment, lead.Status, lead.BitrixResponse,
	)
	if err != nil {
		return 0, fmt.Errorf("create lead: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get lead id: %w", err)
	}

	return id, nil
}

// CountByStatus returns the count of leads with the given status.
func CountByStatus(db *sql.DB, status string) (count int, err error) {
	if err = db.QueryRow(`SELECT COUNT(*) FROM leads WHERE status = ?`, status).Scan(&count); err != nil {
		return 0, fmt.Errorf("count by status: %w", err)
	}
	return
}

// CountCreatedToday returns the count of leads created today (UTC).
func CountCreatedToday(db *sql.DB) (count int, err error) {
	err = db.QueryRow(
		`SELECT COUNT(*) FROM leads WHERE date(created_at) = date('now')`,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count today: %w", err)
	}
	return
}
