package hcms

import (
	"flag"
	"fmt"
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
	var (
		flagPort       int
		flagDBPath     string
		flagUploadPath string
	)

	// Initialize the config using flags.
	cfg := &Config{}
	flag.IntVar(&cfg.Port, "port", flagPort, "HTTP listen port")
	flag.StringVar(&cfg.DBPath, "db", flagDBPath, "SQLite database path")
	flag.StringVar(&cfg.UploadPath, "upload", flagUploadPath, "Upload directory path")
	flag.Parse()

	// Validate Port.
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("port %d is out of valid range 1-65535", cfg.Port)
	}

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
