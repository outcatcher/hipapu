// Copyright (C) 2025  Anton Kachurin
package config

// Config - base app configuration.
type Config struct {
	// corresponding config file path
	filePath string

	Installations []Installation `json:"installations"` // todo: make concurrent-friendly
}

// Installation to be synced.
type Installation struct {
	RepoURL         string `json:"repo_url"`
	LocalPath       string `json:"local_path"`
	KeepLastVersion bool   `json:"keep_last_version,omitempty"`
	SkipSync        bool   `json:"skip_sync,omitempty"`
	// todo: add fields in future versions:
	// Type            InstallationType `json:"type"`
}

// GetInstallations returns tracked installs.
func (c *Config) GetInstallations() []Installation {
	return c.Installations
}
