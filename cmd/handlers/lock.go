// Copyright (C) 2025  Anton Kachurin

package handlers

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// FlagLockfile - handle '--config' flag.
func (h *ActionHandlers) FlagLockfile() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "lock",
		Usage:       "Lockfile path",
		Sources:     cli.ValueSourceChain{},
		Required:    false,
		Value:       defaultConfigPath(),
		Destination: &h.lockPath,
		TakesFile:   true,
	}
}

func defaultConfigPath() string {
	basePath := "."

	usr, _ := user.Current()
	if usr != nil {
		basePath = usr.HomeDir
	}

	return basePath + "/.config/hipapu/lock.json"
}

func (*ActionHandlers) checkAndMigrateLockIfExists(_ io.Reader, out io.Writer) error {
	newPath := defaultConfigPath()

	oldPath := filepath.Clean(filepath.Dir(newPath) + "/config.json")

	stat, _ := os.Stat(oldPath)
	if stat == nil {
		return nil
	}

	_, _ = out.Write([]byte("Old configuration detected in default location\n"))
	_, _ = out.Write([]byte("Configuration will be moved to " + newPath + "\n"))
	_, _ = out.Write([]byte("Backward-compatible fixes will be applied automatically"))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("checkAndMigrateLockIfExists failed: %w", err)
	}

	return nil
}
