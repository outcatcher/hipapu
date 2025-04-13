// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package handlers //nolint:testpackage  // don't want to make useless tests overly complex

import (
	"bytes"
	"testing"

	"github.com/outcatcher/hipapu/cmd/handlers/mocks"
	"github.com/outcatcher/hipapu/internal/config"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	const (
		expectedFilePath = "test-file-path-add"
		expectedRepoPath = "test-repo-path-add"
	)

	expectedInstallations := []config.Installation{
		{
			RepoURL:   expectedRepoPath,
			LocalPath: expectedFilePath,
		},
	}

	appMock := mocks.NewMockapplication(t)
	appMock.On("List").Once().Return(expectedInstallations)

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
}
