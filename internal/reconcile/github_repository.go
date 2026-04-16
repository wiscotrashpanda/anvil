package reconcile

import (
	"context"
	"encoding/json"
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

	updatedPages, err := reconcilePages(ctx, client, repository, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile pages for %s: %w", repositoryID, err)
	}
	if updatedPages {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated GitHub Pages settings for %s", repositoryID))
	}

	updatedProperties, err := reconcileCustomProperties(ctx, client, repository, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile custom properties for %s: %w", repositoryID, err)
	}
	if updatedProperties {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated custom properties for %s", repositoryID))
	}

	branchMessages, branchChanged, err := reconcileBranches(ctx, client, repository, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile branches for %s: %w", repositoryID, err)
	}
	if branchChanged {
		changed = true
	}
	messages = append(messages, branchMessages...)

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

	if spec.Initialization != nil {
		request.GitignoreTemplate = spec.Initialization.GitignoreTemplate
		request.LicenseTemplate = spec.Initialization.LicenseTemplate
		request.IsTemplate = spec.Initialization.IsTemplate
	}

	if strings.EqualFold(account.Type, "Organization") {
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

	if spec.Visibility != nil && !strings.EqualFold(repository.Visibility, dereferenceString(spec.Visibility)) {
		request.Visibility = spec.Visibility
	}

	if spec.Description != nil && dereferenceString(spec.Description) != dereferenceString(repository.Description) {
		request.Description = spec.Description
	}

	if spec.Homepage != nil && dereferenceString(spec.Homepage) != repository.Homepage {
		request.Homepage = spec.Homepage
	}

	if spec.Archived != nil && repository.Archived != *spec.Archived {
		request.Archived = spec.Archived
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
		if spec.MergePolicy.SquashMergeCommitTitle != nil && repository.SquashMergeCommitTitle != dereferenceString(spec.MergePolicy.SquashMergeCommitTitle) {
			request.SquashMergeCommitTitle = spec.MergePolicy.SquashMergeCommitTitle
		}
		if spec.MergePolicy.SquashMergeCommitMessage != nil && repository.SquashMergeCommitMessage != dereferenceString(spec.MergePolicy.SquashMergeCommitMessage) {
			request.SquashMergeCommitMessage = spec.MergePolicy.SquashMergeCommitMessage
		}
		if spec.MergePolicy.MergeCommitTitle != nil && repository.MergeCommitTitle != dereferenceString(spec.MergePolicy.MergeCommitTitle) {
			request.MergeCommitTitle = spec.MergePolicy.MergeCommitTitle
		}
		if spec.MergePolicy.MergeCommitMessage != nil && repository.MergeCommitMessage != dereferenceString(spec.MergePolicy.MergeCommitMessage) {
			request.MergeCommitMessage = spec.MergePolicy.MergeCommitMessage
		}
	}

	if securityRequest := buildSecurityAndAnalysisUpdate(repository.SecurityAndAnalysis, spec.SecurityAndAnalysis); securityRequest != nil {
		request.SecurityAndAnalysis = securityRequest
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

func reconcilePages(ctx context.Context, client *ghapi.Client, repository *ghapi.Repository, spec manifestv1alpha1.GitHubRepositorySpec) (bool, error) {
	if spec.Pages == nil {
		return false, nil
	}

	currentPages, err := client.GetPages(ctx, spec.Owner, spec.Name)
	if err != nil && !ghapi.IsNotFound(err) {
		return false, err
	}

	desiredState := overlayPagesState(currentPagesState(currentPages), spec.Pages)
	if reflect.DeepEqual(desiredState, currentPagesState(currentPages)) {
		return false, nil
	}

	request := desiredState.toRequest()
	if currentPages == nil {
		if err := client.CreatePages(ctx, spec.Owner, spec.Name, request); err != nil {
			return false, err
		}
		return true, nil
	}

	if err := client.UpdatePages(ctx, spec.Owner, spec.Name, request); err != nil {
		return false, err
	}

	return true, nil
}

func reconcileCustomProperties(ctx context.Context, client *ghapi.Client, repository *ghapi.Repository, spec manifestv1alpha1.GitHubRepositorySpec) (bool, error) {
	if spec.CustomProperties == nil {
		return false, nil
	}

	currentValues, err := client.GetCustomPropertyValues(ctx, spec.Owner, spec.Name)
	if err != nil {
		return false, err
	}

	type currentProperty struct {
		Name  string
		Value any
	}

	currentByName := make(map[string]currentProperty, len(currentValues))
	for _, property := range currentValues {
		currentByName[strings.ToLower(property.PropertyName)] = currentProperty{
			Name:  property.PropertyName,
			Value: property.Value,
		}
	}

	desiredByName := make(map[string]manifestv1alpha1.GitHubRepositoryCustomPropertySpec, len(spec.CustomProperties))
	for _, property := range spec.CustomProperties {
		desiredByName[strings.ToLower(property.Name)] = property
	}

	var updates []ghapi.CustomPropertyValue
	for key, property := range desiredByName {
		if sameCustomPropertyValue(currentByName[key].Value, property.Value) {
			continue
		}

		updates = append(updates, ghapi.CustomPropertyValue{
			PropertyName: property.Name,
			Value:        property.Value,
		})
	}

	for key, property := range currentByName {
		if _, ok := desiredByName[key]; ok {
			continue
		}

		updates = append(updates, ghapi.CustomPropertyValue{
			PropertyName: property.Name,
			Value:        nil,
		})

		if property.Value == nil {
			updates = updates[:len(updates)-1]
		}
	}

	sort.Slice(updates, func(i, j int) bool {
		return strings.ToLower(updates[i].PropertyName) < strings.ToLower(updates[j].PropertyName)
	})

	if len(updates) == 0 {
		return false, nil
	}

	if err := client.UpdateCustomPropertyValues(ctx, spec.Owner, spec.Name, updates); err != nil {
		return false, err
	}

	return true, nil
}

func reconcileBranches(ctx context.Context, client *ghapi.Client, repository *ghapi.Repository, spec manifestv1alpha1.GitHubRepositorySpec) ([]string, bool, error) {
	if spec.Branches == nil {
		return nil, false, nil
	}

	protectedBranches, err := client.ListProtectedBranches(ctx, spec.Owner, spec.Name)
	if err != nil {
		return nil, false, err
	}

	protectedSet := make(map[string]struct{}, len(protectedBranches))
	for _, branch := range protectedBranches {
		protectedSet[branch.Name] = struct{}{}
	}

	desiredBranches := make(map[string]manifestv1alpha1.GitHubRepositoryBranchSpec, len(spec.Branches))
	for _, branch := range spec.Branches {
		desiredBranches[branch.Name] = branch
	}

	var messages []string
	changed := false

	for _, branch := range spec.Branches {
		if _, err := client.GetBranch(ctx, spec.Owner, spec.Name, branch.Name); err != nil {
			return messages, changed, fmt.Errorf("get branch %s: %w", branch.Name, err)
		}

		currentProtection, err := client.GetBranchProtection(ctx, spec.Owner, spec.Name, branch.Name)
		if err != nil && !ghapi.IsNotFound(err) {
			return messages, changed, fmt.Errorf("get protection for branch %s: %w", branch.Name, err)
		}
		if ghapi.IsNotFound(err) {
			currentProtection = nil
		}

		if branch.Protection == nil {
			if _, ok := protectedSet[branch.Name]; ok {
				if err := client.DeleteBranchProtection(ctx, spec.Owner, spec.Name, branch.Name); err != nil {
					return messages, changed, fmt.Errorf("delete protection for branch %s: %w", branch.Name, err)
				}
				delete(protectedSet, branch.Name)
				changed = true
				messages = append(messages, fmt.Sprintf("Cleared branch protection for %s/%s#%s", spec.Owner, spec.Name, branch.Name))
			}
			continue
		}

		currentState := currentBranchProtectionState(currentProtection)
		desiredState := overlayBranchProtectionState(currentState, branch.Protection)
		if desiredState.isZero() {
			if _, ok := protectedSet[branch.Name]; ok {
				if err := client.DeleteBranchProtection(ctx, spec.Owner, spec.Name, branch.Name); err != nil {
					return messages, changed, fmt.Errorf("delete protection for branch %s: %w", branch.Name, err)
				}
				delete(protectedSet, branch.Name)
				changed = true
				messages = append(messages, fmt.Sprintf("Cleared branch protection for %s/%s#%s", spec.Owner, spec.Name, branch.Name))
			}
			continue
		}

		if reflect.DeepEqual(currentState, desiredState) {
			delete(protectedSet, branch.Name)
			continue
		}

		if err := client.UpdateBranchProtection(ctx, spec.Owner, spec.Name, branch.Name, desiredState.toRequest()); err != nil {
			return messages, changed, fmt.Errorf("update protection for branch %s: %w", branch.Name, err)
		}

		delete(protectedSet, branch.Name)
		changed = true
		messages = append(messages, fmt.Sprintf("Updated branch protection for %s/%s#%s", spec.Owner, spec.Name, branch.Name))
	}

	var extraProtectedBranches []string
	for branchName := range protectedSet {
		if _, ok := desiredBranches[branchName]; ok {
			continue
		}
		extraProtectedBranches = append(extraProtectedBranches, branchName)
	}
	sort.Strings(extraProtectedBranches)

	for _, branchName := range extraProtectedBranches {
		if err := client.DeleteBranchProtection(ctx, spec.Owner, spec.Name, branchName); err != nil {
			return messages, changed, fmt.Errorf("delete protection for branch %s: %w", branchName, err)
		}

		changed = true
		messages = append(messages, fmt.Sprintf("Cleared branch protection for %s/%s#%s", spec.Owner, spec.Name, branchName))
	}

	return messages, changed, nil
}

func buildSecurityAndAnalysisUpdate(current *ghapi.SecurityAndAnalysis, desired *manifestv1alpha1.GitHubRepositorySecurityAndAnalysisSpec) *ghapi.SecurityAndAnalysis {
	if desired == nil {
		return nil
	}

	update := &ghapi.SecurityAndAnalysis{}
	changed := false

	compareAndSetSecuritySetting := func(currentSetting *ghapi.SecuritySetting, desiredSetting *manifestv1alpha1.GitHubRepositorySecuritySettingSpec, target **ghapi.SecuritySetting) {
		if desiredSetting == nil {
			return
		}
		if currentSecurityStatus(currentSetting) == desiredSetting.Status {
			return
		}

		*target = &ghapi.SecuritySetting{Status: desiredSetting.Status}
		changed = true
	}

	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting { return settings.AdvancedSecurity }), desired.AdvancedSecurity, &update.AdvancedSecurity)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting { return settings.CodeSecurity }), desired.CodeSecurity, &update.CodeSecurity)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting { return settings.SecretScanning }), desired.SecretScanning, &update.SecretScanning)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting {
		return settings.SecretScanningPushProtection
	}), desired.SecretScanningPushProtection, &update.SecretScanningPushProtection)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting {
		return settings.SecretScanningAIDetection
	}), desired.SecretScanningAIDetection, &update.SecretScanningAIDetection)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting {
		return settings.SecretScanningNonProviderPatterns
	}), desired.SecretScanningNonProviderPatterns, &update.SecretScanningNonProviderPatterns)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting {
		return settings.SecretScanningDelegatedAlertDismissal
	}), desired.SecretScanningDelegatedAlertDismissal, &update.SecretScanningDelegatedAlertDismissal)
	compareAndSetSecuritySetting(currentSecurity(current, func(settings *ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting {
		return settings.SecretScanningDelegatedBypass
	}), desired.SecretScanningDelegatedBypass, &update.SecretScanningDelegatedBypass)

	if !changed {
		return nil
	}

	return update
}

type pagesState struct {
	BuildType     *string
	CNAME         *string
	HTTPSEnforced *bool
	Source        *ghapi.PagesSource
}

func currentPagesState(current *ghapi.PagesSite) pagesState {
	if current == nil {
		return pagesState{}
	}

	state := pagesState{
		BuildType: cloneStringPointer(nonEmptyStringPointer(current.BuildType)),
		CNAME:     cloneStringPointer(current.CNAME),
		Source:    clonePagesSource(current.Source),
	}
	if current.HTTPSEnforced {
		state.HTTPSEnforced = boolPointer(true)
	}

	return state
}

func overlayPagesState(current pagesState, desired *manifestv1alpha1.GitHubRepositoryPagesSpec) pagesState {
	state := current
	if desired.BuildType != nil {
		state.BuildType = cloneStringPointer(desired.BuildType)
	}
	if desired.CNAME != nil {
		state.CNAME = cloneStringPointer(desired.CNAME)
	}
	if desired.HTTPSEnforced != nil {
		state.HTTPSEnforced = boolPointer(*desired.HTTPSEnforced)
	}
	if desired.Source != nil {
		state.Source = &ghapi.PagesSource{
			Branch: desired.Source.Branch,
			Path:   desired.Source.Path,
		}
	}

	return state
}

func (s pagesState) toRequest() ghapi.UpdatePagesRequest {
	request := ghapi.UpdatePagesRequest{
		BuildType:     cloneStringPointer(s.BuildType),
		HTTPSEnforced: cloneBoolPointer(s.HTTPSEnforced),
		Source:        clonePagesSource(s.Source),
	}

	if s.CNAME != nil {
		if *s.CNAME == "" {
			request.CNAME = nil
		} else {
			request.CNAME = *s.CNAME
		}
	}

	return request
}

type branchProtectionState struct {
	RequiredStatusChecks        *requiredStatusChecksState
	EnforceAdmins               bool
	PullRequestReviews          *pullRequestReviewsState
	Restrictions                actorAllowanceState
	BypassPullRequestAllowances actorAllowanceState
	RequiredLinearHistory       bool
}

func (s branchProtectionState) isZero() bool {
	return s.RequiredStatusChecks == nil &&
		!s.EnforceAdmins &&
		s.PullRequestReviews == nil &&
		s.Restrictions.isEmpty() &&
		s.BypassPullRequestAllowances.isEmpty() &&
		!s.RequiredLinearHistory
}

func (s branchProtectionState) toRequest() map[string]any {
	request := map[string]any{
		"required_status_checks":         nil,
		"enforce_admins":                 nil,
		"required_pull_request_reviews":  nil,
		"restrictions":                   nil,
		"required_linear_history":        s.RequiredLinearHistory,
		"bypass_pull_request_allowances": s.BypassPullRequestAllowances.toRequest(),
	}

	if s.RequiredStatusChecks != nil {
		request["required_status_checks"] = s.RequiredStatusChecks.toRequest()
	}
	if s.EnforceAdmins {
		request["enforce_admins"] = true
	}
	if s.PullRequestReviews != nil {
		request["required_pull_request_reviews"] = s.PullRequestReviews.toRequest()
	}
	if !s.Restrictions.isEmpty() {
		request["restrictions"] = s.Restrictions.toRequest()
	}

	return request
}

type requiredStatusChecksState struct {
	Strict bool
	Checks []requiredStatusCheckState
}

func (s *requiredStatusChecksState) toRequest() map[string]any {
	checks := make([]map[string]any, 0, len(s.Checks))
	for _, check := range s.Checks {
		entry := map[string]any{"context": check.Context}
		if check.AppID != nil {
			entry["app_id"] = *check.AppID
		}
		checks = append(checks, entry)
	}

	return map[string]any{
		"strict": s.Strict,
		"checks": checks,
	}
}

type requiredStatusCheckState struct {
	Context string
	AppID   *int64
}

type pullRequestReviewsState struct {
	DismissalRestrictions        actorAllowanceState
	DismissStaleReviews          bool
	RequireCodeOwnerReviews      bool
	RequiredApprovingReviewCount int
	RequireLastPushApproval      bool
}

func (s *pullRequestReviewsState) toRequest() map[string]any {
	return map[string]any{
		"dismissal_restrictions":          s.DismissalRestrictions.toRequest(),
		"dismiss_stale_reviews":           s.DismissStaleReviews,
		"require_code_owner_reviews":      s.RequireCodeOwnerReviews,
		"required_approving_review_count": s.RequiredApprovingReviewCount,
		"require_last_push_approval":      s.RequireLastPushApproval,
	}
}

type actorAllowanceState struct {
	Users []string
	Teams []string
	Apps  []string
}

func (s actorAllowanceState) isEmpty() bool {
	return len(s.Users) == 0 && len(s.Teams) == 0 && len(s.Apps) == 0
}

func (s actorAllowanceState) toRequest() map[string]any {
	return map[string]any{
		"users": append([]string(nil), s.Users...),
		"teams": append([]string(nil), s.Teams...),
		"apps":  append([]string(nil), s.Apps...),
	}
}

func currentBranchProtectionState(current *ghapi.BranchProtection) branchProtectionState {
	if current == nil {
		return branchProtectionState{}
	}

	state := branchProtectionState{
		Restrictions:          currentActorAllowance(current.Restrictions),
		RequiredLinearHistory: enabledSetting(current.RequiredLinearHistory),
	}

	if enabledSetting(current.EnforceAdmins) {
		state.EnforceAdmins = true
	}

	if current.RequiredStatusChecks != nil {
		state.RequiredStatusChecks = &requiredStatusChecksState{
			Strict: current.RequiredStatusChecks.Strict,
			Checks: currentRequiredStatusChecks(current.RequiredStatusChecks.Checks),
		}
	}

	if current.RequiredPullRequestReviews != nil {
		state.PullRequestReviews = &pullRequestReviewsState{
			DismissalRestrictions:        currentActorAllowance(current.RequiredPullRequestReviews.DismissalRestrictions),
			DismissStaleReviews:          current.RequiredPullRequestReviews.DismissStaleReviews,
			RequireCodeOwnerReviews:      current.RequiredPullRequestReviews.RequireCodeOwnerReviews,
			RequiredApprovingReviewCount: current.RequiredPullRequestReviews.RequiredApprovingReviewCount,
			RequireLastPushApproval:      current.RequiredPullRequestReviews.RequireLastPushApproval,
		}
		if current.BypassPullRequestAllowances == nil {
			state.BypassPullRequestAllowances = currentActorAllowance(current.RequiredPullRequestReviews.BypassPullRequestAllowances)
		}
	}

	if current.BypassPullRequestAllowances != nil {
		state.BypassPullRequestAllowances = currentActorAllowance(current.BypassPullRequestAllowances)
	}

	return state
}

func overlayBranchProtectionState(current branchProtectionState, desired *manifestv1alpha1.GitHubRepositoryBranchProtectionSpec) branchProtectionState {
	state := current

	if desired.RequiredStatusChecks != nil {
		state.RequiredStatusChecks = &requiredStatusChecksState{
			Strict: desired.RequiredStatusChecks.Strict,
			Checks: desiredRequiredStatusChecks(desired.RequiredStatusChecks.Checks),
		}
	}
	if desired.EnforceAdmins != nil {
		state.EnforceAdmins = *desired.EnforceAdmins
	}
	if desired.PullRequestReviews != nil {
		reviews := state.PullRequestReviews
		if reviews == nil {
			reviews = &pullRequestReviewsState{}
		}
		if desired.PullRequestReviews.DismissalRestrictions != nil {
			reviews.DismissalRestrictions = desiredActorAllowance(desired.PullRequestReviews.DismissalRestrictions)
		}
		if desired.PullRequestReviews.DismissStaleReviews != nil {
			reviews.DismissStaleReviews = *desired.PullRequestReviews.DismissStaleReviews
		}
		if desired.PullRequestReviews.RequireCodeOwnerReviews != nil {
			reviews.RequireCodeOwnerReviews = *desired.PullRequestReviews.RequireCodeOwnerReviews
		}
		if desired.PullRequestReviews.RequiredApprovingReviewCount != nil {
			reviews.RequiredApprovingReviewCount = *desired.PullRequestReviews.RequiredApprovingReviewCount
		}
		if desired.PullRequestReviews.RequireLastPushApproval != nil {
			reviews.RequireLastPushApproval = *desired.PullRequestReviews.RequireLastPushApproval
		}
		state.PullRequestReviews = reviews
	}
	if desired.Restrictions != nil {
		state.Restrictions = desiredActorAllowance(desired.Restrictions)
	}
	if desired.BypassPullRequestAllowances != nil {
		state.BypassPullRequestAllowances = desiredActorAllowance(desired.BypassPullRequestAllowances)
	}
	if desired.RequiredLinearHistory != nil {
		state.RequiredLinearHistory = *desired.RequiredLinearHistory
	}

	return state
}

func desiredRequiredStatusChecks(checks []manifestv1alpha1.GitHubRequiredStatusCheckSpec) []requiredStatusCheckState {
	normalized := make([]requiredStatusCheckState, 0, len(checks))
	for _, check := range checks {
		normalized = append(normalized, requiredStatusCheckState{
			Context: check.Context,
			AppID:   cloneInt64Pointer(check.AppID),
		})
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Context == normalized[j].Context {
			return dereferenceInt64(normalized[i].AppID) < dereferenceInt64(normalized[j].AppID)
		}
		return normalized[i].Context < normalized[j].Context
	})

	return normalized
}

func currentRequiredStatusChecks(checks []ghapi.RequiredStatusCheck) []requiredStatusCheckState {
	normalized := make([]requiredStatusCheckState, 0, len(checks))
	for _, check := range checks {
		normalized = append(normalized, requiredStatusCheckState{
			Context: check.CheckContext(),
			AppID:   cloneInt64Pointer(check.AppID),
		})
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Context == normalized[j].Context {
			return dereferenceInt64(normalized[i].AppID) < dereferenceInt64(normalized[j].AppID)
		}
		return normalized[i].Context < normalized[j].Context
	})

	return normalized
}

func desiredActorAllowance(allowance *manifestv1alpha1.GitHubActorAllowanceSpec) actorAllowanceState {
	if allowance == nil {
		return actorAllowanceState{}
	}

	return actorAllowanceState{
		Users: normalizeStrings(allowance.Users),
		Teams: normalizeStrings(allowance.Teams),
		Apps:  normalizeStrings(allowance.Apps),
	}
}

func currentActorAllowance(allowance *ghapi.ActorAllowance) actorAllowanceState {
	if allowance == nil {
		return actorAllowanceState{}
	}

	state := actorAllowanceState{
		Users: make([]string, 0, len(allowance.Users)),
		Teams: make([]string, 0, len(allowance.Teams)),
		Apps:  make([]string, 0, len(allowance.Apps)),
	}

	for _, user := range allowance.Users {
		state.Users = append(state.Users, user.Login)
	}
	for _, team := range allowance.Teams {
		state.Teams = append(state.Teams, team.Slug)
	}
	for _, app := range allowance.Apps {
		state.Apps = append(state.Apps, app.Slug)
	}

	state.Users = normalizeStrings(state.Users)
	state.Teams = normalizeStrings(state.Teams)
	state.Apps = normalizeStrings(state.Apps)

	return state
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

func sameCustomPropertyValue(current any, desired any) bool {
	return canonicalValue(current) == canonicalValue(desired)
}

func canonicalValue(value any) string {
	normalized := normalizeValue(value)
	payload, _ := json.Marshal(normalized)
	return string(payload)
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case []any:
		normalized := make([]any, 0, len(typed))
		for _, item := range typed {
			normalized = append(normalized, normalizeValue(item))
		}
		return normalized
	case []string:
		normalized := make([]any, 0, len(typed))
		for _, item := range typed {
			normalized = append(normalized, item)
		}
		return normalized
	case map[string]any:
		normalized := make(map[string]any, len(typed))
		for key, item := range typed {
			normalized[key] = normalizeValue(item)
		}
		return normalized
	default:
		return typed
	}
}

func currentSecurity(current *ghapi.SecurityAndAnalysis, selector func(*ghapi.SecurityAndAnalysis) *ghapi.SecuritySetting) *ghapi.SecuritySetting {
	if current == nil {
		return nil
	}
	return selector(current)
}

func currentSecurityStatus(setting *ghapi.SecuritySetting) string {
	if setting == nil {
		return ""
	}
	return setting.Status
}

func dereferenceString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func dereferenceInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBoolPointer(value *bool) *bool {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func clonePagesSource(value *ghapi.PagesSource) *ghapi.PagesSource {
	if value == nil {
		return nil
	}

	return &ghapi.PagesSource{
		Branch: value.Branch,
		Path:   value.Path,
	}
}

func boolPointer(value bool) *bool {
	return &value
}

func nonEmptyStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func featuresValue(features *manifestv1alpha1.GitHubRepositoryFeaturesSpec, accessor func(*manifestv1alpha1.GitHubRepositoryFeaturesSpec) *bool) *bool {
	if features == nil {
		return nil
	}

	return accessor(features)
}

func mergePolicyValue(policy *manifestv1alpha1.GitHubRepositoryMergePolicySpec, accessor func(*manifestv1alpha1.GitHubRepositoryMergePolicySpec) *bool) *bool {
	if policy == nil {
		return nil
	}

	return accessor(policy)
}

func enabledSetting(setting *ghapi.EnabledSetting) bool {
	return setting != nil && setting.Enabled
}
