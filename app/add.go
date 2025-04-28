// Copyright (C) 2025  Anton Kachurin
package app

import (
	"fmt"

	"github.com/outcatcher/hipapu/internal/installations"
)

// Add adds installation to the list. Rewrites lockfile.
func (a *Application) Add(remoteURL, localPath string) error {
	if err := a.lock.Add(installations.Installation{
		RepoURL:   remoteURL,
		LocalPath: localPath,
	}); err != nil {
		return fmt.Errorf("failed to add installation: %w", err)
	}

	return nil
}
