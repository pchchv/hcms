package handlers

import (
	"database/sql"
	"testing"

	"github.com/pchchv/hcms/config"
	"github.com/pchchv/hcms/database"
)

func setupTestServer(t *testing.T) (*Server, *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg := &config.Config{
		Port:       8080,
		DBPath:     ":memory:",
		UploadPath: t.TempDir(),
	}
	s := &Server{
		db:  db,
		cfg: cfg,
	}

	return s, db
}
