// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package handlers //nolint:testpackage  // don't want to make useless tests overly complex

import (
	"bytes"
	"testing"
	"time"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/cmd/handlers/mocks"
	"github.com/outcatcher/hipapu/internal/local"
	"github.com/outcatcher/hipapu/internal/remote"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	const (
		expectedFilePath = "test-file-path-add"
		expectedRepoPath = "test-repo-path-add"
	)

	expectedInstallations := []app.Installation{
		{
			Release: &remote.Release{
				PublishedAt: time.Now(),
				RepoURL:     expectedRepoPath,
			},
			LocalFile: &local.FileInfo{
				FilePath:     expectedFilePath,
				LastModified: time.Now().Add(-1 * time.Second),
			},
		},
	}

	appMock := mocks.NewMockapplication(t)
	appMock.On("List", ctx).Once().Return(expectedInstallations, nil)

	hdl := &ActionHandlers{app: appMock}

	addCmd := hdl.CommandList()

	testBuffer := new(bytes.Buffer)

	addCmd.Writer = testBuffer

	require.NoError(t, addCmd.Action(ctx, addCmd))

	line, err := testBuffer.ReadString('\n')
	require.NoError(t, err)

	require.Equal(t, "Installations:\n", line)

	line2, err := testBuffer.ReadString('\n')
	require.NoError(t, err)

	require.Contains(t, line2, expectedRepoPath)
	require.Contains(t, line2, expectedFilePath)
	require.Contains(t, line2, "HAS UPDATE")
}
