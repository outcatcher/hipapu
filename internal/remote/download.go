// Copyright (C) 2025  Anton Kachurin
package remote

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var (
	errUnexpectedDownloadStatus = errors.New("unexpected download status")
)

// DownloadFile downloads binary file.
func (c *Client) DownloadFile(ctx context.Context, downloadURL string, writer io.Writer) error {
	// todo: add progress tracking
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("error preparing requiest: %w", err)
	}

	req.Header.Add("Accept", "application/octet-stream") // otherwise you get json response instead of file

	resp, err := c.client.Client().Do(req)
	if err != nil {
		return fmt.Errorf("error performing requiest: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w %s", errUnexpectedDownloadStatus, resp.Status)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close download response body", "error", closeErr)
		}
	}()

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}
