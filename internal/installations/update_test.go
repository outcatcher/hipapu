// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package installations_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/outcatcher/hipapu/internal/installations"
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

	lock := new(installations.Lock)

	require.NoError(t, lock.LoadInstallations(path))

	require.Len(t, lock.GetInstallations(), 1)

	require.NoError(t, lock.UpdateVersion())

	require.Len(t, lock.GetInstallations(), 1)

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	require.Len(t, entries, 2) // backup created
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

	lock := new(installations.Lock)

	require.NoError(t, lock.LoadInstallations(path))
	require.ErrorIs(t, lock.UpdateVersion(), installations.ErrUnsupportedVersion)
}
