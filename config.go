package hcms

import (
	"errors"
	"flag"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	Port       int
	DBPath     string
	UploadPath string
}

// Load parses CLI flags and environment variables,
// and checks whether the Port is in the range 1-65535.
func Load() (*Config, error) {
	var (
		flagPort       int
		flagDBPath     string
		flagUploadPath string
	)

	flag.IntVar(&flagPort, "port", 8080, "HTTP listen port")
	flag.StringVar(&flagDBPath, "db", "./cms.db", "SQLite database path")
	flag.StringVar(&flagUploadPath, "upload", "./uploads", "Upload directory path")
	flag.Parse()

	cfg := &Config{
		Port:       flagPort,
		DBPath:     flagDBPath,
		UploadPath: flagUploadPath,
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, errors.New("port " + strconv.Itoa(cfg.Port) + "is out of valid range 1-65535")
	}

	return cfg, nil
}
