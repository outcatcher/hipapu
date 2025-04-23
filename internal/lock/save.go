// Copyright (C) 2025  Anton Kachurin
package lock

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultDirPerm = 0755
	currentVersion = "v1.1"
)

func createBackup(path string) error {
	stat, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat original file for backup: %w", err)
	}

	if stat == nil {
		return nil
	}

	srcFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	backupFile, err := os.Create(path + ".backup." + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer func() {
		_ = backupFile.Close()
	}()

	_, err = io.Copy(backupFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy backup: %w", err)
	}

	return nil
}

// saveInstallations creates new or rewrites old installations file.
func (l *Lock) saveInstallations(withBackup bool) error {
	path, err := filepath.Abs(l.filePath)
	if err != nil {
		return fmt.Errorf("failed to get lockfile abs path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if withBackup {
		if err := createBackup(path); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
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
