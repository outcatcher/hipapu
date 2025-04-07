// Copyright (C) 2025  Anton Kachurin
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const defaultDirPerm = 0755

// SaveConfig creates new or rewrites old configuration file.
func (c *Config) SaveConfig() error {
	path, err := filepath.Abs(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to get config abs path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path) //nolint:gosec  // Abs already calling Clean for result
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			slog.Error("failed to close config file", "error", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to save config : %w", err)
	}

	return nil
}
