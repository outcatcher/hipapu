//go:build test
// +build test

// Copyright (C) 2025  Anton Kachurin

package remote_test

import (
	"github.com/outcatcher/hipapu/internal/remote"
)

const (
	testOwner   = "outcatcher"
	testRepo    = "hipapu-release-testing"
	testRepoURL = "https://github.com/outcatcher/hipapu-release-testing"

	testData          = "6630cfdb-2123-4802-8ce2-f50336d56c18"
	testAssetFileName = "release.txt"
)

var (
	testReleaseName = "TestRelease"
	testReleaseBody = "Release for Testing"
)

func (gs *GithubSuite) TestGetLatestRelease() {
	client, err := remote.New(gs.token)
	gs.Require().NoError(err)

	release, err := client.GetLatestRelease(gs.T().Context(), testRepoURL)
	gs.Require().NoError(err)

	gs.Require().NotNil(release)

	gs.Require().Equal(testReleaseName, release.Name)
	gs.Require().Equal(testOwner, release.Owner)
	gs.Require().Equal(testRepo, release.Repo)
	gs.Require().Equal(testRepoURL, release.RepoURL)
	gs.Require().Equal(testReleaseBody, release.Description)

	gs.Require().Len(release.Assets, 1) // 2 default zips not included

	gs.Require().Equal(testAssetFileName, release.Assets[0].Filename)
}
