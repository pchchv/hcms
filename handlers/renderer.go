package handlers

import "io/fs"

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
