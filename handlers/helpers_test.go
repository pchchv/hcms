package handlers

import "testing"

func TestNewPagination_Basic(t *testing.T) {
	p := newPagination(100, 2, 10, nil)
	if p.Pages != 10 {
		t.Errorf("expected 10 pages, got %d", p.Pages)
	}
	if p.From != 11 {
		t.Errorf("expected from=11, got %d", p.From)
	}
	if p.To != 20 {
		t.Errorf("expected to=20, got %d", p.To)
	}
}

func TestNewPagination_Empty(t *testing.T) {
	p := newPagination(0, 1, 20, nil)
	if p.From != 0 {
		t.Errorf("expected from=0 for empty result, got %d", p.From)
	}
	if p.Pages != 1 {
		t.Errorf("expected pages=1 for empty result, got %d", p.Pages)
	}
}

func TestNewPagination_ClampPage(t *testing.T) {
	p := newPagination(50, 100, 10, nil) // page 100 > pages 5
	if p.Page != 5 {
		t.Errorf("expected page clamped to 5, got %d", p.Page)
	}
}

func TestNewPagination_DefaultLimit(t *testing.T) {
	p := newPagination(100, 1, 0, nil) // limit=0 should default to 20
	if p.Pages != 5 {
		t.Errorf("expected 5 pages with default limit 20, got %d", p.Pages)
	}
}
