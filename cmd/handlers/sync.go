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

// CommandSync handle 'sync' subcommand.
func (h *ActionHandlers) CommandSync() *cli.Command {
	return &cli.Command{
		Name:                  commandNameSync,
		Usage:                 "Synchronize packages from repos",
		EnableShellCompletion: true,
		Action:                h.sync,
		Suggest:               true,
	}
}

func (h *ActionHandlers) sync(ctx context.Context, cmd *cli.Command) error {
	if err := h.app.Synchronize(ctx, cmd.Writer); err != nil {
		if errors.Is(err, app.ErrEmptyInstallationList) {
			_, _ = fmt.Fprintln(cmd.Writer, "Empty installation list. Nothing to synchronize.")

			return nil
		}

		return fmt.Errorf("error during synchnorization: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.Writer, "Sync finished!")

	return nil
}
