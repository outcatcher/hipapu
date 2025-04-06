package main

import (
	"context"
	"os"
	"os/user"

	"github.com/urfave/cli/v3"
)

const (
	commandAdd  = "add"
	commandSync = "sync"
	commandList = "list"
)

//nolint:funlen  // todo: rewrite this abomination
func main() {
	handlers := new(actionHandlers)

	cmd := &cli.Command{
		Name:  "hipapu",
		Usage: "HiPaPu is a tool for automatic updates of binary packages installed from GitHub",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration file path",
				Sources:     cli.ValueSourceChain{},
				Required:    false,
				Value:       defaultConfigPath(),
				Destination: &handlers.configPath,
				Aliases:     []string{"c"},
				TakesFile:   true,
				Config:      cli.StringConfig{},
				OnlyOnce:    true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  commandAdd,
				Usage: "Adds package to the watchlist",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "path",
						Usage:       "(required) File to watch. Will be created if doesn't exist.",
						Required:    true,
						Destination: &handlers.filePath,
						Aliases:     []string{"p"},
						TakesFile:   true,
						OnlyOnce:    true,
					},
					&cli.StringFlag{
						Name:        "repo",
						Usage:       "(required) Repo to watch. Must exist.",
						Required:    true,
						Destination: &handlers.repoPath,
						Aliases:     []string{"r"},
						TakesFile:   true,
						OnlyOnce:    true,
					},
				},
				EnableShellCompletion: true,
				Action:                handlers.add,
				Suggest:               true,
			},
			{
				Name:                  commandSync,
				Usage:                 "Synchronize packages from repos",
				EnableShellCompletion: true,
				Action:                handlers.sync,
				Suggest:               true,
			},
			{
				Name:                  commandList,
				Usage:                 "List existing installations",
				EnableShellCompletion: true,
				Action:                handlers.list,
				Suggest:               true,
			},
		},
		DefaultCommand:        commandList,
		EnableShellCompletion: true,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		println(err.Error())

		os.Exit(1)
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
