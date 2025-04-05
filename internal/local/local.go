package local

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var errCantSyncDir = errors.New("sync of directories is not supported")

// FileInfo - information about local tracked file.
type FileInfo struct {
	Name     string
	FilePath string
	// Zero if file doesn't exist
	LastModified time.Time
}

// Files - service for local operations.
type Files struct{}

// GetFileInfo returns info on the local file.
func (*Files) GetFileInfo(filePath string) (*FileInfo, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get abs path: %w", err)
	}

	result := &FileInfo{
		Name:     filepath.Base(filePath),
		FilePath: filePath,
	}

	stat, err := os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if stat == nil {
		return result, nil
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("%s: %w", absPath, errCantSyncDir)
	}

	result.LastModified = stat.ModTime()

	return result, nil
}
