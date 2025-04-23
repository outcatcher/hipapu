package lock

import (
	"errors"
	"fmt"

	"golang.org/x/mod/semver"
)

var ErrUnsupportedVersion = errors.New("unsupported lockfile verison")

// UpdateVersion updates lockfile format making a backup.
func (l *Lock) UpdateVersion() error {
	if err := l.LoadInstallations(l.filePath); err != nil {
		return fmt.Errorf("failed to load installations: %w", err)
	}

	switch semver.Compare(currentVersion, l.lockData.LockVersion) {
	case 0:
		return nil
	case -1:
		return fmt.Errorf("%w: %s, latest current version is %s",
			ErrUnsupportedVersion, l.lockData.LockVersion, currentVersion)
	}

	if err := l.saveInstallations(true); err != nil {
		return fmt.Errorf("failed to save installations: %w", err)
	}

	return nil
}
