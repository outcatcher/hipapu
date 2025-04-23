// Copyright (C) 2025  Anton Kachurin
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/outcatcher/hipapu/internal/local"
	"github.com/outcatcher/hipapu/internal/lock"
	"github.com/outcatcher/hipapu/internal/remote"
)

// ErrNoLock - error for missing lockfile path.
var ErrNoLock = errors.New("no lock provided")

// Application is a base application state.
//
// Serves as entry point to external consumers.
type Application struct {
	lockfile installationsLock
	remote   remoteClient
	files    localFiles

	logger *slog.Logger
}

// New create new Application instance with given lock and GH token from env.
//
// todo (?): move somewhere else.
func New(lockPath string) (*Application, error) {
	app := new(Application)

	if lockPath == "" {
		return nil, ErrNoLock
	}

	app.logger = initLogger()

	locks := new(lock.Lock)

	err := locks.LoadInstallations(lockPath)
	if err != nil {
		app.logger.Error("lockfile missing or corrupted (%s)")

		return nil, fmt.Errorf("failed to load installations: %w", err)
	}

	app.WithLockfile(locks)

	remote, err := remote.New(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("failed to create GH client: %w", err)
	}

	app.WithRemote(remote)
	app.WithFiles(new(local.Files))

	return app, nil
}

// todo: rewrite
func initLogger() *slog.Logger {
	logFile, err := os.Create("hipapu.log")
	if err != nil {
		return nil
	}

	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})

	logger := slog.New(fileHandler)

	logger.Info("Log started")

	return logger
}

type installationsLock interface {
	// Add adds installation to the list.
	Add(installation lock.Installation) error
	// GetInstallations returns tracked installs.
	GetInstallations() []lock.Installation
}

// WithLockfile sets up lockfile for the app.
func (a *Application) WithLockfile(locks installationsLock) {
	a.lockfile = locks
}

type remoteClient interface {
	// GetLatestRelease - retrieves latest repository release.
	GetLatestRelease(ctx context.Context, repoURL string) (*remote.Release, error)
	// DownloadFile downloads binary file.
	DownloadFile(ctx context.Context, downloadURL string, writer io.Writer) error
}

// WithRemote sets up remote for the app.
func (a *Application) WithRemote(remote remoteClient) {
	a.remote = remote
}

type localFiles interface {
	// GetFileInfo returns info on the local file.
	GetFileInfo(filePath string) (*local.FileInfo, error)
}

// WithFiles sets up file operations for the app.
func (a *Application) WithFiles(files localFiles) {
	a.files = files
}

func (a *Application) log() *slog.Logger {
	if a.logger == nil {
		a.logger = slog.New(slog.DiscardHandler)
	}

	return a.logger
}
