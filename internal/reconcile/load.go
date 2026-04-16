package reconcile

import (
	"fmt"

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

func (p Plan) Messages() []string {
	var messages []string

	for _, repository := range p.githubRepositories {
		messages = append(messages,
			fmt.Sprintf("Reconciling GitHubRepository %s/%s", repository.Manifest.Spec.Owner, repository.Manifest.Spec.Name),
			"Dry run only: no external changes applied",
		)
	}

	return messages
}
