// Copyright (C) 2025  Anton Kachurin
package installations

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

const minID = 0x10000

// Add adds installation to the list. Rewrites configuration file.
func (c *Lock) Add(installation Installation) error {
	if installation.ID == "" {
		installation.ID = strconv.FormatInt(rand.Int64N(minID)+minID, 16) //nolint:gosec
	}

	c.lockData.Installations = append(c.lockData.Installations, installation)

	if err := c.saveInstallations(false); err != nil {
		return fmt.Errorf("failed to save added installation: %w", err)
	}

	return nil
}
