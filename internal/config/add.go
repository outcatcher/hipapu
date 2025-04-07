// Copyright (C) 2025  Anton Kachurin
package config

import "fmt"

// Add adds installation to the list. Rewrites configuration file.
func (c *Config) Add(installation Installation) error {
	c.Installations = append(c.Installations, installation)

	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save added installation: %w", err)
	}

	return nil
}
