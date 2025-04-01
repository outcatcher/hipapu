package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var defaultDownloadLocation = os.TempDir() + "/hihapu"

// LoadConfig loads configuration from file.
func LoadConfig(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get config abs path: %w", err)
	}

	file, err := os.Open(path) //nolint:gosec  // Abs already calling Clean for result
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open config file to read: %w", err)
	}

	configuration := &Config{
		filePath: path,
	}

	// if file exists, load config
	if file != nil {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				slog.Error("failed to close config file (read)", "error", closeErr)
			}
		}()

		if err := json.NewDecoder(file).Decode(configuration); err != nil {
			return nil, fmt.Errorf("failed to read config : %w", err)
		}
	}

	configuration.setDefaults()

	return configuration, nil
}

func (c *Config) setDefaults() {
	if c.DownloadLocation == "" {
		c.DownloadLocation = defaultDownloadLocation
	}
}
