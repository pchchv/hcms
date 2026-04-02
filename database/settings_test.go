package database

import "testing"

func TestGetSettings_NoRow(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	// No settings row yet — should return defaults.
	s, err := Get(d)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if s == nil {
		t.Fatal("Get returned nil")
	}
	if s.SiteName != "My CMS" {
		t.Errorf("expected default SiteName='My CMS', got %q", s.SiteName)
	}
}
