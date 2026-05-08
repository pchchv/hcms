package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	if ip := realIP(req); ip != "1.2.3.4" {
		t.Errorf("expected '1.2.3.4', got %q", ip)
	}
}

func TestRealIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "9.9.9.9")
	if ip := realIP(req); ip != "9.9.9.9" {
		t.Errorf("expected '9.9.9.9', got %q", ip)
	}
}

func TestRealIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:4321"
	if ip := realIP(req); ip != "127.0.0.1" {
		t.Errorf("expected '127.0.0.1', got %q", ip)
	}
}
