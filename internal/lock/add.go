// Copyright (C) 2025  Anton Kachurin
package lock

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

const minID = 0x10000

// SaveInstallations creates new or rewrites old installations file.
func (c *Lock) Add(installation Installation) error {
	if installation.ID == "" {
		installation.ID = strconv.FormatInt(rand.Int64N(minID)+minID, 16) //nolint:gosec
	}

	c.lockData.Installations = append(c.lockData.Installations, installation)

	if err := c.SaveInstallations(); err != nil {
		return fmt.Errorf("failed to save added installation: %w", err)
	}

	return nil
}
