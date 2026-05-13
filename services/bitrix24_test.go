package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pchchv/hcms/models"
	_ "modernc.org/sqlite"
)

// mockBitrixClient implements BitrixClient for testing.
type mockBitrixClient struct {
	calls  []models.Lead
	errMsg string
}

func (m *mockBitrixClient) SendLead(_ context.Context, lead models.Lead, _ string) error {
	m.calls = append(m.calls, lead)
	if m.errMsg != "" {
		return errors.New(m.errMsg)
	}
	return nil
}

func TestHTTPBitrixClient_SendLead(t *testing.T) {
	var received map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &HTTPBitrixClient{}
	lead := models.Lead{
		Name:      "Test User",
		Phone:     "+79991234567",
		Email:     "test@example.com",
		Comment:   "Test comment",
		CreatedAt: time.Now(),
	}

	if err := client.SendLead(context.Background(), lead, srv.URL); err != nil {
		t.Fatalf("SendLead: %v", err)
	}

	fields, ok := received["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'fields' in payload")
	}
	if fields["NAME"] != "Test User" {
		t.Errorf("expected NAME='Test User', got %v", fields["NAME"])
	}
}

func TestHTTPBitrixClient_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := &HTTPBitrixClient{}
	lead := models.Lead{Name: "Test", Phone: "123", CreatedAt: time.Now()}

	if client.SendLead(context.Background(), lead, srv.URL) == nil {
		t.Error("expected error for 500 response")
	}
}

func setupPoolDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY DEFAULT 1,
			site_name TEXT NOT NULL DEFAULT 'My CMS',
			admin_email TEXT NOT NULL DEFAULT '',
			admin_password TEXT NOT NULL DEFAULT '',
			bitrix24_webhook TEXT NOT NULL DEFAULT '',
			bitrix24_enabled INTEGER NOT NULL DEFAULT 0,
			CHECK(id = 1)
		);
		CREATE TABLE IF NOT EXISTS leads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			comment TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'new',
			bitrix_response TEXT NOT NULL DEFAULT '',
			bitrix_sent_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		INSERT INTO settings (id, bitrix24_webhook, bitrix24_enabled)
		VALUES (1, 'https://example.com/webhook', 1);
	`); err != nil {
		t.Fatalf("setup schema: %v", err)
	}
	return db
}
