package app

import "fmt"

// UpdateLockVersion updates version if needed and creates backup.
func (a *Application) UpdateLockVersion() error {
	if err := a.lockfile.UpdateVersion(); err != nil {
		return fmt.Errorf("failed to update lock version: %w", err)
	}

	return nil
}
