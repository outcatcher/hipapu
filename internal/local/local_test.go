//go:build test
// +build test

package local_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/outcatcher/hipapu/internal/local"
	"github.com/stretchr/testify/require"
)

const (
	testData     = "6630cfdb-2123-4802-8ce2-f50336d56c18"
	testFilename = "test.txt"
)

func TestGetFileInfo(t *testing.T) {
	t.Parallel()

	fullPath := t.TempDir() + "/" + testFilename

	file, err := os.Create(fullPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, os.Remove(fullPath))
	})

	_, err = file.WriteString(testData)
	require.NoError(t, err)

	require.NoError(t, file.Close())

	updated := time.Now().Truncate(time.Millisecond)

	t.Log(file.Name())

	files := local.Files{}

	info, err := files.GetFileInfo(file.Name())
	require.NoError(t, err)

	require.True(t, strings.HasSuffix(info.FilePath, "/"+testFilename))
	require.InDelta(t, updated.Unix(), info.LastModified.Unix(), 1)
}
