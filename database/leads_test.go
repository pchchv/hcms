package database

import (
	"testing"

	"github.com/pchchv/hcms/models"
)

func TestCountByStatus(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, _ = CreateLead(d, &models.Lead{Name: "A", Phone: "1", Status: models.StatusNew})
	_, _ = CreateLead(d, &models.Lead{Name: "B", Phone: "2", Status: models.StatusNew})
	id, _ := CreateLead(d, &models.Lead{Name: "C", Phone: "3", Status: models.StatusNew})
	_ = UpdateLeadStatus(d, int(id), models.StatusSent)
	newCount, err := CountByStatus(d, models.StatusNew)
	if err != nil {
		t.Fatalf("CountByStatus new: %v", err)
	}
	if newCount != 2 {
		t.Errorf("expected 2 new, got %d", newCount)
	}

	sentCount, err := CountByStatus(d, models.StatusSent)
	if err != nil {
		t.Fatalf("CountByStatus sent: %v", err)
	}
	if sentCount != 1 {
		t.Errorf("expected 1 sent, got %d", sentCount)
	}
}
