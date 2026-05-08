package middleware

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pchchv/hcms/database"
)

func TestGetSession_NoValue(t *testing.T) {
	s := GetSession(context.Background())
	if s != nil {
		t.Error("expected nil session from empty context")
	}
}

func openAuthDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return db
}
