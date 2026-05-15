package handlers

import (
	"bytes"

	"github.com/yuin/goldmark"
)

var md = goldmark.New()

// renderMarkdown converts a Markdown string to HTML.
func renderMarkdown(src string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
