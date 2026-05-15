package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// templatesDir is the relative path to templates from the handlers package.
const templatesDir = "../../internal/templates"

func TestNewRenderer_ReturnsRenderer(t *testing.T) {
	r := NewRenderer(templatesDir)
	if r == nil {
		t.Error("expected non-nil renderer")
	}

	if r.dir != templatesDir {
		t.Errorf("expected dir %q, got %q", templatesDir, r.dir)
	}
}

func TestRenderer_Page_Renders(t *testing.T) {
	r := NewRenderer(templatesDir)
	rr := httptest.NewRecorder()
	data := map[string]any{
		"Page":          "dashboard",
		"Title":         "Dashboard",
		"Session":       nil,
		"CSRFToken":     "",
		"Settings":      nil,
		"NewLeadsCount": 0,
		"Flash":         nil,
		// dashboard-specific
		"TotalLeads":    0,
		"NewLeadsToday": 0,
		"TotalNews":     0,
		"BitrixEnabled": false,
		"RecentLeads":   nil,
	}
	r.Page(rr, "dashboard", data)
	if rr.Code != 200 {
		t.Errorf("expected 200, got %d: %s", rr.Code, rr.Body.String()[:min(200, rr.Body.Len())])
	}

	if !strings.Contains(rr.Body.String(), "<html") {
		t.Errorf("expected HTML response")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
