// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package app_test

import (
	"errors"
	"testing"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/app/mocks"
	"github.com/stretchr/testify/require"
)

func TestUpdateLockfileVersion(t *testing.T) {
	t.Parallel()

	filePath := t.TempDir() + "/test.lock"

	newApp, err := app.New(filePath)
	require.NoError(t, err)

	mockLock := mocks.NewMockinstallationsLock(t)
	mockLock.On("UpdateVersion").Return(nil)

	newApp.WithLock(mockLock)

	require.NoError(t, newApp.UpdateLockfileVersion())
}

func TestUpdateLockfileVersion_err(t *testing.T) {
	t.Parallel()

	filePath := t.TempDir() + "/test.lock"

	newApp, err := app.New(filePath)
	require.NoError(t, err)

	testErr := errors.New("123")

	mockLock := mocks.NewMockinstallationsLock(t)
	mockLock.On("UpdateVersion").Return(testErr)

	newApp.WithLock(mockLock)

	require.ErrorIs(t, newApp.UpdateLockfileVersion(), testErr)
}
