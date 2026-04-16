package reconcile

import (
	"context"

	ghapi "github.com/emkaytec/anvil/internal/github"
	"github.com/emkaytec/anvil/internal/manifest"
)

type Plan struct {
	githubRepositories []manifest.LoadedGitHubRepositoryManifest
}

func Load(path string) (Plan, error) {
	loaded, err := manifest.LoadDir(path)
	if err != nil {
		return Plan{}, err
	}

	return Plan{
		githubRepositories: loaded.GitHubRepositories,
	}, nil
}

func Run(ctx context.Context, path string) ([]string, error) {
	plan, err := Load(path)
	if err != nil {
		return nil, err
	}

	client, err := NewGitHubClient()
	if err != nil {
		return nil, err
	}

	return plan.Run(ctx, client)
}

func NewGitHubClient() (*ghapi.Client, error) {
	return ghapi.NewClientFromEnv()
}

func (p Plan) Run(ctx context.Context, client *ghapi.Client) ([]string, error) {
	var messages []string

	for _, repository := range p.githubRepositories {
		repositoryMessages, err := reconcileGitHubRepository(ctx, client, repository)
		messages = append(messages, repositoryMessages...)
		if err != nil {
			return messages, err
		}
	}

	return messages, nil
}
