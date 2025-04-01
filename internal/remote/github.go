package remote

import (
	"context"
	"fmt"
	"net/http"
	"time"

	github "github.com/google/go-github/v70/github"
)

var ghClientTimeout = 10 * time.Second

// GitHubClient - GH client w/o authentication.
type GitHubClient struct {
	client *github.Client
}

// New creates new GH client.
func New() (*GitHubClient, error) {
	httpClient := &http.Client{
		Timeout: ghClientTimeout,
	}

	return &GitHubClient{
		client: github.NewClient(httpClient),
	}, nil
}

// Asset - release attachment.
type Asset struct {
	Filename    string
	DownloadURL string
}

// Release - github release info.
type Release struct {
	Name        string
	Description string
	PublishedAt time.Time
	Assets      []Asset
}

// GetLatestRelease - retrieves latest repository release.
func (c *GitHubClient) GetLatestRelease(ctx context.Context, owner, repo string) (*Release, error) {
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
			Filename:    *asset.Name,
			DownloadURL: *asset.BrowserDownloadURL,
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
