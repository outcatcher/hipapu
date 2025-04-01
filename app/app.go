package app

import (
	"errors"
	"fmt"

	"github.com/outcatcher/hipapu/internal/config"
	"github.com/outcatcher/hipapu/internal/remote"
)

// ErrNoConfig - error for missing configuration path.
var ErrNoConfig = errors.New("no config provided")

// Application is a base application state.
//
// Serves as entry point to external consumers.
type Application struct {
	config *config.Config

	// todo: remake both remote and local with interfaces for smoother testing
	remote *remote.GitHubClient
}

// New create new Application instance with given config.
func New(configPath string) (*Application, error) {
	if configPath == "" {
		return nil, ErrNoConfig
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load app configuration: %w", err)
	}

	remote, err := remote.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create GH client: %w", err)
	}

	return &Application{
		config: cfg,
		remote: remote,
	}, nil
}
