//go:build test
// +build test

package remote_test

import (
	"strings"
	"testing"

	"github.com/outcatcher/hipapu/internal/remote"
	"github.com/stretchr/testify/require"
)

const (
	// todo: create/remove release during test setup/terdown
	testOwner         = "outcatcher"
	testRepo          = "hipapu"
	testReleaseName   = "TestRelease"
	testAssetFileName = "release.txt"
)

func TestGetLatestRelease(t *testing.T) {
	t.Parallel()

	client, err := remote.New()
	require.NoError(t, err)

	release, err := client.GetLatestRelease(t.Context(), testOwner, testRepo)
	require.NoError(t, err)

	require.NotNil(t, release)

	require.Equal(t, testReleaseName, release.Name)
	require.Len(t, release.Assets, 1) // 2 default zips not included

	require.Equal(t, testAssetFileName, release.Assets[0].Filename)
	require.True(t, strings.HasSuffix(release.Assets[0].DownloadURL, testAssetFileName))
}
