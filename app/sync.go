package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/outcatcher/hipapu/internal/config"
)

// Public sync errors. Self-explanatory.
var (
	ErrEmptyInstallationList = errors.New("installation list is empty")
	ErrMissingAsset          = errors.New("no asset with given name found")
)

// Synchronize downloads all new releases replacing local files.
func (a *Application) Synchronize(ctx context.Context) error {
	if len(a.config.GetInstallations()) == 0 {
		return ErrEmptyInstallationList
	}

	var errs error

	for _, installation := range a.config.GetInstallations() {
		// todo: parrallelize
		if err := a.syncInstallation(ctx, installation); err != nil {
			errs = errors.Join(errs, err)

			continue
		}
	}

	return errs
}

//nolint:cyclop  // rewriting makes it less readable
func (a *Application) syncInstallation(ctx context.Context, installation config.Installation) error {
	file, err := a.files.GetFileInfo(installation.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	urlParts := strings.Split(installation.RepoURL, "/")
	owner, repo := urlParts[len(urlParts)-2], urlParts[len(urlParts)-1]

	release, err := a.remote.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get release info: %w", err)
	}

	if !release.PublishedAt.After(file.LastModified) {
		// nothing to do, local version seems to be newer
		return nil
	}

	downloadURL := ""

	for _, asset := range release.Assets {
		if asset.Filename == file.Name {
			downloadURL = asset.DownloadURL

			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("%w: %s", ErrMissingAsset, file.Name)
	}

	tmpFilePath := file.FilePath + ".download"

	tmpFile, err := os.Create(filepath.Clean(tmpFilePath))
	if err != nil {
		return fmt.Errorf("failed to create tmp file: %w", err)
	}

	if err := a.remote.DownloadFile(ctx, downloadURL, tmpFile); err != nil {
		return fmt.Errorf("failed to dowload to tmp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close tmp file: %w", err)
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old file: %w", err)
	}

	if err := os.Rename(tmpFilePath, file.FilePath); err != nil {
		return fmt.Errorf("failed to rename tmp file: %w", err)
	}

	return nil
}
