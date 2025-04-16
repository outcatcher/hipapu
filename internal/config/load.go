// Copyright (C) 2025  Anton Kachurin
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

const selfRepo = "https://github.com/outcatcher/hipapu"

// LoadConfig loads configuration from file.
func LoadConfig(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get config abs path: %w", err)
	}

	file, err := os.Open(path) //nolint:gosec  // Abs already calling Clean for result
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open config file to read: %w", err)
	}

	configuration := &Config{
		filePath: path,
	}

	// if file exists, load config
	if file != nil {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				slog.Error("failed to close config file (read)", "error", closeErr)
			}
		}()

		if err := json.NewDecoder(file).Decode(configuration); err != nil {
			return nil, fmt.Errorf("failed to read config : %w", err)
		}
	}

	configuration.setDefaults()

	return configuration, nil
}

func (c *Config) setDefaults() {
	hasSelf := slices.ContainsFunc(c.Installations, func(inst Installation) bool {
		return inst.RepoURL == selfRepo
	})

	if !hasSelf {
		c.Installations = appendSelf(c.Installations)
	}
}

func appendSelf(installs []Installation) []Installation {
	cwd, err := os.Getwd()
	if err != nil {
		return installs
	}

	return append(installs, Installation{
		RepoURL:         selfRepo,
		LocalPath:       filepath.Join(cwd, os.Args[0]), // register self-update
		KeepLastVersion: true,                           // safety measure for broken releases
		SkipSync:        true,                           // todo: remove after KeepLastVersion implemented
	})
}
