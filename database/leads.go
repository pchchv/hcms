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
