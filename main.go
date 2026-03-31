package hcms

import (
	"log"
	"os"
)

func main() {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Ensure uploads directory exists.
	if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
		log.Fatalf("create upload dir: %v", err)
	}
}
