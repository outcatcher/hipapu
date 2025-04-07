// Copyright (C) 2025  Anton Kachurin
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/outcatcher/hipapu/internal/config"
	"github.com/schollz/progressbar/v3"
)

// Public sync errors. Self-explanatory.
var (
	ErrEmptyInstallationList = errors.New("installation list is empty")
	ErrMissingAsset          = errors.New("no asset with given name found")
)

// Synchronize runs synchronization of all new releases replacing local files reporting the progress.
func (a *Application) Synchronize(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	//todo: move to app shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		cancel()
	}()

	installations := a.config.GetInstallations()

	if len(installations) == 0 {
		return ErrEmptyInstallationList
	}

	a.log().InfoContext(ctx, "found installactions", "count", len(installations))

	var errs error

	for i, installation := range installations {
		fmt.Printf("Synchronizing installation #%d of %d\n", i+1, len(installations))

		// todo: parrallelize
		if err := a.syncInstallation(ctx, installation); err != nil {
			errs = errors.Join(errs, err)

			continue
		}
	}

	return errs
}

//nolint:cyclop,funlen  // rewriting makes it less readable
func (a *Application) syncInstallation(
	ctx context.Context, installation config.Installation,
) error {
	file, err := a.files.GetFileInfo(installation.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	urlParts := strings.Split(installation.RepoURL, "/")
	owner, repo := urlParts[len(urlParts)-2], urlParts[len(urlParts)-1]

	a.log().Info("Starting sync of installation", "owner", owner, "repo", repo, "local path", installation.LocalPath)

	release, err := a.remote.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get release info: %w", err)
	}

	if !release.PublishedAt.After(file.LastModified) {
		a.log().Info(
			"Current installation is up to date",
			"published at", release.PublishedAt.Format(time.RFC3339),
			"last modified", file.LastModified.Format(time.RFC3339),
		)

		// nothing to do, local version seems to be newer
		return nil
	}

	downloadURL := ""
	totalSize := 0

	for _, asset := range release.Assets {
		if asset.Filename == file.Name {
			downloadURL = asset.DownloadURL
			totalSize = asset.TotalSize

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

	a.log().Info("Download started", "download URL", downloadURL, "total size, MiB", totalSize/1024/1024) //nolint:mnd

	bar := progressbar.DefaultBytes(-1, "Downloading to"+tmpFilePath)
	bar.ChangeMax(totalSize)

	// cleanup on cancel
	go func() {
		<-ctx.Done()

		_ = tmpFile.Close()
	}()

	writer := io.MultiWriter(bar, tmpFile)

	if err := a.remote.DownloadFile(ctx, downloadURL, writer); err != nil {
		return fmt.Errorf("failed to dowload to tmp file: %w", err)
	}

	a.log().Info("Download finished", "download URL", downloadURL)

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close tmp file: %w", err)
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old file: %w", err)
	}

	if err := os.Rename(tmpFilePath, file.FilePath); err != nil {
		return fmt.Errorf("failed to rename tmp file: %w", err)
	}

	a.log().Info("Finished sync of installation",
		"owner", owner,
		"repo", repo,
		"local path", installation.LocalPath,
		"download URL", downloadURL,
	)

	return nil
}
