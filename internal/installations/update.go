package installations

import (
	"errors"
	"fmt"

	"golang.org/x/mod/semver"
)

// ErrUnsupportedVersion - unsupported lockfile version.
var ErrUnsupportedVersion = errors.New("unsupported lockfile version")

// UpdateVersion updates lockfile format making a backup.
func (l *Lock) UpdateVersion() error {
	switch semver.Compare(currentVersion, l.lockData.LockVersion) {
	case 0: // up to date
		return nil
	case -1: // newer than existing
		return fmt.Errorf("%w: %s, latest current version is %s",
			ErrUnsupportedVersion, l.lockData.LockVersion, currentVersion)
	}

	if err := l.saveInstallations(true); err != nil {
		return fmt.Errorf("failed to save installations: %w", err)
	}

	return nil
}
