// Copyright (C) 2025  Anton Kachurin
package lock

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	defaultDirPerm = 0755
	currentVersion = "1.1"
)

// SaveInstallations creates new or rewrites old installations file.
func (l *Lock) SaveInstallations() error {
	path, err := filepath.Abs(l.filePath)
	if err != nil {
		return fmt.Errorf("failed to get lockfile abs path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path) //nolint:gosec  // Abs already calling Clean for result
	if err != nil {
		return fmt.Errorf("failed to create lockfile: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			slog.Error("failed to close lockfile", "error", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	l.lockData.LockVersion = currentVersion

	if err := encoder.Encode(l.lockData); err != nil {
		return fmt.Errorf("failed to save lockfile : %w", err)
	}

	return nil
}
