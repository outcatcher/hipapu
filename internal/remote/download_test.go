// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package remote_test

import (
	"os"

	"github.com/outcatcher/hipapu/internal/remote"
)

func (gs *GithubSuite) TestDownloadFile() {
	ctx := gs.T().Context()

	remoteClient, err := remote.New(gs.token)
	gs.Require().NoError(err)

	release, err := remoteClient.GetLatestRelease(ctx, testOwner, testRepo)
	gs.Require().NoError(err)

	filePath := gs.T().TempDir() + "/" + release.Assets[0].Filename

	gs.T().Cleanup(func() {
		gs.Require().NoError(os.Remove(filePath))
	})

	file, err := os.Create(filePath)
	gs.Require().NoError(err)

	// closing file after that is important, so using 'assert' instead of 'require'
	gs.Assert().NoError(remoteClient.DownloadFile(ctx, release.Assets[0].DownloadURL, file))

	gs.Require().NoError(file.Close())

	data, err := os.ReadFile(filePath)
	gs.Require().NoError(err)

	gs.Require().Equal(testData, string(data))
}
