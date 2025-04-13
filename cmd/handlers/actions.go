// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/internal/config"
	"github.com/urfave/cli/v3"
)

// DefaultCommandName - name of default command.
const DefaultCommandName = commandNameList

type application interface {
	List() []config.Installation
	Add(remoteURL, localPath string) error
	Synchronize(ctx context.Context) error
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
