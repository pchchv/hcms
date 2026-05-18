package handlers

import "net/url"

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

// cloneValues creates a deep copy of url.Values.
func cloneValues(v url.Values) (out url.Values) {
	for k, vals := range v {
		out[k] = append([]string(nil), vals...)
	}
	return
}
