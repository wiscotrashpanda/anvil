package reconcile

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	manifestv1alpha1 "github.com/emkaytec/alloy/manifest/v1alpha1"
	ghapi "github.com/emkaytec/anvil/internal/github"
	"github.com/emkaytec/anvil/internal/manifest"
)

func reconcileGitHubRepository(ctx context.Context, client *ghapi.Client, loaded manifest.LoadedGitHubRepositoryManifest) ([]string, error) {
	spec := loaded.Manifest.Spec
	repositoryID := fmt.Sprintf("%s/%s", spec.Owner, spec.Name)

	messages := []string{fmt.Sprintf("Reconciling GitHubRepository %s", repositoryID)}

	repository, created, err := ensureRepository(ctx, client, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile %s: %w", repositoryID, err)
	}

	changed := created
	if created {
		messages = append(messages, fmt.Sprintf("Created GitHub repository %s", repositoryID))
	}

	repository, updatedRepositorySettings, err := reconcileRepositorySettings(ctx, client, repository, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile repository settings for %s: %w", repositoryID, err)
	}
	if updatedRepositorySettings {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated repository settings for %s", repositoryID))
	}

	updatedTopics, err := reconcileTopics(ctx, client, repository, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile topics for %s: %w", repositoryID, err)
	}
	if updatedTopics {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated repository topics for %s", repositoryID))
	}

	if !changed {
		messages = append(messages, fmt.Sprintf("GitHubRepository %s is up to date", repositoryID))
	}

	return messages, nil
}

func ensureRepository(ctx context.Context, client *ghapi.Client, spec manifestv1alpha1.GitHubRepositorySpec) (*ghapi.Repository, bool, error) {
	repository, err := client.GetRepository(ctx, spec.Owner, spec.Name)
	if err == nil {
		return repository, false, nil
	}
	if !ghapi.IsNotFound(err) {
		return nil, false, err
	}

	account, err := client.GetAccount(ctx, spec.Owner)
	if err != nil {
		return nil, false, fmt.Errorf("detect owner type: %w", err)
	}

	request := ghapi.CreateRepositoryRequest{
		Name:             spec.Name,
		Visibility:       spec.Visibility,
		Description:      spec.Description,
		Homepage:         spec.Homepage,
		AutoInit:         spec.AutoInit,
		HasIssues:        featuresValue(spec.Features, func(features *manifestv1alpha1.GitHubRepositoryFeaturesSpec) *bool { return features.HasIssues }),
		HasProjects:      featuresValue(spec.Features, func(features *manifestv1alpha1.GitHubRepositoryFeaturesSpec) *bool { return features.HasProjects }),
		HasWiki:          featuresValue(spec.Features, func(features *manifestv1alpha1.GitHubRepositoryFeaturesSpec) *bool { return features.HasWiki }),
		AllowSquashMerge: mergePolicyValue(spec.MergePolicy, func(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool { return policy.AllowSquashMerge }),
		AllowMergeCommit: mergePolicyValue(spec.MergePolicy, func(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool { return policy.AllowMergeCommit }),
		AllowRebaseMerge: mergePolicyValue(spec.MergePolicy, func(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool { return policy.AllowRebaseMerge }),
		AllowAutoMerge:   mergePolicyValue(spec.MergePolicy, func(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool { return policy.AllowAutoMerge }),
		DeleteBranchOnMerge: mergePolicyValue(spec.MergePolicy, func(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool {
			return policy.DeleteBranchOnMerge
		}),
	}

	if account.Type == "Organization" {
		repository, err = client.CreateOrganizationRepository(ctx, spec.Owner, request)
	} else {
		repository, err = client.CreateUserRepository(ctx, request)
	}
	if err != nil {
		return nil, false, err
	}

	return repository, true, nil
}

func reconcileRepositorySettings(ctx context.Context, client *ghapi.Client, repository *ghapi.Repository, spec manifestv1alpha1.GitHubRepositorySpec) (*ghapi.Repository, bool, error) {
	request := ghapi.UpdateRepositoryRequest{}

	if spec.Visibility != nil && repository.Visibility != dereferenceString(spec.Visibility) {
		request.Visibility = spec.Visibility
	}

	if spec.Description != nil && dereferenceString(repository.Description) != dereferenceString(spec.Description) {
		request.Description = spec.Description
	}

	if spec.Homepage != nil && repository.Homepage != dereferenceString(spec.Homepage) {
		request.Homepage = spec.Homepage
	}

	if spec.DefaultBranch != nil && repository.DefaultBranch != dereferenceString(spec.DefaultBranch) {
		request.DefaultBranch = spec.DefaultBranch
	}

	if spec.Features != nil {
		if spec.Features.HasIssues != nil && repository.HasIssues != *spec.Features.HasIssues {
			request.HasIssues = spec.Features.HasIssues
		}
		if spec.Features.HasProjects != nil && repository.HasProjects != *spec.Features.HasProjects {
			request.HasProjects = spec.Features.HasProjects
		}
		if spec.Features.HasWiki != nil && repository.HasWiki != *spec.Features.HasWiki {
			request.HasWiki = spec.Features.HasWiki
		}
	}

	if spec.MergePolicy != nil {
		if spec.MergePolicy.AllowSquashMerge != nil && repository.AllowSquashMerge != *spec.MergePolicy.AllowSquashMerge {
			request.AllowSquashMerge = spec.MergePolicy.AllowSquashMerge
		}
		if spec.MergePolicy.AllowMergeCommit != nil && repository.AllowMergeCommit != *spec.MergePolicy.AllowMergeCommit {
			request.AllowMergeCommit = spec.MergePolicy.AllowMergeCommit
		}
		if spec.MergePolicy.AllowRebaseMerge != nil && repository.AllowRebaseMerge != *spec.MergePolicy.AllowRebaseMerge {
			request.AllowRebaseMerge = spec.MergePolicy.AllowRebaseMerge
		}
		if spec.MergePolicy.AllowAutoMerge != nil && repository.AllowAutoMerge != *spec.MergePolicy.AllowAutoMerge {
			request.AllowAutoMerge = spec.MergePolicy.AllowAutoMerge
		}
		if spec.MergePolicy.AllowUpdateBranch != nil && repository.AllowUpdateBranch != *spec.MergePolicy.AllowUpdateBranch {
			request.AllowUpdateBranch = spec.MergePolicy.AllowUpdateBranch
		}
		if spec.MergePolicy.DeleteBranchOnMerge != nil && repository.DeleteBranchOnMerge != *spec.MergePolicy.DeleteBranchOnMerge {
			request.DeleteBranchOnMerge = spec.MergePolicy.DeleteBranchOnMerge
		}
	}

	if request.IsZero() {
		return repository, false, nil
	}

	updatedRepository, err := client.UpdateRepository(ctx, spec.Owner, spec.Name, request)
	if err != nil {
		return nil, false, err
	}

	return updatedRepository, true, nil
}

func reconcileTopics(ctx context.Context, client *ghapi.Client, repository *ghapi.Repository, spec manifestv1alpha1.GitHubRepositorySpec) (bool, error) {
	if spec.Topics == nil {
		return false, nil
	}

	desiredTopics := normalizeStrings(spec.Topics)
	currentTopics := normalizeStrings(repository.Topics)
	if reflect.DeepEqual(desiredTopics, currentTopics) {
		return false, nil
	}

	if err := client.ReplaceTopics(ctx, spec.Owner, spec.Name, desiredTopics); err != nil {
		return false, err
	}

	return true, nil
}

func featuresValue(spec *manifestv1alpha1.GitHubRepositoryFeaturesSpec, selector func(*manifestv1alpha1.GitHubRepositoryFeaturesSpec) *bool) *bool {
	if spec == nil {
		return nil
	}

	return selector(spec)
}

func mergePolicyValue(spec *manifestv1alpha1.GitHubRepositoryMergePolicySpec, selector func(*manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool) *bool {
	if spec == nil {
		return nil
	}

	return selector(spec)
}

func dereferenceString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func normalizeStrings(values []string) []string {
	if values == nil {
		return nil
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		normalized = append(normalized, strings.ToLower(strings.TrimSpace(value)))
	}
	sort.Strings(normalized)
	return normalized
}
