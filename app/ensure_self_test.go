// Copyright (C) 2025  Anton Kachurin

//go:build test
// +build test

package app_test

import (
	"testing"

	"github.com/outcatcher/hipapu/app"
	"github.com/outcatcher/hipapu/app/mocks"
	"github.com/outcatcher/hipapu/internal/installations"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const selfRepo = "https://github.com/outcatcher/hipapu"

func TestEnsureSelf(t *testing.T) {
	t.Parallel()

	lockMock := mocks.NewMockinstallationsLock(t)
	lockMock.On("GetInstallations").Return(nil)
	lockMock.
		On("Add", mock.AnythingOfType("installations.Installation")).
		Run(func(args mock.Arguments) {
			inst, ok := args.Get(0).(installations.Installation)
			require.True(t, ok)

			require.Equal(t, inst.ID, "0")
			require.Equal(t, inst.RepoURL, selfRepo)
		}).
		Return(nil)

	app := new(app.Application)
	app.WithLock(lockMock)

	require.NoError(t, app.EnsureSelf())
}

func TestEnsureSelf_existing(t *testing.T) {
	t.Parallel()

	lockMock := mocks.NewMockinstallationsLock(t)
	lockMock.On("GetInstallations").Return([]installations.Installation{
		{
			RepoURL: selfRepo,
		},
	})

	app := new(app.Application)
	app.WithLock(lockMock)

	require.NoError(t, app.EnsureSelf())
}
