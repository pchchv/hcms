package database

import "testing"

func TestGenerateID_Unique(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 50; i++ {
		id, err := GenerateID()
		if err != nil {
			t.Fatalf("GenerateID: %v", err)
		}
		if len(id) != 64 {
			t.Errorf("expected 64-char hex ID, got len=%d", len(id))
		}
		if ids[id] {
			t.Errorf("duplicate session ID generated")
		}
		ids[id] = true
	}
}
