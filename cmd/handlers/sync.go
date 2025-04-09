// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/outcatcher/hipapu/app"
	"github.com/urfave/cli/v3"
)

const commandNameSync = "sync"

// SyncCommand handle 'sync' subcommand.
func (h *ActionHandlers) SyncCommand() *cli.Command {
	return &cli.Command{
		Name:                  commandNameSync,
		Usage:                 "Synchronize packages from repos",
		EnableShellCompletion: true,
		Action:                h.sync,
		Suggest:               true,
	}
}

func (h *ActionHandlers) sync(ctx context.Context, _ *cli.Command) error {
	application, err := app.New(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	if err := application.Synchronize(ctx); err != nil {
		if errors.Is(err, app.ErrEmptyInstallationList) {
			fmt.Println("Empty installation list. Nothing to synchronize.")

			return nil
		}

		return fmt.Errorf("error during synchnorization: %w", err)
	}

	println("Sync finished!")

	return nil
}
