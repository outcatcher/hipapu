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
}

// List lists all existing installations.
//
// It would be better to add a type for return, but it's such a waste of code.
func (a *Application) List(ctx context.Context) ([]Installation, error) {
	installations := a.config.GetInstallations()

	result := make([]Installation, len(installations))

	for idx, installation := range installations {
		release, err := a.remote.GetLatestRelease(ctx, installation.RepoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get release info: %w", err)
		}

		file, err := a.files.GetFileInfo(installation.LocalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation file info: %w", err)
		}

		result[idx] = Installation{
			Release:   release,
			LocalFile: file,
		}
	}

	return result, nil
}
