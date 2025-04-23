// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package lock_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/outcatcher/hipapu/internal/lock"
	"github.com/stretchr/testify/require"
)

const (
	repoURL     = "https://github.com/outcatcher/hipapu-release-testing"
	oldFileData = `{
  "download_location": "/tmp/hihapu",
  "installations": [
    {
      "repo_url": "%s",
      "local_path": "%s"
    }
  ]
}
`
	superNewFileData = `{
  "version": "v2200"
}
`
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	path := tmpDir + "/test.lock"
	localPath := tmpDir + "/test.bin"

	origFile, err := os.Create(path)
	require.NoError(t, err)

	_, err = io.WriteString(origFile, fmt.Sprintf(oldFileData, repoURL, localPath))
	require.NoError(t, err)

	require.NoError(t, origFile.Close())

	locks := new(lock.Lock)

	require.NoError(t, locks.LoadInstallations(path))

	require.Len(t, locks.GetInstallations(), 2)

	require.NoError(t, locks.UpdateVersion())

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	require.Len(t, entries, 2)

	require.NoError(t, locks.UpdateVersion())

	require.Len(t, entries, 2)
}
func TestUpdateNew(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	path := tmpDir + "/test.lock"

	origFile, err := os.Create(path)
	require.NoError(t, err)

	_, err = io.WriteString(origFile, superNewFileData)
	require.NoError(t, err)

	require.NoError(t, origFile.Close())

	locks := new(lock.Lock)

	require.NoError(t, locks.LoadInstallations(path))

	require.Len(t, locks.GetInstallations(), 1)

	require.ErrorIs(t, locks.UpdateVersion(), lock.ErrUnsupportedVersion)
}
