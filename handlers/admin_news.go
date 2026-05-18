package handlers

import "net/http"

// HandleNewsPreview handles POST /admin/news/preview.
// Converts Markdown to HTML for HTMX live preview.
func (s *Server) HandleNewsPreview(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("description")
	html, err := renderMarkdown(content)
	if err != nil {
		http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
