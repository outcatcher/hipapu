// Copyright (C) 2025  Anton Kachurin
package remote

import (
	"context"
	"fmt"
	"strings"
)

// GetLatestRelease - retrieves latest repository release.
func (c *Client) GetLatestRelease(ctx context.Context, repoURL string) (*Release, error) {
	urlParts := strings.Split(repoURL, "/")
	owner, repo := urlParts[len(urlParts)-2], urlParts[len(urlParts)-1]

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
		Owner:       owner,
		Repo:        repo,
		RepoURL:     repoURL,
	}

	return result, nil
}
