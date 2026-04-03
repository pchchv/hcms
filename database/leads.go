package database

import (
	"database/sql"
	"fmt"
	"time"

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

// UpdateLeadStatus updates only the status field of a lead.
func UpdateLeadStatus(db *sql.DB, id int, status string) error {
	if _, err := db.Exec(`UPDATE leads SET status = ? WHERE id = ?`, status, id); err != nil {
		return fmt.Errorf("update lead status: %w", err)
	}
	return nil
}

// UpdateLeadBitrix updates the status, bitrix_response, and bitrix_sent_at fields.
func UpdateLeadBitrix(db *sql.DB, id int, status, response string, sentAt time.Time) error {
	_, err := db.Exec(
		`UPDATE leads SET status = ?, bitrix_response = ?, bitrix_sent_at = ? WHERE id = ?`,
		status, response, sentAt.UTC(), id,
	)
	if err != nil {
		return fmt.Errorf("update lead bitrix: %w", err)
	}
	return nil
}

// DeleteLead removes a lead by ID.
func DeleteLead(db *sql.DB, id int) error {
	if _, err := db.Exec(`DELETE FROM leads WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete lead: %w", err)
	}
	return nil
}
