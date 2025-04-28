// Copyright (C) 2025  Anton Kachurin
package installations

// Lock - installations lockfile manager.
type Lock struct {
	// corresponding lockfile path
	filePath string

	lockData lockData
}

// GetInstallations returns tracked installs.
func (l *Lock) GetInstallations() []Installation {
	return l.lockData.Installations
}

type lockData struct {
	LockVersion   string         `json:"version"`
	Installations []Installation `json:"installations"`
}

// Installation to be synced.
type Installation struct {
	ID              string `json:"id"`
	RepoURL         string `json:"repo_url"`
	LocalPath       string `json:"local_path"`
	KeepLastVersion bool   `json:"keep_last_version,omitempty"`
	SkipSync        bool   `json:"skip_sync,omitempty"`
	// todo: add fields in future versions:
	// Type            InstallationType `json:"type"`
}
