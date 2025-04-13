// Copyright (C) 2025  Anton Kachurin

package handlers

import (
	"os/user"

	"github.com/urfave/cli/v3"
)

// FlagConfig - handle '--config' flag.
func (h *ActionHandlers) FlagConfig() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "config",
		Usage:       "Configuration file path",
		Sources:     cli.ValueSourceChain{},
		Required:    false,
		Value:       defaultConfigPath(),
		Destination: &h.configPath,
		Aliases:     []string{"c"},
	}
}

func defaultConfigPath() string {
	basePath := "."

	usr, _ := user.Current()
	if usr != nil {
		basePath = usr.HomeDir
	}

	return basePath + "/.config/hipapu/config.json"
}
