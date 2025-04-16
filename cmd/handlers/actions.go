// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"
	"io"

	"github.com/outcatcher/hipapu/app"
	"github.com/urfave/cli/v3"
)

// DefaultCommandName - name of default command.
const DefaultCommandName = commandNameList

type application interface {
	// List lists all existing installations.
	List(ctx context.Context) ([]app.Installation, error)
	// Add adds installation to the list. Rewrites configuration file.
	Add(remoteURL, localPath string) error
	// Synchronize runs synchronization of all new releases replacing local files reporting the progress.
	Synchronize(ctx context.Context, writer io.Writer) error
}

// ActionHandlers handle CLI actions.
type ActionHandlers struct {
	filePath, repoPath, configPath string

	app application
}

// Before is a before function for the command handlers.
func (h *ActionHandlers) Before(ctx context.Context, _ *cli.Command) (context.Context, error) {
	application, err := app.New(h.configPath)
	if err != nil {
		return ctx, fmt.Errorf("failed to init app: %w", err)
	}

	h.app = application

	return ctx, nil
}
