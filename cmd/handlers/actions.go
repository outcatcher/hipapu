// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"io"

	"github.com/outcatcher/hipapu/app"
)

// DefaultCommandName - name of default command.
const DefaultCommandName = commandNameList

type application interface {
	// List lists all existing installations.
	List(ctx context.Context) ([]app.Installation, error)
	// Add adds installation to the list. Rewrites lockfile.
	Add(remoteURL, localPath string) error
	// Synchronize runs synchronization of all new releases replacing local files reporting the progress.
	Synchronize(ctx context.Context, writer io.Writer) error
	// UpdateLockVersion updates version if needed and creates backup.
	UpdateLockfileVersion() error
}

// ActionHandlers handle CLI actions.
type ActionHandlers struct {
	filePath, repoPath, lockPath string

	app application
}
