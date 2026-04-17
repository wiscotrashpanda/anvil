package reconcile

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	manifestv1alpha1 "github.com/emkaytec/alloy/manifest/v1alpha1"
	hcpapi "github.com/emkaytec/anvil/internal/hcpterraform"
	"github.com/emkaytec/anvil/internal/manifest"
)

func reconcileHCPTerraformWorkspace(ctx context.Context, client *hcpapi.Client, loaded manifest.LoadedHCPTerraformWorkspaceManifest) ([]string, error) {
	spec := loaded.Manifest.Spec
	workspaceID := fmt.Sprintf("%s/%s", spec.Organization, spec.Name)

	messages := []string{fmt.Sprintf("Reconciling HCPTerraformWorkspace %s", workspaceID)}

	if err := validateSupportedHCPTerraformWorkspaceSpec(spec); err != nil {
		return messages, fmt.Errorf("reconcile %s: %w", workspaceID, err)
	}

	workspace, created, err := ensureHCPTerraformWorkspace(ctx, client, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile %s: %w", workspaceID, err)
	}

	changed := created
	if created {
		messages = append(messages, fmt.Sprintf("Created HCP Terraform workspace %s", workspaceID))
	}

	workspace, updatedSettings, err := reconcileHCPTerraformWorkspaceSettings(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile workspace settings for %s: %w", workspaceID, err)
	}
	if updatedSettings {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated workspace settings for %s", workspaceID))
	}

	updatedTags, err := reconcileHCPTerraformWorkspaceTags(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile tags for %s: %w", workspaceID, err)
	}
	if updatedTags {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated workspace tags for %s", workspaceID))
	}

	updatedTagBindings, err := reconcileHCPTerraformWorkspaceTagBindings(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile tag bindings for %s: %w", workspaceID, err)
	}
	if updatedTagBindings {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated workspace tag bindings for %s", workspaceID))
	}

	updatedConsumers, err := reconcileHCPTerraformWorkspaceRemoteStateConsumers(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile remote state consumers for %s: %w", workspaceID, err)
	}
	if updatedConsumers {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated remote state consumers for %s", workspaceID))
	}

	updatedVariables, err := reconcileHCPTerraformWorkspaceVariables(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile variables for %s: %w", workspaceID, err)
	}
	if updatedVariables {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated workspace variables for %s", workspaceID))
	}

	updatedVariableSets, err := reconcileHCPTerraformWorkspaceVariableSets(ctx, client, workspace, spec)
	if err != nil {
		return messages, fmt.Errorf("reconcile variable sets for %s: %w", workspaceID, err)
	}
	if updatedVariableSets {
		changed = true
		messages = append(messages, fmt.Sprintf("Updated workspace variable set assignments for %s", workspaceID))
	}

	if !changed {
		messages = append(messages, fmt.Sprintf("HCPTerraformWorkspace %s is up to date", workspaceID))
	}

	return messages, nil
}

func validateSupportedHCPTerraformWorkspaceSpec(spec manifestv1alpha1.HCPTerraformWorkspaceSpec) error {
	if spec.SSHKeyID != nil {
		return fmt.Errorf("spec.sshKeyID is not supported by the current anvil HCP Terraform reconciler yet")
	}
	if len(spec.RunTriggers) > 0 {
		return fmt.Errorf("spec.runTriggers is not supported by the current anvil HCP Terraform reconciler yet")
	}
	if len(spec.TeamAccess) > 0 {
		return fmt.Errorf("spec.teamAccess is not supported by the current anvil HCP Terraform reconciler yet")
	}
	if len(spec.Notifications) > 0 {
		return fmt.Errorf("spec.notifications is not supported by the current anvil HCP Terraform reconciler yet")
	}

	return nil
}

func ensureHCPTerraformWorkspace(ctx context.Context, client *hcpapi.Client, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (*hcpapi.Workspace, bool, error) {
	workspace, err := client.GetWorkspace(ctx, spec.Organization, spec.Name)
	if err == nil {
		return workspace, false, nil
	}
	if !hcpapi.IsNotFound(err) {
		return nil, false, err
	}

	workspace, err = client.CreateWorkspace(ctx, spec.Organization, spec.Name, hcpWorkspaceRequestFromSpec(spec))
	if err != nil {
		return nil, false, err
	}

	return workspace, true, nil
}

func reconcileHCPTerraformWorkspaceSettings(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (*hcpapi.Workspace, bool, error) {
	request := hcpWorkspaceUpdateRequest(workspace, spec)
	if request == nil {
		return workspace, false, nil
	}

	updatedWorkspace, err := client.UpdateWorkspace(ctx, workspace.ID, *request)
	if err != nil {
		return nil, false, err
	}

	return updatedWorkspace, true, nil
}

func reconcileHCPTerraformWorkspaceTags(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (bool, error) {
	if spec.Tags == nil {
		return false, nil
	}

	currentTags, err := client.ListWorkspaceTags(ctx, workspace.ID)
	if err != nil {
		return false, err
	}

	currentByName := make(map[string]hcpapi.WorkspaceTag, len(currentTags))
	for _, tag := range currentTags {
		currentByName[strings.ToLower(strings.TrimSpace(tag.Name))] = tag
	}

	desiredNames := normalizeTrimmedStrings(spec.Tags)
	var toAdd []string
	for _, name := range desiredNames {
		if _, ok := currentByName[strings.ToLower(name)]; ok {
			continue
		}
		toAdd = append(toAdd, name)
	}

	desiredSet := make(map[string]struct{}, len(desiredNames))
	for _, name := range desiredNames {
		desiredSet[strings.ToLower(name)] = struct{}{}
	}

	var toRemove []string
	for key, tag := range currentByName {
		if _, ok := desiredSet[key]; ok {
			continue
		}
		toRemove = append(toRemove, tag.ID)
	}

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return false, nil
	}

	if err := client.AddWorkspaceTags(ctx, workspace.ID, toAdd); err != nil {
		return false, err
	}
	if err := client.RemoveWorkspaceTags(ctx, workspace.ID, toRemove); err != nil {
		return false, err
	}

	return true, nil
}

func reconcileHCPTerraformWorkspaceTagBindings(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (bool, error) {
	if spec.TagBindings == nil {
		return false, nil
	}

	currentBindings, err := client.ListWorkspaceTagBindings(ctx, workspace.ID)
	if err != nil {
		return false, err
	}

	desiredBindings := workspaceTagBindingsFromSpec(spec.TagBindings)
	if reflect.DeepEqual(normalizeWorkspaceTagBindings(currentBindings), normalizeWorkspaceTagBindings(desiredBindings)) {
		return false, nil
	}

	if err := client.ReplaceWorkspaceTagBindings(ctx, workspace.ID, desiredBindings); err != nil {
		return false, err
	}

	return true, nil
}

func reconcileHCPTerraformWorkspaceRemoteStateConsumers(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (bool, error) {
	if spec.RemoteStateConsumerIDs == nil {
		return false, nil
	}

	currentIDs, err := client.ListRemoteStateConsumers(ctx, workspace.ID)
	if err != nil {
		return false, err
	}

	desiredIDs := normalizeTrimmedStrings(spec.RemoteStateConsumerIDs)
	currentIDs = normalizeTrimmedStrings(currentIDs)

	toAdd := differencePreservingOriginal(desiredIDs, currentIDs)
	toRemove := differencePreservingOriginal(currentIDs, desiredIDs)

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return false, nil
	}

	if err := client.AddRemoteStateConsumers(ctx, workspace.ID, toAdd); err != nil {
		return false, err
	}
	if err := client.RemoveRemoteStateConsumers(ctx, workspace.ID, toRemove); err != nil {
		return false, err
	}

	return true, nil
}

func reconcileHCPTerraformWorkspaceVariables(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (bool, error) {
	if spec.Variables == nil {
		return false, nil
	}

	currentVariables, err := client.ListVariables(ctx, workspace.ID)
	if err != nil {
		return false, err
	}

	currentByKey := make(map[string]hcpapi.WorkspaceVariable, len(currentVariables))
	for _, variable := range currentVariables {
		currentByKey[workspaceVariableMapKey(variable.Category, variable.Key)] = variable
	}

	desiredByKey := make(map[string]manifestv1alpha1.HCPTerraformWorkspaceVariableSpec, len(spec.Variables))
	for _, variable := range spec.Variables {
		desiredByKey[workspaceVariableMapKey(variable.Category, variable.Key)] = variable
	}

	changed := false

	for key, desired := range desiredByKey {
		current, ok := currentByKey[key]
		target := hcpapi.WorkspaceVariable{
			Key:         desired.Key,
			Value:       desired.Value,
			Description: dereferenceString(desired.Description),
			Category:    desired.Category,
			HCL:         desired.HCL != nil && *desired.HCL,
			Sensitive:   desired.Sensitive != nil && *desired.Sensitive,
		}

		if !ok {
			if err := client.CreateVariable(ctx, workspace.ID, target); err != nil {
				return false, err
			}
			changed = true
			continue
		}

		if hcpWorkspaceVariableMatches(current, target) {
			continue
		}

		if err := client.UpdateVariable(ctx, workspace.ID, current.ID, target); err != nil {
			return false, err
		}
		changed = true
	}

	for key, current := range currentByKey {
		if _, ok := desiredByKey[key]; ok {
			continue
		}

		if err := client.DeleteVariable(ctx, workspace.ID, current.ID); err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}

func reconcileHCPTerraformWorkspaceVariableSets(ctx context.Context, client *hcpapi.Client, workspace *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) (bool, error) {
	if spec.VariableSetIDs == nil {
		return false, nil
	}

	currentVariableSets, err := client.ListWorkspaceVariableSets(ctx, workspace.ID)
	if err != nil {
		return false, err
	}

	currentIDs := make([]string, 0, len(currentVariableSets))
	for _, variableSet := range currentVariableSets {
		currentIDs = append(currentIDs, variableSet.ID)
	}

	desiredIDs := normalizeTrimmedStrings(spec.VariableSetIDs)
	currentIDs = normalizeTrimmedStrings(currentIDs)

	toAdd := differencePreservingOriginal(desiredIDs, currentIDs)
	toRemove := differencePreservingOriginal(currentIDs, desiredIDs)

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return false, nil
	}

	for _, variableSetID := range toAdd {
		if err := client.AddVariableSetToWorkspace(ctx, variableSetID, workspace.ID); err != nil {
			return false, err
		}
	}

	for _, variableSetID := range toRemove {
		if err := client.RemoveVariableSetFromWorkspace(ctx, variableSetID, workspace.ID); err != nil {
			return false, err
		}
	}

	return true, nil
}

func hcpWorkspaceRequestFromSpec(spec manifestv1alpha1.HCPTerraformWorkspaceSpec) hcpapi.WorkspaceRequest {
	return hcpapi.WorkspaceRequest{
		ProjectID:                   spec.ProjectID,
		Description:                 spec.Description,
		TerraformVersion:            spec.TerraformVersion,
		WorkingDirectory:            spec.WorkingDirectory,
		ExecutionMode:               spec.ExecutionMode,
		AgentPoolID:                 spec.AgentPoolID,
		AllowDestroyPlan:            spec.AllowDestroyPlan,
		AssessmentsEnabled:          spec.AssessmentsEnabled,
		AutoApply:                   spec.AutoApply,
		AutoApplyRunTrigger:         spec.AutoApplyRunTrigger,
		AutoDestroyAt:               spec.AutoDestroyAt,
		AutoDestroyActivityDuration: spec.AutoDestroyActivityDuration,
		FileTriggersEnabled:         spec.FileTriggersEnabled,
		GlobalRemoteState:           spec.GlobalRemoteState,
		ProjectRemoteState:          spec.ProjectRemoteState,
		QueueAllRuns:                spec.QueueAllRuns,
		SourceName:                  spec.SourceName,
		SourceURL:                   spec.SourceURL,
		SpeculativeEnabled:          spec.SpeculativeEnabled,
		TriggerPatterns:             spec.TriggerPatterns,
		TriggerPrefixes:             spec.TriggerPrefixes,
		SettingOverwrites:           hcpSettingOverwritesFromSpec(spec.SettingOverwrites),
		VCSRepo:                     hcpVCSRepoFromSpec(spec.VCSRepo),
		TagBindings:                 workspaceTagBindingsFromSpec(spec.TagBindings),
	}
}

func hcpWorkspaceUpdateRequest(current *hcpapi.Workspace, spec manifestv1alpha1.HCPTerraformWorkspaceSpec) *hcpapi.WorkspaceRequest {
	request := hcpapi.WorkspaceRequest{}
	changed := false

	if spec.ProjectID != nil && current.ProjectID != *spec.ProjectID {
		request.ProjectID = spec.ProjectID
		changed = true
	}
	if spec.Description != nil && current.Description != *spec.Description {
		request.Description = spec.Description
		changed = true
	}
	if spec.TerraformVersion != nil && current.TerraformVersion != *spec.TerraformVersion {
		request.TerraformVersion = spec.TerraformVersion
		changed = true
	}
	if spec.WorkingDirectory != nil && current.WorkingDirectory != *spec.WorkingDirectory {
		request.WorkingDirectory = spec.WorkingDirectory
		changed = true
	}
	if spec.ExecutionMode != nil && current.ExecutionMode != *spec.ExecutionMode {
		request.ExecutionMode = spec.ExecutionMode
		changed = true
	}
	if spec.AgentPoolID != nil && current.AgentPoolID != *spec.AgentPoolID {
		request.AgentPoolID = spec.AgentPoolID
		changed = true
	}
	if spec.AllowDestroyPlan != nil && current.AllowDestroyPlan != *spec.AllowDestroyPlan {
		request.AllowDestroyPlan = spec.AllowDestroyPlan
		changed = true
	}
	if spec.AssessmentsEnabled != nil && current.AssessmentsEnabled != *spec.AssessmentsEnabled {
		request.AssessmentsEnabled = spec.AssessmentsEnabled
		changed = true
	}
	if spec.AutoApply != nil && current.AutoApply != *spec.AutoApply {
		request.AutoApply = spec.AutoApply
		changed = true
	}
	if spec.AutoApplyRunTrigger != nil && current.AutoApplyRunTrigger != *spec.AutoApplyRunTrigger {
		request.AutoApplyRunTrigger = spec.AutoApplyRunTrigger
		changed = true
	}
	if spec.AutoDestroyAt != nil && current.AutoDestroyAt != *spec.AutoDestroyAt {
		request.AutoDestroyAt = spec.AutoDestroyAt
		changed = true
	}
	if spec.AutoDestroyActivityDuration != nil && current.AutoDestroyActivityDuration != *spec.AutoDestroyActivityDuration {
		request.AutoDestroyActivityDuration = spec.AutoDestroyActivityDuration
		changed = true
	}
	if spec.FileTriggersEnabled != nil && current.FileTriggersEnabled != *spec.FileTriggersEnabled {
		request.FileTriggersEnabled = spec.FileTriggersEnabled
		changed = true
	}
	if spec.GlobalRemoteState != nil && current.GlobalRemoteState != *spec.GlobalRemoteState {
		request.GlobalRemoteState = spec.GlobalRemoteState
		changed = true
	}
	if spec.ProjectRemoteState != nil && current.ProjectRemoteState != *spec.ProjectRemoteState {
		request.ProjectRemoteState = spec.ProjectRemoteState
		changed = true
	}
	if spec.QueueAllRuns != nil && current.QueueAllRuns != *spec.QueueAllRuns {
		request.QueueAllRuns = spec.QueueAllRuns
		changed = true
	}
	if spec.SourceName != nil && current.SourceName != *spec.SourceName {
		request.SourceName = spec.SourceName
		changed = true
	}
	if spec.SourceURL != nil && current.SourceURL != *spec.SourceURL {
		request.SourceURL = spec.SourceURL
		changed = true
	}
	if spec.SpeculativeEnabled != nil && current.SpeculativeEnabled != *spec.SpeculativeEnabled {
		request.SpeculativeEnabled = spec.SpeculativeEnabled
		changed = true
	}
	if spec.TriggerPatterns != nil && !reflect.DeepEqual(normalizeTrimmedStringsPreserveCase(current.TriggerPatterns), normalizeTrimmedStringsPreserveCase(spec.TriggerPatterns)) {
		request.TriggerPatterns = spec.TriggerPatterns
		changed = true
	}
	if spec.TriggerPrefixes != nil && !reflect.DeepEqual(normalizeTrimmedStringsPreserveCase(current.TriggerPrefixes), normalizeTrimmedStringsPreserveCase(spec.TriggerPrefixes)) {
		request.TriggerPrefixes = spec.TriggerPrefixes
		changed = true
	}

	desiredOverwrites := hcpSettingOverwritesFromSpec(spec.SettingOverwrites)
	if spec.SettingOverwrites != nil && !reflect.DeepEqual(current.SettingOverwrites, desiredOverwrites) {
		request.SettingOverwrites = desiredOverwrites
		changed = true
	}

	desiredVCSRepo := hcpVCSRepoFromSpec(spec.VCSRepo)
	if spec.VCSRepo != nil && !reflect.DeepEqual(current.VCSRepo, desiredVCSRepo) {
		request.VCSRepo = desiredVCSRepo
		changed = true
	}

	if !changed {
		return nil
	}

	return &request
}

func hcpSettingOverwritesFromSpec(spec *manifestv1alpha1.HCPTerraformWorkspaceSettingOverwrites) *hcpapi.WorkspaceSettingOverwrites {
	if spec == nil {
		return nil
	}

	return &hcpapi.WorkspaceSettingOverwrites{
		ExecutionMode:       spec.ExecutionMode != nil && *spec.ExecutionMode,
		AgentPoolID:         spec.AgentPoolID != nil && *spec.AgentPoolID,
		AutoApply:           spec.AutoApply != nil && *spec.AutoApply,
		FileTriggersEnabled: spec.FileTriggersEnabled != nil && *spec.FileTriggersEnabled,
		GlobalRemoteState:   spec.GlobalRemoteState != nil && *spec.GlobalRemoteState,
		QueueAllRuns:        spec.QueueAllRuns != nil && *spec.QueueAllRuns,
		SpeculativeEnabled:  spec.SpeculativeEnabled != nil && *spec.SpeculativeEnabled,
		TerraformVersion:    spec.TerraformVersion != nil && *spec.TerraformVersion,
		TriggerPatterns:     spec.TriggerPatterns != nil && *spec.TriggerPatterns,
		TriggerPrefixes:     spec.TriggerPrefixes != nil && *spec.TriggerPrefixes,
		WorkingDirectory:    spec.WorkingDirectory != nil && *spec.WorkingDirectory,
	}
}

func hcpVCSRepoFromSpec(spec *manifestv1alpha1.HCPTerraformWorkspaceVCSRepoSpec) *hcpapi.WorkspaceVCSRepo {
	if spec == nil {
		return nil
	}

	return &hcpapi.WorkspaceVCSRepo{
		Identifier:        dereferenceString(spec.Identifier),
		Branch:            dereferenceString(spec.Branch),
		IngressSubmodules: spec.IngressSubmodules != nil && *spec.IngressSubmodules,
		OAuthTokenID:      dereferenceString(spec.OAuthTokenID),
		TagsRegex:         dereferenceString(spec.TagsRegex),
	}
}

func workspaceTagBindingsFromSpec(specs []manifestv1alpha1.HCPTerraformWorkspaceTagBindingSpec) []hcpapi.WorkspaceTagBinding {
	if specs == nil {
		return nil
	}

	bindings := make([]hcpapi.WorkspaceTagBinding, 0, len(specs))
	for _, spec := range specs {
		bindings = append(bindings, hcpapi.WorkspaceTagBinding{
			Key:   spec.Key,
			Value: spec.Value,
		})
	}

	return bindings
}

func normalizeWorkspaceTagBindings(bindings []hcpapi.WorkspaceTagBinding) []hcpapi.WorkspaceTagBinding {
	if bindings == nil {
		return nil
	}

	normalized := make([]hcpapi.WorkspaceTagBinding, 0, len(bindings))
	for _, binding := range bindings {
		normalized = append(normalized, hcpapi.WorkspaceTagBinding{
			Key:   strings.ToLower(strings.TrimSpace(binding.Key)),
			Value: strings.TrimSpace(binding.Value),
		})
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Key == normalized[j].Key {
			return normalized[i].Value < normalized[j].Value
		}
		return normalized[i].Key < normalized[j].Key
	})

	return normalized
}

func workspaceVariableMapKey(category string, key string) string {
	return strings.ToLower(strings.TrimSpace(category)) + ":" + strings.TrimSpace(key)
}

func hcpWorkspaceVariableMatches(current hcpapi.WorkspaceVariable, desired hcpapi.WorkspaceVariable) bool {
	if current.Key != desired.Key ||
		current.Description != desired.Description ||
		current.Category != desired.Category ||
		current.HCL != desired.HCL ||
		current.Sensitive != desired.Sensitive {
		return false
	}

	if current.Sensitive || desired.Sensitive {
		return false
	}

	return current.Value == desired.Value
}

func normalizeTrimmedStrings(values []string) []string {
	if values == nil {
		return nil
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	sort.Strings(normalized)
	return normalized
}

func normalizeTrimmedStringsPreserveCase(values []string) []string {
	return normalizeTrimmedStrings(values)
}

func differencePreservingOriginal(left []string, right []string) []string {
	rightSet := make(map[string]struct{}, len(right))
	for _, value := range right {
		rightSet[value] = struct{}{}
	}

	var diff []string
	for _, value := range left {
		if _, ok := rightSet[value]; ok {
			continue
		}
		diff = append(diff, value)
	}

	return diff
}
