package services

import (
	"context"
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
