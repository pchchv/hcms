package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// funcMap provides template helper functions.
var funcMap = template.FuncMap{
	// formatDate formats a time.Time as "02.01.2006"
	"formatDate": func(t time.Time) string {
		return t.Format("02.01.2006")
	},
	// formatDateTime formats a time.Time as "02.01.2006 15:04"
	"formatDateTime": func(t time.Time) string {
		return t.Format("02.01.2006 15:04")
	},
	// dec decrements an integer
	"dec": func(i int) int { return i - 1 },
	// inc increments an integer
	"inc": func(i int) int { return i + 1 },
	// slice returns a rune-safe substring
	"slice": func(s string, start, end int) string {
		runes := []rune(s)
		if start < 0 {
			start = 0
		}

		if end > len(runes) {
			end = len(runes)
		}

		if start > end {
			return ""
		}

		return string(runes[start:end])
	},
	// safeHTML marks a string as safe HTML (skips escaping)
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
	// contains checks if a string contains a substring
	"contains": strings.Contains,
	// eq is an alias for equality comparison (built-in eq works, but explicit for clarity)
	"eq": func(a, b any) bool { return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b) },
}

// Renderer loads and renders HTML templates.
// When fsys is non-nil it reads from the embedded FS, otherwise it reads from dir on disk.
type Renderer struct {
	dir  string
	fsys fs.FS // non-nil when using embedded assets
}

// NewRenderer creates a Renderer that loads templates from dir on disk (dev mode).
func NewRenderer(dir string) *Renderer {
	return &Renderer{dir: dir}
}

// Page renders a full admin page using the base layout.
// Executes the "base" template.
func (r *Renderer) Page(w http.ResponseWriter, page string, data any) {
	t, err := r.parseLayout(page)
	if err != nil {
		log.Printf("renderer.Page parse error (page=%s): %v", page, err)
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("renderer.Page execute error (page=%s): %v", page, err)
	}
}

// Partial renders a named template from a page file (for HTMX responses).
func (r *Renderer) Partial(w http.ResponseWriter, page, tmplName string, data any) {
	t, err := r.parseLayout(page)
	if err != nil {
		log.Printf("renderer.Partial parse error (page=%s tmpl=%s): %v", page, tmplName, err)
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, tmplName, data); err != nil {
		log.Printf("renderer.Partial execute error (page=%s tmpl=%s): %v", page, tmplName, err)
	}
}

// layoutPaths returns the slash-separated template paths for a full admin page.
func (r *Renderer) layoutPaths(page string) []string {
	return []string{
		"layouts/base.html",
		"partials/nav.html",
		"partials/flash.html",
		"partials/pagination.html",
		"admin/" + page + ".html",
	}
}

// parseLayout parses all layout + page templates and returns the parsed template.
func (r *Renderer) parseLayout(page string) (*template.Template, error) {
	paths := r.layoutPaths(page)
	t := template.New("").Funcs(funcMap)
	if r.fsys != nil {
		return t.ParseFS(r.fsys, paths...)
	}

	// disk mode:
	// convert slash paths to OS paths under r.dir
	fsPaths := make([]string, len(paths))
	for i, p := range paths {
		fsPaths[i] = filepath.Join(r.dir, filepath.FromSlash(p))
	}

	return t.ParseFiles(fsPaths...)
}
