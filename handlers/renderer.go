package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
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
