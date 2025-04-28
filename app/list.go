// Copyright (C) 2025  Anton Kachurin
package app

import (
	"context"
	"fmt"

	"github.com/outcatcher/hipapu/internal/local"
	"github.com/outcatcher/hipapu/internal/remote"
)

// Installation to be synced.
type Installation struct {
	// Remote repo release information
	Release *remote.Release
	// Local file info
	LocalFile *local.FileInfo

	// Skipping sync for installation
	SkipSync bool
}

// List lists all existing installations.
func (a *Application) List(ctx context.Context) ([]Installation, error) {
	installations := a.lock.GetInstallations()

	result := make([]Installation, len(installations))

	for i, installation := range installations {
		release, err := a.remote.GetLatestRelease(ctx, installation.RepoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get release info: %w", err)
		}

		file, err := a.files.GetFileInfo(installation.LocalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation file info: %w", err)
		}

		result[i] = Installation{
			Release:   release,
			LocalFile: file,
			SkipSync:  installation.SkipSync,
		}
	}

	return result, nil
}
