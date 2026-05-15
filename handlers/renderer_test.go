package handlers

import "testing"

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
