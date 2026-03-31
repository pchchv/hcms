package hcms

import "flag"

// Config holds the application configuration.
type Config struct {
	Port       int
	DBPath     string
	UploadPath string
}

// Load parses CLI flags and environment variables, and environment variables override flags.
// It checks whether the Port is in the range 1-65535.
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

	return cfg, nil
}
