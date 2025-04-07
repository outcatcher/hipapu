// Copyright (C) 2025  Anton Kachurin
package app

import (
	"github.com/outcatcher/hipapu/internal/config"
)

// List lists all existing installations.
//
// It would be better to add a type for return, but it's such a waste of code.
func (a *Application) List() []config.Installation {
	return a.config.GetInstallations()
}
