//go:build test
// +build test

package remote_test

import (
	"os"
	"testing"

	"github.com/outcatcher/hipapu/internal/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL      = "https://github.com/outcatcher/hipapu/releases/download/testrelease/release.txt"
	testFileName = "release.txt"
)

func TestDownloadFile(t *testing.T) {
	t.Parallel()

	filePath := t.TempDir() + "/" + testFileName

	t.Cleanup(func() {
		require.NoError(t, os.Remove(filePath))
	})

	file, err := os.Create(filePath)
	require.NoError(t, err)

	// closing file after that is important, so using 'assert' instead of 'require'
	assert.NoError(t, remote.DownloadFile(t.Context(), testURL, file))

	require.NoError(t, file.Close())

	stat, err := os.Stat(filePath)
	require.NoError(t, err)

	require.Positive(t, stat.Size(), "file is empty")
}
