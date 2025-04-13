// Copyright (C) 2025  Anton Kachurin
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/outcatcher/hipapu/cmd/handlers"
	"github.com/urfave/cli/v3"
)

const (
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

func main() {
	hdl := new(handlers.ActionHandlers)

	cmd := &cli.Command{
		Name:  "hipapu",
		Usage: "HiPaPu is a tool for automatic updates of binary packages installed from GitHub",
		Flags: []cli.Flag{
			hdl.FlagConfig(),
		},
		Commands: []*cli.Command{
			hdl.CommandAdd(),
			hdl.CommandSync(),
			hdl.CommandList(),
		},
		Before:                hdl.Before,
		DefaultCommand:        handlers.DefaultCommandName,
		EnableShellCompletion: true,
		Copyright:             copyright,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sig

		cancel()
	}()

	if err := cmd.Run(ctx, os.Args); err != nil {
		println(err.Error())

		os.Exit(1)
	}
}
