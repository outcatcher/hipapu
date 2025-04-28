// Copyright (C) 2025  Anton Kachurin
package app

import "fmt"

// UpdateLockfileVersion updates version if needed and creates backup.
func (a *Application) UpdateLockfileVersion() error {
	if err := a.lock.UpdateVersion(); err != nil {
		return fmt.Errorf("failed to update lock version: %w", err)
	}

	return nil
}
