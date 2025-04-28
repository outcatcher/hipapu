// Copyright (C) 2025  Anton Kachurin
package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/outcatcher/hipapu/app"
	"github.com/urfave/cli/v3"
)

// Before is a before function for the command handlers.
func (h *ActionHandlers) Before(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	if err := h.checkAndMigrateLockIfExists(cmd.Reader, cmd.Writer); err != nil {
		return ctx, fmt.Errorf("failed migrate old default config to lock: %w", err)
	}

	application, err := app.New(h.lockPath)
	if err != nil {
		return ctx, fmt.Errorf("failed to init app: %w", err)
	}

	h.app = application

	if err := h.app.UpdateLockfileVersion(); err != nil {
		return ctx, fmt.Errorf("failed to update lockfile: %w", err)
	}

	return ctx, nil
}

func (*ActionHandlers) checkAndMigrateLockIfExists(_ io.Reader, out io.Writer) error {
	newPath := defaultLockfilePath()

	oldPath := filepath.Clean(filepath.Dir(newPath) + "/config.json")

	stat, _ := os.Stat(oldPath)
	if stat == nil {
		return nil
	}

	_, _ = out.Write([]byte("Old configuration detected in default location\n"))
	_, _ = out.Write([]byte("Configuration will be moved to " + newPath + "\n"))
	_, _ = out.Write([]byte("Backward-compatible fixes will be applied automatically"))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	return nil
}
