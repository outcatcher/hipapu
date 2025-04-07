// Copyright (C) 2025  Anton Kachurin
package remote

import (
	"context"
	"fmt"
)

// GetLatestRelease - retrieves latest repository release.
func (c *Client) GetLatestRelease(ctx context.Context, owner, repo string) (*Release, error) {
	release, _, err := c.client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("error getting release: %w", err)
	}

	resultAssets := make([]Asset, 0, len(release.Assets))

	for _, asset := range release.Assets {
		if asset == nil {
			continue
		}

		resultAssets = append(resultAssets, Asset{
			Filename:    asset.GetName(),
			DownloadURL: asset.GetURL(),
			TotalSize:   asset.GetSize(),
		})
	}

	result := &Release{
		Name:        release.GetName(),
		Description: release.GetBody(),
		PublishedAt: release.PublishedAt.Time,
		Assets:      resultAssets,
	}

	return result, nil
}
