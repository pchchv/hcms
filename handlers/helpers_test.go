package handlers

import (
	"net/url"
	"strings"
	"testing"
)

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

func TestPageNumbers_FewPages(t *testing.T) {
	p := PaginationData{Page: 1, Pages: 5}
	nums := p.PageNumbers()
	if len(nums) != 5 {
		t.Errorf("expected 5 page numbers, got %d", len(nums))
	}

	for i, n := range nums {
		if n != i+1 {
			t.Errorf("expected page %d, got %d", i+1, n)
		}
	}
}

func TestPageNumbers_ManyPages_Start(t *testing.T) {
	p := PaginationData{Page: 1, Pages: 20}
	nums := p.PageNumbers()
	// should start with 1 and end with 20,
	// contain -1 for ellipsis
	if nums[0] != 1 {
		t.Error("first page should be 1")
	}

	if nums[len(nums)-1] != 20 {
		t.Error("last page should be 20")
	}

	var hasEllipsis bool
	for _, n := range nums {
		if n == -1 {
			hasEllipsis = true
		}
	}

	if !hasEllipsis {
		t.Error("expected ellipsis (-1) in page numbers")
	}
}

func TestPageNumbers_ManyPages_Middle(t *testing.T) {
	p := PaginationData{Page: 10, Pages: 20}
	nums := p.PageNumbers()
	if nums[0] != 1 {
		t.Error("first should always be 1")
	}

	if nums[len(nums)-1] != 20 {
		t.Error("last should always be 20")
	}

	// should have two ellipses
	var ellipses int
	for _, n := range nums {
		if n == -1 {
			ellipses++
		}
	}

	if ellipses != 2 {
		t.Errorf("expected 2 ellipses for middle page, got %d", ellipses)
	}
}

func TestPageNumbers_ManyPages_End(t *testing.T) {
	p := PaginationData{Page: 20, Pages: 20}
	nums := p.PageNumbers()
	if nums[0] != 1 {
		t.Error("first page should be 1")
	}
	if nums[len(nums)-1] != 20 {
		t.Error("last page should be 20")
	}
}

func TestQueryWithPage(t *testing.T) {
	q := url.Values{}
	q.Set("status", "new")
	q.Set("page", "1")
	p := PaginationData{Page: 1, Pages: 5, query: q}
	result := p.QueryWithPage(3)
	if !strings.Contains(result, "page=3") {
		t.Errorf("expected page=3 in query, got %q", result)
	}

	if !strings.Contains(result, "status=new") {
		t.Errorf("expected status=new preserved in query, got %q", result)
	}
}

func TestCloneValues_Independence(t *testing.T) {
	orig := url.Values{"key": {"value1", "value2"}}
	clone := cloneValues(orig)
	clone["key"][0] = "modified"
	if orig["key"][0] == "modified" {
		t.Error("clone should not share slices with original")
	}
}
