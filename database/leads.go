package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pchchv/hcms/models"
)

// LeadsFilter defines filtering options for listing leads.
type LeadsFilter struct {
	Search string
	Status string
	Limit  int
	Offset int
}

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

// GetLead returns a lead by its ID.
func GetLead(db *sql.DB, id int) (l *models.Lead, err error) {
	row := db.QueryRow(
		`SELECT id, name, phone, email, comment, status, bitrix_response, bitrix_sent_at, created_at
		 FROM leads WHERE id = ?`,
		id,
	)
	l, err = scanLead(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get lead %d: %w", id, err)
	}
	return
}

// ListRecentLeads returns the N most recently created leads.
func ListRecentLeads(db *sql.DB, n int) ([]models.Lead, error) {
	rows, err := db.Query(
		`SELECT id, name, phone, email, comment, status, bitrix_response, bitrix_sent_at, created_at
		 FROM leads ORDER BY created_at DESC LIMIT ?`,
		n,
	)
	if err != nil {
		return nil, fmt.Errorf("list recent leads: %w", err)
	}
	defer rows.Close()

	leads := make([]models.Lead, 0, n)
	for rows.Next() {
		l, err := scanLead(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("scan lead: %w", err)
		}
		leads = append(leads, *l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return leads, nil
}

// DeleteLead removes a lead by ID.
func DeleteLead(db *sql.DB, id int) error {
	if _, err := db.Exec(`DELETE FROM leads WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete lead: %w", err)
	}
	return nil
}

// scanLead scans a lead row from a *sql.Rows or *sql.Row.
func scanLead(scan func(...any) error) (*models.Lead, error) {
	var l models.Lead
	var bitrixSentAt sql.NullTime
	err := scan(
		&l.ID, &l.Name, &l.Phone, &l.Email, &l.Comment,
		&l.Status, &l.BitrixResponse, &bitrixSentAt, &l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if bitrixSentAt.Valid {
		t := bitrixSentAt.Time
		l.BitrixSentAt = &t
	}

	return &l, nil
}
