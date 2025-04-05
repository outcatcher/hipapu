package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/outcatcher/hipapu/internal/config"
	"github.com/outcatcher/hipapu/internal/local"
	"github.com/outcatcher/hipapu/internal/remote"
)

// ErrNoConfig - error for missing configuration path.
var ErrNoConfig = errors.New("no config provided")

// Application is a base application state.
//
// Serves as entry point to external consumers.
type Application struct {
	config cfg
	remote remoteClient
	files  localFiles
}

type cfg interface {
	// Add adds installation to the list.
	Add(installation config.Installation) error
	// GetInstallations returns tracked installs.
	GetInstallations() []config.Installation
}

type remoteClient interface {
	// GetLatestRelease - retrieves latest repository release.
	GetLatestRelease(ctx context.Context, owner, repo string) (*remote.Release, error)
	// DownloadFile downloads binary file.
	DownloadFile(ctx context.Context, downloadURL string, writer io.Writer) error
}

type localFiles interface {
	// GetFileInfo returns info on the local file.
	GetFileInfo(filePath string) (*local.FileInfo, error)
}

// New create new Application instance with given config and GH token from env.
//
// todo (?): move somewhere else.
func New(configPath string) (*Application, error) {
	app := new(Application)

	if configPath == "" {
		return nil, ErrNoConfig
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load app configuration: %w", err)
	}

	app.WithConfig(cfg)

	remote, err := remote.New(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("failed to create GH client: %w", err)
	}

	app.WithRemote(remote)
	app.WithFiles(new(local.LocalFiles))

	return app, nil
}

// WithConfig sets up configuration for the app.
func (a *Application) WithConfig(cfg cfg) {
	a.config = cfg
}

// WithRemote sets up remote for the app.
func (a *Application) WithRemote(remote remoteClient) {
	a.remote = remote
}

// WithRemote sets up file operations for the app.
func (a *Application) WithFiles(files localFiles) {
	a.files = files
}
