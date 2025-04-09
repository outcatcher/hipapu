package handlers

import (
	"os/user"

	"github.com/urfave/cli/v3"
)

func (h *ActionHandlers) ConfigFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "config",
		Usage:       "Configuration file path",
		Sources:     cli.ValueSourceChain{},
		Required:    false,
		Value:       defaultConfigPath(),
		Destination: &h.configPath,
		Aliases:     []string{"c"},
		TakesFile:   true,
		OnlyOnce:    true,
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
