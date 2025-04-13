// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

const commandNameAdd = "add"

// CommandAdd handle 'add' subcommand.
func (h *ActionHandlers) CommandAdd() *cli.Command {
	return &cli.Command{
		Name:  commandNameAdd,
		Usage: "Adds package to the watchlist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Usage:       "(required) File to watch. Will be created if doesn't exist.",
				Required:    true,
				Destination: &h.filePath,
				Aliases:     []string{"p"},
				TakesFile:   true,
				OnlyOnce:    true,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "(required) Repo to watch. Must exist.",
				Required:    true,
				Destination: &h.repoPath,
				Aliases:     []string{"r"},
				TakesFile:   true,
				OnlyOnce:    true,
			},
		},
		EnableShellCompletion: true,
		Action:                h.addHandler,
		Suggest:               true,
	}
}

func (h *ActionHandlers) addHandler(_ context.Context, cmd *cli.Command) error {
	if err := h.app.Add(h.repoPath, h.filePath); err != nil {
		return fmt.Errorf("error during installation addition: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.Writer, "Added!")

	return nil
}
