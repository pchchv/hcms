package handlers

import (
	"database/sql"

	"github.com/pchchv/hcms/config"
	"github.com/pchchv/hcms/services"
)

// Server holds the application dependencies for HTTP handlers.
type Server struct {
	db        *sql.DB
	cfg       *config.Config
	renderer  *Renderer
	bitrix    *services.BitrixPool
	uploader  *services.Uploader
	version   string
	buildTime string
}
