package handlers

import "testing"

func setupAuthServer(t *testing.T) *Server {
	t.Helper()
	s, database := setupTestServer(t)
	s.renderer = NewRenderer(templatesDir)
	_ = database
	return s
}
