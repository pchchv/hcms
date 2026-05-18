package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pchchv/hcms/database"
	"github.com/pchchv/hcms/middleware"
)

// PaginationData holds computed pagination metadata.
type PaginationData struct {
	Page  int
	Pages int
	Total int
	From  int
	To    int
	query url.Values
}

// newPagination constructs a PaginationData from total count,
// current page, limit per page, and the current URL query values.
func newPagination(total, page, limit int, q url.Values) PaginationData {
	if limit <= 0 {
		limit = 20
	}

	pages := total / limit
	if total%limit != 0 {
		pages++
	}

	if pages < 1 {
		pages = 1
	}

	if page < 1 {
		page = 1
	}

	if page > pages {
		page = pages
	}

	to := page * limit
	if to > total {
		to = total
	}

	from := (page-1)*limit + 1
	if total == 0 {
		from = 0
	}

	return PaginationData{
		Page:  page,
		Pages: pages,
		Total: total,
		From:  from,
		To:    to,
		query: q,
	}
}

// QueryWithPage returns the current URL query string with the page parameter replaced.
func (p PaginationData) QueryWithPage(pg int) string {
	q := cloneValues(p.query)
	q.Set("page", strconv.Itoa(pg))
	return "?" + q.Encode()
}

// PageNumbers returns a slice of page numbers to display,
// with -1 representing ellipsis.
func (p PaginationData) PageNumbers() []int {
	if p.Pages <= 7 {
		pages := make([]int, p.Pages)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	result := append([]int{}, 1)
	if p.Page > 3 {
		result = append(result, -1) // ellipsis
	}

	start := p.Page - 1
	if start < 2 {
		start = 2
	}

	end := p.Page + 1
	if end > p.Pages-1 {
		end = p.Pages - 1
	}

	for i := start; i <= end; i++ {
		result = append(result, i)
	}

	if p.Page < p.Pages-2 {
		result = append(result, -1) // ellipsis
	}

	return append(result, p.Pages)
}

// cloneValues creates a deep copy of url.Values.
func cloneValues(v url.Values) url.Values {
	out := make(url.Values)
	for k, vals := range v {
		out[k] = append([]string(nil), vals...)
	}
	return out
}

// baseData builds the common data map needed by all admin pages.
// It queries settings, new lead count, CSRF token, and flash messages.
func baseData(r *http.Request, w http.ResponseWriter, db *sql.DB, page, title string) map[string]any {
	var csrfToken string
	session := middleware.GetSession(r.Context())
	if session != nil {
		csrfToken = middleware.GenerateToken(session.ID)
	}

	settings, _ := database.Get(db)
	newLeadsCount, _ := database.CountByStatus(db, "new")
	flash := GetFlash(r, w)
	data := map[string]any{
		"Page":          page,
		"Title":         title,
		"Session":       session,
		"CSRFToken":     csrfToken,
		"Settings":      settings,
		"NewLeadsCount": newLeadsCount,
		"Flash":         flash,
	}
	return data
}
