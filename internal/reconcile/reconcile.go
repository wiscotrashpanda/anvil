package reconcile

import (
	"context"
	"fmt"

	ghapi "github.com/emkaytec/anvil/internal/github"
	hcpapi "github.com/emkaytec/anvil/internal/hcpterraform"
	"github.com/emkaytec/anvil/internal/manifest"
)

type Plan struct {
	githubRepositories     []manifest.LoadedGitHubRepositoryManifest
	hcpTerraformWorkspaces []manifest.LoadedHCPTerraformWorkspaceManifest
}

func Load(path string) (Plan, error) {
	loaded, err := manifest.LoadDir(path)
	if err != nil {
		return Plan{}, err
	}

	return Plan{
		githubRepositories:     loaded.GitHubRepositories,
		hcpTerraformWorkspaces: loaded.HCPTerraformWorkspaces,
	}, nil
}

func Run(ctx context.Context, path string) ([]string, error) {
	plan, err := Load(path)
	if err != nil {
		return nil, err
	}

	var githubClient *ghapi.Client
	if len(plan.githubRepositories) > 0 {
		var err error
		githubClient, err = NewGitHubClient()
		if err != nil {
			return nil, err
		}
	}

	var hcpTerraformClient *hcpapi.Client
	if len(plan.hcpTerraformWorkspaces) > 0 {
		var err error
		hcpTerraformClient, err = NewHCPTerraformClient()
		if err != nil {
			return nil, err
		}
	}

	return plan.Run(ctx, githubClient, hcpTerraformClient)
}

func NewGitHubClient() (*ghapi.Client, error) {
	return ghapi.NewClientFromEnv()
}

func NewHCPTerraformClient() (*hcpapi.Client, error) {
	return hcpapi.NewClientFromEnv()
}

func (p Plan) HasGitHubRepositories() bool {
	return len(p.githubRepositories) > 0
}

func (p Plan) HasHCPTerraformWorkspaces() bool {
	return len(p.hcpTerraformWorkspaces) > 0
}

func (p Plan) Run(ctx context.Context, githubClient *ghapi.Client, hcpTerraformClient *hcpapi.Client) ([]string, error) {
	var messages []string

	for _, repository := range p.githubRepositories {
		if githubClient == nil {
			return messages, fmt.Errorf("missing GitHub client for GitHubRepository reconciliation")
		}

		repositoryMessages, err := reconcileGitHubRepository(ctx, githubClient, repository)
		messages = append(messages, repositoryMessages...)
		if err != nil {
			return messages, err
		}
	}

	for _, workspace := range p.hcpTerraformWorkspaces {
		if hcpTerraformClient == nil {
			return messages, fmt.Errorf("missing HCP Terraform client for HCPTerraformWorkspace reconciliation")
		}

		workspaceMessages, err := reconcileHCPTerraformWorkspace(ctx, hcpTerraformClient, workspace)
		messages = append(messages, workspaceMessages...)
		if err != nil {
			return messages, err
		}
	}

	return messages, nil
}
