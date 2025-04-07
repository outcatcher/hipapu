// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package config_test

import (
	"testing"

	"github.com/outcatcher/hipapu/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	path := tmpDir + "/test.config"
	testInstall := config.Installation{
		RepoURL:   "https://github.com/outcatcher/hipapu",
		LocalPath: tmpDir + "/test.txt",
	}

	cfg, err := config.LoadConfig(path)
	require.NoError(t, err)

	require.NoError(t, cfg.Add(testInstall))

	newCfg, err := config.LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, newCfg)

	require.Len(t, newCfg.Installations, 1)
	require.Equal(t, testInstall, newCfg.Installations[0])
}
