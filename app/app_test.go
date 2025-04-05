//go:build test
// +build test

package app_test

import (
	"testing"
	"time"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/app/mocks"
	"github.com/outcatcher/hipapu/internal/config"
	"github.com/outcatcher/hipapu/internal/local"
	"github.com/outcatcher/hipapu/internal/remote"
	"github.com/stretchr/testify/mock"
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

	apk := new(app.Application)

	const (
		expectecdURL = "https://github.com/outcatcher/asdfasdf"

		expectedConfigPath    = "./config.cfg"
		expectedLocalFilename = "localFilePath"

		expectedTestOwner = "outcatcher"
		expectedTestRepo  = "asdfasdf"

		expectecdDownloadURL = "https://adfsgijnasdfgj.test"
	)

	expectedLocalPath := t.TempDir() + "/" + expectedLocalFilename

	mockCfg := mocks.NewMockcfg(t)
	mockCfg.On("Add", config.Installation{
		RepoURL:   expectecdURL,
		LocalPath: expectedLocalPath,
	}).Return(nil)
	mockCfg.On("GetInstallations").Return([]config.Installation{
		{
			RepoURL:   expectecdURL,
			LocalPath: expectedLocalPath,
		},
	})

	apk.WithConfig(mockCfg)

	ctx := t.Context()

	mockRemote := mocks.NewMockremoteClient(t)
	mockRemote.
		On("GetLatestRelease", ctx, expectedTestOwner, expectedTestRepo).
		Return(&remote.Release{
			Name:        "123",
			Description: "234",
			PublishedAt: time.Now(),
			Assets: []remote.Asset{{
				Filename:    expectedLocalFilename,
				DownloadURL: expectecdDownloadURL,
			}},
		}, nil)
	mockRemote.
		On("DownloadFile", ctx, expectecdDownloadURL, mock.AnythingOfType("*os.File")).
		Return(nil)

	apk.WithRemote(mockRemote)

	mockLocal := mocks.NewMocklocalFiles(t)
	mockLocal.On("GetFileInfo", expectedLocalPath).Return(&local.FileInfo{
		Name:     expectedLocalFilename,
		FilePath: expectedLocalPath,
	}, nil)

	apk.WithFiles(mockLocal)

	require.NoError(t, apk.Add(expectecdURL, expectedLocalPath))

	require.NoError(t, apk.Synchronize(ctx))
}
