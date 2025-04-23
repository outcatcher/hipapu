// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package lock_test

import (
	"testing"

	"github.com/outcatcher/hipapu/internal/lock"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	path := tmpDir + "/test.lock"
	testInstall := lock.Installation{
		ID:        "test",
		RepoURL:   "https://github.com/outcatcher/hipapu-release-testing/",
		LocalPath: tmpDir + "/test.txt",
	}

	locks := new(lock.Lock)

	require.NoError(t, locks.LoadInstallations(path))

	require.Len(t, locks.GetInstallations(), 1, "no self installation")

	require.NoError(t, locks.Add(testInstall))

	require.NoError(t, locks.LoadInstallations(path))
	require.NotNil(t, locks)

	require.Len(t, locks.GetInstallations(), 2)
	require.Contains(t, locks.GetInstallations(), testInstall)
}
