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

// HandleNewsCreate handles POST /admin/news.
func (s *Server) HandleNewsCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
	}

	n, formErrors := parseNewsForm(r)
	if len(formErrors) > 0 {
		data := baseData(r, w, s.db, "news_form", "Добавить новость")
		data["News"] = n
		data["IsNew"] = true
		data["Errors"] = formErrors
		s.renderer.Page(w, "news_form", data)
		return
	}

	// рandle image upload
	if s.uploader != nil {
		file, fileHeader, err := r.FormFile("image")
		if err == nil && file != nil && fileHeader != nil {
			defer file.Close()

			path, err := s.uploader.Save(fileHeader)
			if err != nil {
				log.Printf("Warning: image upload failed: %v", err)
			} else {
				n.Image = path
				log.Printf("Image uploaded: %s", path)
			}
		}
	}

	if _, err := database.CreateNews(s.db, n); err != nil {
		log.Printf("Error creating news: %v", err)
		data := baseData(r, w, s.db, "news_form", "Добавить новость")
		data["News"] = n
		data["IsNew"] = true
		data["Errors"] = map[string]string{"general": "Не удалось сохранить новость: " + err.Error()}
		s.renderer.Page(w, "news_form", data)
		return
	}

	SetFlash(w, Flash{Type: "success", Message: "Новость успешно добавлена"})
	http.Redirect(w, r, "/admin/news", http.StatusFound)
}

// HandleNewsUpdate handles POST /admin/news/{id} with _method=PUT.
func (s *Server) HandleNewsUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	existing, err := database.GetNews(s.db, id)
	if err != nil || existing == nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		r.ParseForm()
	}

	n, formErrors := parseNewsForm(r)
	n.ID = id
	if len(formErrors) > 0 {
		data := baseData(r, w, s.db, "news_form", "Редактировать новость")
		data["News"] = n
		data["IsNew"] = false
		data["Errors"] = formErrors
		s.renderer.Page(w, "news_form", data)
		return
	}

	// handle image upload
	if s.uploader != nil {
		if _, fileHeader, err := r.FormFile("image"); err == nil && fileHeader != nil {
			path, err := s.uploader.Save(fileHeader)
			if err == nil {
				// delete old image
				if existing.Image != "" {
					_ = s.uploader.Delete(existing.Image)
				}
				n.Image = path
			}
		} else {
			// keep existing image
			n.Image = existing.Image
		}
	} else {
		n.Image = existing.Image
	}

	if err := database.UpdateNews(s.db, n); err != nil {
		data := baseData(r, w, s.db, "news_form", "Редактировать новость")
		data["News"] = n
		data["IsNew"] = false
		data["Errors"] = map[string]string{"general": "Не удалось обновить новость: " + err.Error()}
		s.renderer.Page(w, "news_form", data)
		return
	}

	SetFlash(w, Flash{Type: "success", Message: "Новость успешно обновлена"})
	http.Redirect(w, r, "/admin/news", http.StatusFound)
}

// parseNewsForm parses and validates news form data.
// Returns the populated News struct and any field errors.
func parseNewsForm(r *http.Request) (*models.News, map[string]string) {
	errs := make(map[string]string)
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		errs["title"] = "Заголовок обязателен"
	} else if len([]rune(title)) > 500 {
		errs["title"] = "Заголовок не должен превышать 500 символов"
	}

	var date time.Time
	dateStr := r.FormValue("date")
	if dateStr == "" {
		errs["date"] = "Дата обязательна"
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			errs["date"] = "Неверный формат даты (ГГГГ-ММ-ДД)"
		}
	}

	announce := strings.TrimSpace(r.FormValue("announce"))
	description := strings.TrimSpace(r.FormValue("description"))
	n := &models.News{
		Date:        date,
		Title:       title,
		Announce:    announce,
		Description: description,
	}

	return n, errs
}
