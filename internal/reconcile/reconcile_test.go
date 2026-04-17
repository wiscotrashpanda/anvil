package reconcile

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	manifestv1alpha1 "github.com/emkaytec/alloy/manifest/v1alpha1"
	ghapi "github.com/emkaytec/anvil/internal/github"
	hcpapi "github.com/emkaytec/anvil/internal/hcpterraform"
	"github.com/emkaytec/anvil/internal/manifest"
)

func TestPlanRunCreatesRepository(t *testing.T) {
	t.Parallel()

	var createdRequest map[string]any

	client := newGitHubTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusNotFound, map[string]any{"message": "Not Found"})
		case r.Method == http.MethodGet && r.URL.Path == "/users/example-org":
			return jsonResponse(http.StatusOK, map[string]any{
				"login": "example-org",
				"type":  "Organization",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/orgs/example-org/repos":
			decodeJSON(t, r, &createdRequest)
			return jsonResponse(http.StatusCreated, map[string]any{
				"name":                   "example-repo",
				"full_name":              "example-org/example-repo",
				"visibility":             "private",
				"description":            "Created by Anvil",
				"homepage":               "",
				"default_branch":         "",
				"topics":                 []string{},
				"owner":                  map[string]any{"login": "example-org", "type": "Organization"},
				"has_issues":             true,
				"has_projects":           false,
				"has_wiki":               true,
				"allow_squash_merge":     true,
				"allow_merge_commit":     false,
				"allow_rebase_merge":     true,
				"allow_auto_merge":       true,
				"allow_update_branch":    false,
				"delete_branch_on_merge": true,
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			return nil
		}
	})

	visibility := "private"
	description := "Created by Anvil"
	spec := manifestv1alpha1.GitHubRepositorySpec{
		Owner:       "example-org",
		Name:        "example-repo",
		Visibility:  &visibility,
		Description: &description,
		AutoInit:    true,
		Features: &manifestv1alpha1.GitHubRepositoryFeaturesSpec{
			HasIssues:   boolPtr(true),
			HasProjects: boolPtr(false),
			HasWiki:     boolPtr(true),
		},
		MergePolicy: &manifestv1alpha1.GitHubRepositoryMergePolicySpec{
			AllowSquashMerge:    boolPtr(true),
			AllowMergeCommit:    boolPtr(false),
			AllowRebaseMerge:    boolPtr(true),
			AllowAutoMerge:      boolPtr(true),
			DeleteBranchOnMerge: boolPtr(true),
		},
	}

	messages, err := newGitHubRepositoryPlan(spec).Run(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if createdRequest["name"] != "example-repo" {
		t.Fatalf("expected create request to include repo name, got %#v", createdRequest)
	}
	if createdRequest["visibility"] != "private" {
		t.Fatalf("expected create request to include visibility, got %#v", createdRequest)
	}
	if createdRequest["auto_init"] != true {
		t.Fatalf("expected create request to include auto_init, got %#v", createdRequest)
	}

	expectedMessages := []string{
		"Reconciling GitHubRepository example-org/example-repo",
		"Created GitHub repository example-org/example-repo",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func TestPlanRunReconcilesRepositoryDrift(t *testing.T) {
	t.Parallel()

	var patchRepositoryRequest map[string]any
	var replaceTopicsRequest map[string]any

	client := newGitHubTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusOK, map[string]any{
				"name":                   "example-repo",
				"full_name":              "example-org/example-repo",
				"visibility":             "public",
				"description":            "Old description",
				"homepage":               "",
				"default_branch":         "master",
				"topics":                 []string{"legacy"},
				"owner":                  map[string]any{"login": "example-org", "type": "Organization"},
				"has_issues":             false,
				"has_projects":           true,
				"has_wiki":               false,
				"allow_squash_merge":     true,
				"allow_merge_commit":     true,
				"allow_rebase_merge":     false,
				"allow_auto_merge":       false,
				"allow_update_branch":    false,
				"delete_branch_on_merge": false,
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/repos/example-org/example-repo":
			decodeJSON(t, r, &patchRepositoryRequest)
			return jsonResponse(http.StatusOK, map[string]any{
				"name":                   "example-repo",
				"full_name":              "example-org/example-repo",
				"visibility":             "private",
				"description":            "Managed description",
				"homepage":               "https://example.com",
				"default_branch":         "main",
				"topics":                 []string{"legacy"},
				"owner":                  map[string]any{"login": "example-org", "type": "Organization"},
				"has_issues":             true,
				"has_projects":           false,
				"has_wiki":               true,
				"allow_squash_merge":     false,
				"allow_merge_commit":     false,
				"allow_rebase_merge":     true,
				"allow_auto_merge":       true,
				"allow_update_branch":    true,
				"delete_branch_on_merge": true,
			})
		case r.Method == http.MethodPut && r.URL.Path == "/repos/example-org/example-repo/topics":
			decodeJSON(t, r, &replaceTopicsRequest)
			return jsonResponse(http.StatusOK, map[string]any{"names": []string{"anvil", "managed"}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			return nil
		}
	})

	visibility := "private"
	description := "Managed description"
	homepage := "https://example.com"
	defaultBranch := "main"

	spec := manifestv1alpha1.GitHubRepositorySpec{
		Owner:         "example-org",
		Name:          "example-repo",
		Visibility:    &visibility,
		Description:   &description,
		Homepage:      &homepage,
		DefaultBranch: &defaultBranch,
		Topics:        []string{"anvil", "managed"},
		Features: &manifestv1alpha1.GitHubRepositoryFeaturesSpec{
			HasIssues:   boolPtr(true),
			HasProjects: boolPtr(false),
			HasWiki:     boolPtr(true),
		},
		MergePolicy: &manifestv1alpha1.GitHubRepositoryMergePolicySpec{
			AllowSquashMerge:    boolPtr(false),
			AllowMergeCommit:    boolPtr(false),
			AllowRebaseMerge:    boolPtr(true),
			AllowAutoMerge:      boolPtr(true),
			AllowUpdateBranch:   boolPtr(true),
			DeleteBranchOnMerge: boolPtr(true),
		},
	}

	messages, err := newGitHubRepositoryPlan(spec).Run(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if patchRepositoryRequest["visibility"] != "private" {
		t.Fatalf("expected repository patch to include visibility, got %#v", patchRepositoryRequest)
	}
	if patchRepositoryRequest["homepage"] != "https://example.com" {
		t.Fatalf("expected repository patch to include homepage, got %#v", patchRepositoryRequest)
	}
	if !reflect.DeepEqual(replaceTopicsRequest["names"], []any{"anvil", "managed"}) {
		t.Fatalf("expected topics replace request, got %#v", replaceTopicsRequest)
	}

	expectedMessages := []string{
		"Reconciling GitHubRepository example-org/example-repo",
		"Updated repository settings for example-org/example-repo",
		"Updated repository topics for example-org/example-repo",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func TestPlanRunCreatesHCPTerraformWorkspace(t *testing.T) {
	t.Parallel()

	var createWorkspaceRequest map[string]any

	client := newHCPTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/organizations/example-org/workspaces/example-workspace":
			return jsonResponse(http.StatusNotFound, map[string]any{
				"errors": []map[string]any{{"detail": "not found"}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/organizations/example-org/workspaces":
			decodeJSON(t, r, &createWorkspaceRequest)
			return jsonResponse(http.StatusCreated, map[string]any{
				"data": map[string]any{
					"id":   "ws-123",
					"type": "workspaces",
					"attributes": map[string]any{
						"name":                           "example-workspace",
						"description":                    "Managed by Anvil",
						"terraform-version":              "1.14.8",
						"working-directory":              "terraform",
						"execution-mode":                 "remote",
						"allow-destroy-plan":             true,
						"assessments-enabled":            false,
						"auto-apply":                     false,
						"auto-apply-run-trigger":         false,
						"auto-destroy-at":                "",
						"auto-destroy-activity-duration": "",
						"file-triggers-enabled":          true,
						"global-remote-state":            false,
						"project-remote-state":           false,
						"queue-all-runs":                 false,
						"source-name":                    "",
						"source-url":                     "",
						"speculative-enabled":            true,
						"trigger-patterns":               []string{},
						"trigger-prefixes":               []string{},
					},
					"relationships": map[string]any{
						"project": map[string]any{
							"data": map[string]any{"id": "prj-123", "type": "projects"},
						},
					},
				},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			return nil
		}
	})

	description := "Managed by Anvil"
	terraformVersion := "1.14.8"
	workingDirectory := "terraform"
	executionMode := "remote"
	allowDestroyPlan := true
	speculativeEnabled := true
	fileTriggersEnabled := true
	projectID := "prj-123"

	spec := manifestv1alpha1.HCPTerraformWorkspaceSpec{
		Organization:        "example-org",
		Name:                "example-workspace",
		ProjectID:           &projectID,
		Description:         &description,
		TerraformVersion:    &terraformVersion,
		WorkingDirectory:    &workingDirectory,
		ExecutionMode:       &executionMode,
		AllowDestroyPlan:    &allowDestroyPlan,
		FileTriggersEnabled: &fileTriggersEnabled,
		SpeculativeEnabled:  &speculativeEnabled,
	}

	messages, err := newHCPTerraformWorkspacePlan(spec).Run(context.Background(), nil, client)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	attrs := createWorkspaceRequest["data"].(map[string]any)["attributes"].(map[string]any)
	if attrs["name"] != "example-workspace" {
		t.Fatalf("expected create request to include workspace name, got %#v", createWorkspaceRequest)
	}
	if attrs["terraform-version"] != "1.14.8" {
		t.Fatalf("expected create request to include terraform version, got %#v", createWorkspaceRequest)
	}

	expectedMessages := []string{
		"Reconciling HCPTerraformWorkspace example-org/example-workspace",
		"Created HCP Terraform workspace example-org/example-workspace",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func TestPlanRunReconcilesHCPTerraformWorkspaceDrift(t *testing.T) {
	t.Parallel()

	var updateWorkspaceRequest map[string]any
	var addTagsRequest map[string]any
	var removeTagsRequest map[string]any
	var replaceTagBindingsRequest map[string]any
	var addConsumersRequest map[string]any
	var removeConsumersRequest map[string]any
	var createVariableRequest map[string]any
	var updateVariableRequest map[string]any
	var addVariableSetCalls []string
	var removeVariableSetCalls []string

	client := newHCPTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/organizations/example-org/workspaces/example-workspace":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":   "ws-123",
					"type": "workspaces",
					"attributes": map[string]any{
						"name":                           "example-workspace",
						"description":                    "Old description",
						"terraform-version":              "1.5.0",
						"working-directory":              "",
						"execution-mode":                 "remote",
						"agent-pool-id":                  "",
						"allow-destroy-plan":             false,
						"assessments-enabled":            false,
						"auto-apply":                     false,
						"auto-apply-run-trigger":         false,
						"auto-destroy-at":                "",
						"auto-destroy-activity-duration": "",
						"file-triggers-enabled":          false,
						"global-remote-state":            false,
						"project-remote-state":           false,
						"queue-all-runs":                 true,
						"source-name":                    "",
						"source-url":                     "",
						"speculative-enabled":            false,
						"trigger-patterns":               []string{"legacy/**"},
						"trigger-prefixes":               []string{"legacy/"},
						"setting-overwrites":             map[string]any{"terraform-version": false},
						"vcs-repo":                       map[string]any{"identifier": "example/repo", "oauth-token-id": "ot-old", "branch": "main", "ingress-submodules": false, "tags-regex": ""},
					},
					"relationships": map[string]any{
						"project": map[string]any{
							"data": map[string]any{"id": "prj-old", "type": "projects"},
						},
					},
				},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v2/workspaces/ws-123":
			decodeJSON(t, r, &updateWorkspaceRequest)
			return jsonResponse(http.StatusOK, map[string]any{
				"data": map[string]any{
					"id":   "ws-123",
					"type": "workspaces",
					"attributes": map[string]any{
						"name":                           "example-workspace",
						"description":                    "Managed by Anvil",
						"terraform-version":              "1.14.8",
						"working-directory":              "terraform",
						"execution-mode":                 "remote",
						"allow-destroy-plan":             true,
						"assessments-enabled":            true,
						"auto-apply":                     true,
						"auto-apply-run-trigger":         false,
						"auto-destroy-at":                "",
						"auto-destroy-activity-duration": "",
						"file-triggers-enabled":          true,
						"global-remote-state":            false,
						"project-remote-state":           false,
						"queue-all-runs":                 false,
						"source-name":                    "anvil",
						"source-url":                     "https://example.com/anvil",
						"speculative-enabled":            true,
						"trigger-patterns":               []string{"terraform/**/*.tf"},
						"trigger-prefixes":               []string{"terraform/"},
						"setting-overwrites":             map[string]any{"terraform-version": true},
						"vcs-repo":                       map[string]any{"identifier": "example/repo", "oauth-token-id": "ot-123", "branch": "main", "ingress-submodules": false, "tags-regex": "^v.*$"},
					},
					"relationships": map[string]any{
						"project": map[string]any{
							"data": map[string]any{"id": "prj-123", "type": "projects"},
						},
					},
				},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/tags":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": []map[string]any{
					{"id": "tag-old", "type": "tags", "attributes": map[string]any{"name": "legacy"}},
					{"id": "tag-keep", "type": "tags", "attributes": map[string]any{"name": "platform"}},
				},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/tags":
			decodeJSON(t, r, &addTagsRequest)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/tags":
			decodeJSON(t, r, &removeTagsRequest)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/workspaces/ws-123/tag-bindings":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": []map[string]any{
					{"id": "tb-1", "type": "tag-bindings", "attributes": map[string]any{"key": "env", "value": "dev"}},
				},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v2/workspaces/ws-123/tag-bindings":
			decodeJSON(t, r, &replaceTagBindingsRequest)
			return jsonResponse(http.StatusOK, map[string]any{"data": []map[string]any{}})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/remote-state-consumers":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": []map[string]any{
					{"id": "ws-old", "type": "workspaces"},
				},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/remote-state-consumers":
			decodeJSON(t, r, &addConsumersRequest)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/workspaces/ws-123/relationships/remote-state-consumers":
			decodeJSON(t, r, &removeConsumersRequest)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/workspaces/ws-123/vars":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": []map[string]any{
					{"id": "var-old", "type": "vars", "attributes": map[string]any{"key": "AWS_REGION", "value": "us-west-2", "description": "", "category": "env", "hcl": false, "sensitive": false}},
					{"id": "var-stale", "type": "vars", "attributes": map[string]any{"key": "LEGACY", "value": "true", "description": "", "category": "env", "hcl": false, "sensitive": false}},
				},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v2/workspaces/ws-123/vars/var-old":
			decodeJSON(t, r, &updateVariableRequest)
			return jsonResponse(http.StatusOK, map[string]any{"data": map[string]any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/workspaces/ws-123/vars":
			decodeJSON(t, r, &createVariableRequest)
			return jsonResponse(http.StatusCreated, map[string]any{"data": map[string]any{}})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/workspaces/ws-123/vars/var-stale":
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/workspaces/ws-123/varsets":
			return jsonResponse(http.StatusOK, map[string]any{
				"data": []map[string]any{
					{"id": "varset-old", "type": "varsets"},
				},
			})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v2/varsets/") && strings.HasSuffix(r.URL.Path, "/relationships/workspaces"):
			addVariableSetCalls = append(addVariableSetCalls, r.URL.Path)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v2/varsets/") && strings.HasSuffix(r.URL.Path, "/relationships/workspaces"):
			removeVariableSetCalls = append(removeVariableSetCalls, r.URL.Path)
			return jsonResponse(http.StatusNoContent, nil)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			return nil
		}
	})

	description := "Managed by Anvil"
	terraformVersion := "1.14.8"
	workingDirectory := "terraform"
	executionMode := "remote"
	allowDestroyPlan := true
	assessmentsEnabled := true
	autoApply := true
	fileTriggersEnabled := true
	queueAllRuns := false
	sourceName := "anvil"
	sourceURL := "https://example.com/anvil"
	speculativeEnabled := true
	projectID := "prj-123"
	oauthTokenID := "ot-123"
	branch := "main"
	tagsRegex := "^v.*$"

	spec := manifestv1alpha1.HCPTerraformWorkspaceSpec{
		Organization:           "example-org",
		Name:                   "example-workspace",
		ProjectID:              &projectID,
		Description:            &description,
		TerraformVersion:       &terraformVersion,
		WorkingDirectory:       &workingDirectory,
		ExecutionMode:          &executionMode,
		AllowDestroyPlan:       &allowDestroyPlan,
		AssessmentsEnabled:     &assessmentsEnabled,
		AutoApply:              &autoApply,
		FileTriggersEnabled:    &fileTriggersEnabled,
		QueueAllRuns:           &queueAllRuns,
		SourceName:             &sourceName,
		SourceURL:              &sourceURL,
		SpeculativeEnabled:     &speculativeEnabled,
		Tags:                   []string{"platform", "managed"},
		TagBindings:            []manifestv1alpha1.HCPTerraformWorkspaceTagBindingSpec{{Key: "env", Value: "prod"}},
		TriggerPatterns:        []string{"terraform/**/*.tf"},
		TriggerPrefixes:        []string{"terraform/"},
		RemoteStateConsumerIDs: []string{"ws-new"},
		VCSRepo: &manifestv1alpha1.HCPTerraformWorkspaceVCSRepoSpec{
			Identifier:   stringPtr("example/repo"),
			OAuthTokenID: &oauthTokenID,
			Branch:       &branch,
			TagsRegex:    &tagsRegex,
		},
		Variables: []manifestv1alpha1.HCPTerraformWorkspaceVariableSpec{
			{Key: "AWS_REGION", Category: "env", Value: "us-east-1"},
			{Key: "account_id", Category: "terraform", Value: "\"123456789012\"", HCL: boolPtr(true)},
		},
		VariableSetIDs: []string{"varset-new"},
	}

	messages, err := newHCPTerraformWorkspacePlan(spec).Run(context.Background(), nil, client)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	attrs := updateWorkspaceRequest["data"].(map[string]any)["attributes"].(map[string]any)
	if attrs["terraform-version"] != "1.14.8" {
		t.Fatalf("expected workspace patch to include terraform-version, got %#v", updateWorkspaceRequest)
	}
	if attrs["source-name"] != "anvil" {
		t.Fatalf("expected workspace patch to include source-name, got %#v", updateWorkspaceRequest)
	}

	addedTags := addTagsRequest["data"].([]any)
	if len(addedTags) != 1 {
		t.Fatalf("expected one tag add request, got %#v", addTagsRequest)
	}

	removedTags := removeTagsRequest["data"].([]any)
	if len(removedTags) != 1 {
		t.Fatalf("expected one tag removal request, got %#v", removeTagsRequest)
	}

	createdVarAttrs := createVariableRequest["data"].(map[string]any)["attributes"].(map[string]any)
	if createdVarAttrs["key"] != "account_id" {
		t.Fatalf("expected variable create request for account_id, got %#v", createVariableRequest)
	}

	updatedVarAttrs := updateVariableRequest["data"].(map[string]any)["attributes"].(map[string]any)
	if updatedVarAttrs["value"] != "us-east-1" {
		t.Fatalf("expected variable update request to set us-east-1, got %#v", updateVariableRequest)
	}

	if !containsPath(addVariableSetCalls, "/api/v2/varsets/varset-new/relationships/workspaces") {
		t.Fatalf("expected varset-new assignment, got %v", addVariableSetCalls)
	}
	if !containsPath(removeVariableSetCalls, "/api/v2/varsets/varset-old/relationships/workspaces") {
		t.Fatalf("expected varset-old removal, got %v", removeVariableSetCalls)
	}

	expectedMessages := []string{
		"Reconciling HCPTerraformWorkspace example-org/example-workspace",
		"Updated workspace settings for example-org/example-workspace",
		"Updated workspace tags for example-org/example-workspace",
		"Updated workspace tag bindings for example-org/example-workspace",
		"Updated remote state consumers for example-org/example-workspace",
		"Updated workspace variables for example-org/example-workspace",
		"Updated workspace variable set assignments for example-org/example-workspace",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func TestPlanRunRejectsUnsupportedHCPTerraformWorkspaceSurfaceClearly(t *testing.T) {
	t.Parallel()

	spec := manifestv1alpha1.HCPTerraformWorkspaceSpec{
		Organization: "example-org",
		Name:         "example-workspace",
		Notifications: []manifestv1alpha1.HCPTerraformWorkspaceNotificationSpec{
			{Name: "run-events", DestinationType: "generic", URL: stringPtr("https://hooks.example.com")},
		},
	}

	messages, err := newHCPTerraformWorkspacePlan(spec).Run(context.Background(), nil, newHCPTestClient(t, func(r *http.Request) *http.Response {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		return nil
	}))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "spec.notifications is not supported") {
		t.Fatalf("expected unsupported notifications error, got %v", err)
	}

	expectedMessages := []string{
		"Reconciling HCPTerraformWorkspace example-org/example-workspace",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func newGitHubRepositoryPlan(spec manifestv1alpha1.GitHubRepositorySpec) Plan {
	return Plan{
		githubRepositories: []manifest.LoadedGitHubRepositoryManifest{
			{
				Path: "repo.yaml",
				Manifest: manifestv1alpha1.NewGitHubRepositoryManifest(
					manifestv1alpha1.Metadata{Name: spec.Name},
					spec,
				),
			},
		},
	}
}

func newHCPTerraformWorkspacePlan(spec manifestv1alpha1.HCPTerraformWorkspaceSpec) Plan {
	return Plan{
		hcpTerraformWorkspaces: []manifest.LoadedHCPTerraformWorkspaceManifest{
			{
				Path: "workspace.yaml",
				Manifest: manifestv1alpha1.NewHCPTerraformWorkspaceManifest(
					manifestv1alpha1.Metadata{Name: spec.Name},
					spec,
				),
			},
		},
	}
}

func newGitHubTestClient(t *testing.T, responder func(*http.Request) *http.Response) *ghapi.Client {
	t.Helper()
	return ghapi.NewClient("https://api.github.test", "test-token", &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return responder(r), nil
		}),
	})
}

func newHCPTestClient(t *testing.T, responder func(*http.Request) *http.Response) *hcpapi.Client {
	t.Helper()
	return hcpapi.NewClient("https://app.terraform.test/api/v2", "test-token", &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return responder(r), nil
		}),
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(status int, body any) *http.Response {
	var reader io.ReadCloser
	if body == nil {
		reader = io.NopCloser(bytes.NewReader(nil))
	} else {
		payload, _ := json.Marshal(body)
		reader = io.NopCloser(bytes.NewReader(payload))
	}

	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       reader,
	}
}

func decodeJSON(t *testing.T, r *http.Request, target any) {
	t.Helper()

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func containsPath(values []string, path string) bool {
	for _, value := range values {
		if value == path {
			return true
		}
	}

	return false
}
