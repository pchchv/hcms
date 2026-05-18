package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pchchv/hcms/database"
)

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

// HandleNewsDelete handles DELETE /admin/news/{id}.
func (s *Server) HandleNewsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// try from form path
		if err := r.ParseForm(); err == nil {
			idStr = strings.TrimPrefix(r.URL.Path, "/admin/news/")
			if id, err = strconv.Atoi(idStr); err != nil {
				http.Error(w, "Invalid news ID", http.StatusBadRequest)
				return
			}
		}
	}

	deleted, err := database.DeleteNews(s.db, id)
	if err != nil {
		http.Error(w, "Failed to delete news", http.StatusInternalServerError)
		return
	}

	// clean up uploaded image
	if deleted != nil && deleted.Image != "" && s.uploader != nil {
		_ = s.uploader.Delete(deleted.Image)
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	SetFlash(w, Flash{Type: "success", Message: "Новость удалена"})
	http.Redirect(w, r, "/admin/news", http.StatusFound)
}
