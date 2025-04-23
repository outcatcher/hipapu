// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package handlers //nolint:testpackage  // don't want to make useless tests overly complex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	const expectedConfigPath = "test-config-path-add"

	hdl := &ActionHandlers{}

	require.NoError(t, hdl.FlagLockfile().Set("", expectedConfigPath))

	require.Equal(t, expectedConfigPath, hdl.lockPath)
}
