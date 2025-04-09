// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"

	"github.com/outcatcher/hipapu/app"
	"github.com/urfave/cli/v3"
)

const commandNameList = "list"

// ListCommand handle 'list' subcommand.
func (h *ActionHandlers) ListCommand() *cli.Command {
	return &cli.Command{
		Name:                  commandNameList,
		Usage:                 "List existing installations",
		EnableShellCompletion: true,
		Action:                h.list,
		Suggest:               true,
	}
}

func (h *ActionHandlers) list(context.Context, *cli.Command) error {
	application, err := app.New(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	installations := application.List()

	fmt.Println("Installations:")

	for i, installation := range installations {
		fmt.Printf("  %d) %s <---> %s\n", i+1, installation.RepoURL, installation.LocalPath)
	}

	return nil
}
