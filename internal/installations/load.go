// Copyright (C) 2025  Anton Kachurin
package installations

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

const selfRepo = "https://github.com/outcatcher/hipapu"

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

	l.setDefaults()
	l.filePath = path

	return nil
}

func (l *Lock) setDefaults() {
	hasSelf := slices.ContainsFunc(l.lockData.Installations, func(inst Installation) bool {
		return inst.RepoURL == selfRepo
	})

	if !hasSelf {
		l.lockData.Installations = appendSelf(l.lockData.Installations)
	}
}

func appendSelf(installs []Installation) []Installation {
	selfPath, err := os.Executable() // os.Args[0] won't work if binary is in PATH
	if err != nil {
		return installs
	}

	return append(installs, Installation{
		ID:              "0",
		RepoURL:         selfRepo,
		LocalPath:       selfPath, // register self-update
		KeepLastVersion: true,     // safety measure for broken releases
		SkipSync:        true,     // todo: remove after KeepLastVersion implemented
	})
}
