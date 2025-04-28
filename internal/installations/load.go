// Copyright (C) 2025  Anton Kachurin
package installations

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// LoadInstallations loads installations data from file.
func (l *Lock) LoadInstallations(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get lockfile abs path: %w", err)
	}

	file, err := os.Open(path) //nolint:gosec  // Abs already calling Clean for result
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to open lockfile file to read: %w", err)
	}

	// if file exists, load lockfile
	if file != nil {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				slog.Error("failed to close lockfile (read)", "error", closeErr)
			}
		}()

		if err := json.NewDecoder(file).Decode(&l.lockData); err != nil {
			return fmt.Errorf("failed to read lockfile : %w", err)
		}
	}

	l.filePath = path

	return nil
}
