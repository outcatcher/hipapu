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

func (h *ActionHandlers) list(_ context.Context, cmd *cli.Command) error {
	installations := h.app.List()

	_, _ = fmt.Fprintln(cmd.Writer, "Installations:")

	for i, installation := range installations {
		_, _ = fmt.Fprintf(cmd.Writer, "  %d) %s <---> %s\n", i+1, installation.RepoURL, installation.LocalPath)
	}

	return nil
}
