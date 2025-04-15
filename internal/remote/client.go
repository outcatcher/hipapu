// Copyright (C) 2025  Anton Kachurin
package remote

import (
	"net/http"
	"time"

	github "github.com/google/go-github/v70/github"
)

// Asset - release attachment.
type Asset struct {
	Filename    string
	DownloadURL string
	TotalSize   int
}

// Release - github release info.
type Release struct {
	RepoURL     string
	Owner       string
	Repo        string
	Name        string
	Description string
	PublishedAt time.Time
	Assets      []Asset
}

// Client - GH client w/o authentication.
type Client struct {
	client *github.Client
}

// New creates new GH client.
func New(token string) (*Client, error) {
	httpClient := &http.Client{} // HTTP client without timeout

	ghClient := github.NewClient(httpClient)

	if token != "" {
		ghClient = ghClient.WithAuthToken(token)
	}

	return &Client{
		client: ghClient,
	}, nil
}
