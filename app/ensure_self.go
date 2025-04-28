// Copyright (C) 2025  Anton Kachurin
package app

import (
	"fmt"
	"os"
	"slices"

	"github.com/outcatcher/hipapu/internal/installations"
)

const selfRepo = "https://github.com/outcatcher/hipapu"

func (app *Application) EnsureSelf() error {
	installs := app.lock.GetInstallations()

	hasSelf := slices.ContainsFunc(installs, func(inst installations.Installation) bool {
		return inst.RepoURL == selfRepo
	})

	if hasSelf {
		return nil
	}

	selfPath, err := os.Executable() // os.Args[0] won't work if binary is in PATH
	if err != nil {
		return fmt.Errorf("error ensuring self: %w", err)
	}

	if err := app.lock.Add(installations.Installation{
		ID:              "0",
		RepoURL:         selfRepo,
		LocalPath:       selfPath, // register self-update
		KeepLastVersion: true,     // safety measure for broken releases
		SkipSync:        true,     // todo: remove after KeepLastVersion implemented
	}); err != nil {
		return fmt.Errorf("error ensuring self: %w", err)
	}

	return nil
}
