//go:build test
// +build test

package app_test

import (
	"os"
	"testing"

	"github.com/outcatcher/hipapu/app"
	"github.com/stretchr/testify/require"
)

func TestAppNewNoConfig(t *testing.T) {
	t.Parallel()

	_, err := app.New("")
	require.ErrorIs(t, err, app.ErrNoConfig)
}

func TestAppNewConfig(t *testing.T) {
	t.Parallel()

	filePath := t.TempDir() + "/test.config"

	_, err := app.New(filePath)
	require.NoError(t, err)
}

func TestAppWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	configFilePath := tmpDir + "/test.config"
	localFilePath := tmpDir + "/release.txt"

	apk, err := app.New(configFilePath)
	require.NoError(t, err)

	require.NoError(t, apk.Add("https://github.com/outcatcher/hipapu", localFilePath))

	require.FileExists(t, configFilePath)
	require.NoFileExists(t, localFilePath)

	require.NoError(t, apk.Synchronize(t.Context()))
	require.FileExists(t, localFilePath)

	stat, err := os.Stat(localFilePath)
	require.NoError(t, err)

	require.Positive(t, stat.Size(), "local file not downloaded")
}
