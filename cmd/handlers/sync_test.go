// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package handlers //nolint:testpackage  // don't want to make useless tests overly complex

import (
	"io"
	"testing"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/cmd/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	appMock := mocks.NewMockapplication(t)
	appMock.On("Synchronize", ctx).Once().Return(nil)

	hdl := &ActionHandlers{app: appMock}

	syncCmd := hdl.CommandSync()
	syncCmd.Writer = io.Discard

	require.NoError(t, syncCmd.Action(ctx, syncCmd))
}

func TestSync_emptyList(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	appMock := mocks.NewMockapplication(t)
	appMock.On("Synchronize", ctx).Once().Return(app.ErrEmptyInstallationList)

	hdl := &ActionHandlers{app: appMock}

	syncCmd := hdl.CommandSync()
	syncCmd.Writer = io.Discard

	require.NoError(t, syncCmd.Action(ctx, syncCmd))
}
