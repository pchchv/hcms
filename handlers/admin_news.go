package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pchchv/hcms/database"
	"github.com/pchchv/hcms/models"
)

// HandleNews handles GET /admin/news.
func (s *Server) HandleNews(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit
	filter := database.NewsFilter{
		Search: q.Get("search"),
		Limit:  limit,
		Offset: offset,
	}

	result, err := database.ListNews(s.db, filter)
	if err != nil {
		log.Printf("Error listing news: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("News list: total=%d, returned=%d", result.Total, len(result.News))
	pagination := newPagination(result.Total, page, limit, q)
	data := baseData(r, w, s.db, "news", "Новости")
	data["NewsList"] = result.News
	data["Pagination"] = pagination
	data["Filter"] = filter
	s.renderer.Page(w, "news", data)
}

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

// HandleNewsNew handles GET /admin/news/new.
func (s *Server) HandleNewsNew(w http.ResponseWriter, r *http.Request) {
	data := baseData(r, w, s.db, "news_form", "Добавить новость")
	data["News"] = &models.News{Date: time.Now()}
	data["IsNew"] = true
	s.renderer.Page(w, "news_form", data)
}

// HandleNewsEdit handles GET /admin/news/{id}/edit.
func (s *Server) HandleNewsEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	n, err := database.GetNews(s.db, id)
	if err != nil || n == nil {
		http.NotFound(w, r)
		return
	}

	data := baseData(r, w, s.db, "news_form", "Редактировать новость")
	data["News"] = n
	data["IsNew"] = false
	s.renderer.Page(w, "news_form", data)
}
