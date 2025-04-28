// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package installations_test

import (
	"testing"

	"github.com/outcatcher/hipapu/internal/installations"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	path := tmpDir + "/test.lock"
	testInstall := installations.Installation{
		ID:        "test",
		RepoURL:   "https://github.com/outcatcher/hipapu-release-testing/",
		LocalPath: tmpDir + "/test.txt",
	}

	lock := new(installations.Lock)

	require.NoError(t, lock.LoadInstallations(path))

	require.Empty(t, lock.GetInstallations())

	require.NoError(t, lock.Add(testInstall))

	require.NoError(t, lock.LoadInstallations(path))
	require.NotNil(t, lock)

	require.Len(t, lock.GetInstallations(), 1)
	require.Contains(t, lock.GetInstallations(), testInstall)
}
