// Copyright (C) 2025  Anton Kachurin
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Public sync errors. Self-explanatory.
var (
	ErrEmptyInstallationList = errors.New("installation list is empty")
	ErrMissingAsset          = errors.New("no asset with given name found")
)

// Synchronize runs synchronization of all new releases replacing local files reporting the progress.
func (a *Application) Synchronize(ctx context.Context, writer io.Writer) error {
	installations, err := a.List(ctx)
	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}

	if len(installations) == 0 {
		return ErrEmptyInstallationList
	}

	a.log().InfoContext(ctx, "found installactions", "count", len(installations))

	var errs error

	for i, installation := range installations {
		if installation.SkipSync {
			_, _ = fmt.Fprintf(
				writer,
				"Skipping sync for installation #%d of %d (%s)\n",
				i+1,
				len(installations),
				installation.Release.RepoURL,
			)

			continue
		}

		// todo: output is a CLI interaction, needs to be moved out to cmd somehow
		_, _ = fmt.Fprintf(
			writer,
			"Synchronizing installation #%d of %d (%s)\n",
			i+1,
			len(installations),
			installation.Release.RepoURL,
		)

		// todo: parrallelize
		if err := a.syncInstallation(ctx, installation, writer); err != nil {
			errs = errors.Join(errs, err)

			continue
		}
	}

	return errs
}

//nolint:cyclop,funlen  // rewriting makes it less readable
func (a *Application) syncInstallation(
	ctx context.Context, installation Installation, extWriter io.Writer,
) error {
	a.log().InfoContext(ctx,
		"Starting sync of installation",
		"owner", installation.Release.Owner,
		"repo", installation.Release.Repo,
		"local path", installation.LocalFile.FilePath,
	)

	release, file := installation.Release, installation.LocalFile

	if !release.PublishedAt.After(installation.LocalFile.LastModified) {
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

	// todo: progress bar is a CLI interaction, needs to be moved out to cmd somehow
	bar := newProgressBar(extWriter, totalSize, "Downloading to "+tmpFilePath)

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
		"owner", release.Owner,
		"repo", release.Repo,
		"local path", file.FilePath,
		"download URL", downloadURL,
	)

	return nil
}

func newProgressBar(writer io.Writer, maxBytes int, description string) *progressbar.ProgressBar {
	//nolint:mnd  // temporary location for progress bar initialization
	return progressbar.NewOptions(
		maxBytes,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowTotalBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(writer, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
}
