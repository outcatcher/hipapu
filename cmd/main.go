// Copyright (C) 2025  Anton Kachurin
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

	copyright = `(C) 2025  Anton Kachurin

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
`
)

//nolint:funlen  // todo: rewrite this abomination
func main() {
	handlers := new(actionHandlers)

	cmd := &cli.Command{
		Name:      "hipapu",
		Usage:     "HiPaPu is a tool for automatic updates of binary packages installed from GitHub",
		Copyright: copyright,
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
