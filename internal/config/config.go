package config

// Config - base app configuration.
type Config struct {
	// corresponding config file path
	filePath string

	DownloadLocation string         `json:"download_location"`
	Installations    []Installation `json:"installations"` // todo: make concurrent-friendly
}

// Installation to be synced.
type Installation struct {
	RepoURL   string `json:"repo_url"`
	LocalPath string `json:"local_path"`
	// todo: add fields in future versions:
	// Type            InstallationType `json:"type"`
	// OldVersionCount uint8            `json:"old_version_count"`
	// SkipSync        bool             `json:"skip_sync"`
}
