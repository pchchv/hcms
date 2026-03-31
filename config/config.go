package config

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	// Read ENV.
	// If there aren't any, use the hard-coded defaults.
	envPort, err := getEnvInt("CMS_PORT", 8080)
	if err != nil {
		return nil, err
	}
	envDB := getEnv("CMS_DB_PATH", "./cms.db")
	envUpload := getEnv("CMS_UPLOAD_PATH", "./uploads")

	// Initialize the config using flags.
	// The values ​​from ENV are now the default.
	// If user enters the --port flag, it will OVERWRITE the value from ENV.
	cfg := &Config{}
	flag.IntVar(&cfg.Port, "port", envPort, "HTTP listen port")
	flag.StringVar(&cfg.DBPath, "db", envDB, "SQLite database path")
	flag.StringVar(&cfg.UploadPath, "upload", envUpload, "Upload directory path")
	flag.Parse()

	// Validate Port.
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("port %d is out of valid range 1-65535", cfg.Port)
	}

	// Infrastructure preparation.
	if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory %q: %w", cfg.UploadPath, err)
	}

	// Log the final parameters.
	log.Printf("Configuration loaded: Port=%d, DB=%s, Uploads=%s", cfg.Port, cfg.DBPath, cfg.UploadPath)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid environment variable %s=%q: %w", key, v, err)
	}

	return i, nil
}
