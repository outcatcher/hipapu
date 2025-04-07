// Copyright (C) 2025  Anton Kachurin
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/outcatcher/hipapu/app"
	"github.com/urfave/cli/v3"
)

type actionHandlers struct {
	filePath, repoPath, configPath string
}

func (h *actionHandlers) add(context.Context, *cli.Command) error {
	application, err := app.New(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	if err := application.Add(h.repoPath, h.filePath); err != nil {
		return fmt.Errorf("error during installation addition: %w", err)
	}

	return nil
}

func (h *actionHandlers) sync(ctx context.Context, _ *cli.Command) error {
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

	return nil
}

func (h *actionHandlers) list(context.Context, *cli.Command) error {
	application, err := app.New(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	installations := application.List()

	fmt.Println("Installations:")

	for i, installation := range installations {
		fmt.Printf("#%d) %s <---> %s\n", i+1, installation.RepoURL, installation.LocalPath)
	}

	return nil
}
