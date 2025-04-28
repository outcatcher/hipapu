// Copyright (C) 2025  Anton Kachurin

package handlers

import (
	"os/user"

	"github.com/urfave/cli/v3"
)

// FlagLockfile - handle '--config' flag.
func (h *ActionHandlers) FlagLockfile() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "lock",
		Usage:       "Lockfile path",
		Sources:     cli.ValueSourceChain{},
		Required:    false,
		Value:       defaultLockfilePath(),
		Destination: &h.lockPath,
		TakesFile:   true,
	}
}

func defaultLockfilePath() string {
	basePath := "."

	usr, _ := user.Current()
	if usr != nil {
		basePath = usr.HomeDir
	}

	return basePath + "/.config/hipapu/lock.json"
}
