// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package handlers //nolint:testpackage  // don't want to make useless tests overly complex

import (
	"io"
	"testing"

	"github.com/outcatcher/hipapu/cmd/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	const (
		expectedFilePath = "test-file-path-add"
		expectedRepoPath = "test-repo-path-add"
	)

	appMock := mocks.NewMockapplication(t)
	appMock.On("Add", expectedRepoPath, expectedFilePath).Once().Return(nil)

	hdl := &ActionHandlers{
		filePath: expectedFilePath,
		repoPath: expectedRepoPath,
		lockPath: "test-lock-path-add",
		app:      appMock,
	}

	addCmd := hdl.CommandAdd()
	addCmd.Writer = io.Discard

	require.NoError(t, addCmd.Action(ctx, addCmd))
}
