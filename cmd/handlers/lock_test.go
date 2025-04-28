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

	const expectedLockPath = "test-lock-path-add"

	hdl := &ActionHandlers{}

	require.NoError(t, hdl.FlagLockfile().Set("", expectedLockPath))

	require.Equal(t, expectedLockPath, hdl.lockPath)
}
