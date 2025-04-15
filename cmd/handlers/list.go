// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

const commandNameList = "list"

// CommandList handle 'list' subcommand.
func (h *ActionHandlers) CommandList() *cli.Command {
	return &cli.Command{
		Name:                  commandNameList,
		Usage:                 "List existing installations",
		EnableShellCompletion: true,
		Action:                h.list,
		Suggest:               true,
	}
}

func (h *ActionHandlers) list(ctx context.Context, cmd *cli.Command) error {
	installations, err := h.app.List(ctx)
	if err != nil {
		return fmt.Errorf("error handling `list` command: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.Writer, "Installations:")

	for i, inst := range installations {
		statusString := fmt.Sprintf("  %d) %s <---> %s", i+1, inst.Release.RepoURL, inst.LocalFile.FilePath)

		if inst.Release.PublishedAt.After(inst.LocalFile.LastModified) {
			statusString += " (HAS UPDATES)"
		}

		_, _ = fmt.Fprintln(cmd.Writer, statusString)
	}

	return nil
}
