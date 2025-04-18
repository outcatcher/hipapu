// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package app_test

import (
	"bytes"
	"io"
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
		expectecdURL    = "https://github.com/outcatcher/asdfasdf"
		expectedSkipURL = expectecdURL + ".skip.me"

		expectedConfigPath    = "./config.cfg"
		expectedLocalFilename = "localFilePath"

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
		{
			RepoURL:   expectedSkipURL,
			LocalPath: expectedLocalPath,
			SkipSync:  true,
		},
	})

	apk.WithConfig(mockCfg)

	ctx := t.Context()

	mockRemote := mocks.NewMockremoteClient(t)
	mockRemote.
		On("GetLatestRelease", ctx, expectecdURL).
		Return(&remote.Release{
			Name:        "123",
			Description: "234",
			PublishedAt: time.Now(),
			Assets: []remote.Asset{{
				Filename:    expectedLocalFilename,
				DownloadURL: expectecdDownloadURL,
			}},
			RepoURL: expectecdURL,
		}, nil)
	mockRemote.
		On("GetLatestRelease", ctx, expectedSkipURL).
		Return(&remote.Release{
			Name:        "123",
			Description: "234",
			PublishedAt: time.Now(),
			Assets: []remote.Asset{{
				Filename:    expectedLocalFilename,
				DownloadURL: expectecdDownloadURL,
			}},
			RepoURL: expectedSkipURL,
		}, nil)
	mockRemote.
		On("DownloadFile", ctx, expectecdDownloadURL, mock.Anything).
		Return(nil)

	apk.WithRemote(mockRemote)

	mockLocal := mocks.NewMocklocalFiles(t)
	mockLocal.On("GetFileInfo", expectedLocalPath).Return(&local.FileInfo{
		Name:     expectedLocalFilename,
		FilePath: expectedLocalPath,
	}, nil)

	apk.WithFiles(mockLocal)

	require.NoError(t, apk.Add(expectecdURL, expectedLocalPath))

	outBuffer := new(bytes.Buffer)

	require.NoError(t, apk.Synchronize(ctx, outBuffer))

	line1, err := outBuffer.ReadString('\n')
	require.NoError(t, err)

	require.Contains(t, line1, expectecdURL)

	line2, err := io.ReadAll(outBuffer)
	require.NoError(t, err)

	require.Contains(t, string(line2), "Skipping")
	require.Contains(t, string(line2), expectedSkipURL)
}
