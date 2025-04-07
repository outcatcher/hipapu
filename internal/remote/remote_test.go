// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package remote_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	github "github.com/google/go-github/v70/github"
	"github.com/stretchr/testify/suite"
)

type GithubSuite struct {
	suite.Suite

	token string

	client *github.Client

	assetID   int64
	releaseID int64
}

func (gs *GithubSuite) SetupSuite() {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	gs.Require().True(ok)

	gs.token = token

	gs.client = github.NewClient(http.DefaultClient).WithAuthToken(gs.token)

	ctx := gs.T().Context()

	latest := "true"

	release, _, err := gs.client.Repositories.CreateRelease(
		ctx,
		testOwner,
		testRepo,
		&github.RepositoryRelease{
			Name:       &testReleaseName,
			Body:       &testReleaseBody,
			MakeLatest: &latest,
			TagName:    &testReleaseName,
		},
	)
	gs.Require().NoError(err)

	filePath := gs.T().TempDir() + "/" + testAssetFileName

	file, err := os.Create(filePath)
	gs.Require().NoError(err)

	_, err = io.WriteString(file, testData)
	gs.Require().NoError(err)

	gs.Require().NoError(file.Close())

	readFile, err := os.Open(filePath)
	gs.Require().NoError(err)

	asset, _, err := gs.client.Repositories.UploadReleaseAsset(
		ctx,
		testOwner,
		testRepo,
		release.GetID(),
		&github.UploadOptions{
			Name:  testAssetFileName,
			Label: "Example file",
		},
		readFile,
	)
	gs.Require().NoError(err)

	gs.assetID = asset.GetID()
	gs.releaseID = release.GetID()

	gs.T().Log("Created test release")
}

func (gs *GithubSuite) TearDownSuite() {
	_, err := gs.client.Repositories.DeleteReleaseAsset(gs.T().Context(), testOwner, testRepo, gs.assetID)
	gs.Require().NoError(err)

	_, err = gs.client.Repositories.DeleteRelease(gs.T().Context(), testOwner, testRepo, gs.releaseID)
	gs.Require().NoError(err)

	gs.T().Log("Removed test release")
}

func TestWorkflow(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GithubSuite))
}
