package github

import (
	"context"
	"fmt"

	githubapi "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"

	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

type Connector struct {
	client *githubapi.Client
	config connectors.Config
}

func New() *Connector {
	return &Connector{}
}

func (c *Connector) Name() string {
	return "github"
}

func (c *Connector) Connect(ctx context.Context, config connectors.Config) error {
	token := config["token"]

	if token == "" {
		return fmt.Errorf("github token is required")
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	c.client = githubapi.NewClient(
		oauth2.NewClient(ctx, tokenSource),
	)

	c.config = config

	return nil
}

func (c *Connector) Disconnect(ctx context.Context) error {
	c.client = nil
	c.config = nil

	return nil
}

func (c *Connector) ListResources(ctx context.Context) ([]connectors.Resource, error) {
	if c.client == nil {
		return nil, fmt.Errorf("github connector is not connected")
	}

	organization := c.config["organization"]

	var repositories []*githubapi.Repository
	var err error

	if organization != "" {
		repositories, _, err = c.client.Repositories.ListByOrg(
			ctx,
			organization,
			&githubapi.RepositoryListByOrgOptions{
				ListOptions: githubapi.ListOptions{
					PerPage: 100,
				},
			},
		)
	} else {
		repositories, _, err = c.client.Repositories.List(
			ctx,
			"",
			&githubapi.RepositoryListOptions{
				ListOptions: githubapi.ListOptions{
					PerPage: 100,
				},
			},
		)
	}

	if err != nil {
		return nil, fmt.Errorf("list github repositories: %w", err)
	}

	resources := make([]connectors.Resource, 0, len(repositories))

	for _, repository := range repositories {
		if repository == nil {
			continue
		}

		metadata := map[string]string{}

		if repository.HTMLURL != nil {
			metadata["url"] = repository.GetHTMLURL()
		}

		if repository.Language != nil {
			metadata["language"] = repository.GetLanguage()
		}

		resources = append(resources, connectors.Resource{
			ID:       repository.GetFullName(),
			Name:     repository.GetName(),
			Type:     "github_repository",
			Metadata: metadata,
		})
	}

	return resources, nil
}
