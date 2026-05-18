package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFlash_NoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	flash := GetFlash(req, rr)
	if flash != nil {
		t.Error("expected nil flash when no cookie")
	}
}

func TestGetFlash_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: flashCookieName, Value: "not-valid-json"})
	rr := httptest.NewRecorder()
	flash := GetFlash(req, rr)
	if flash != nil {
		t.Error("expected nil flash for invalid JSON")
	}
}
