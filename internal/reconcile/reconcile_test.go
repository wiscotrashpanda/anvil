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
	"github.com/emkaytec/anvil/internal/manifest"
)

func TestPlanRunCreatesRepository(t *testing.T) {
	t.Parallel()

	var createdRequest map[string]any

	client := newTestClient(t, func(r *http.Request) *http.Response {
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
				"archived":               false,
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

	messages, err := newTestPlan(spec).Run(context.Background(), client)
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
	var createPagesRequest map[string]any
	var updatePropertiesRequest map[string]any
	var updateBranchProtectionRequest map[string]any
	var deletedProtectionPaths []string

	client := newTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusOK, map[string]any{
				"name":                   "example-repo",
				"full_name":              "example-org/example-repo",
				"visibility":             "public",
				"description":            "Old description",
				"homepage":               "",
				"archived":               false,
				"default_branch":         "main",
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
				"security_and_analysis": map[string]any{
					"advanced_security": map[string]any{"status": "disabled"},
				},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/repos/example-org/example-repo":
			decodeJSON(t, r, &patchRepositoryRequest)
			return jsonResponse(http.StatusOK, map[string]any{
				"name":                        "example-repo",
				"full_name":                   "example-org/example-repo",
				"visibility":                  "private",
				"description":                 "Managed description",
				"homepage":                    "https://example.com",
				"archived":                    false,
				"default_branch":              "main",
				"topics":                      []string{"legacy"},
				"owner":                       map[string]any{"login": "example-org", "type": "Organization"},
				"has_issues":                  true,
				"has_projects":                false,
				"has_wiki":                    true,
				"allow_squash_merge":          false,
				"allow_merge_commit":          false,
				"allow_rebase_merge":          true,
				"allow_auto_merge":            true,
				"allow_update_branch":         true,
				"delete_branch_on_merge":      true,
				"squash_merge_commit_title":   "PR_TITLE",
				"squash_merge_commit_message": "PR_BODY",
				"merge_commit_title":          "PR_TITLE",
				"merge_commit_message":        "PR_BODY",
				"security_and_analysis": map[string]any{
					"advanced_security": map[string]any{"status": "enabled"},
				},
			})
		case r.Method == http.MethodPut && r.URL.Path == "/repos/example-org/example-repo/topics":
			decodeJSON(t, r, &replaceTopicsRequest)
			return jsonResponse(http.StatusOK, map[string]any{"names": []string{"anvil", "managed"}})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/pages":
			return jsonResponse(http.StatusNotFound, map[string]any{"message": "Not Found"})
		case r.Method == http.MethodPost && r.URL.Path == "/repos/example-org/example-repo/pages":
			decodeJSON(t, r, &createPagesRequest)
			return jsonResponse(http.StatusCreated, map[string]any{
				"build_type":     "legacy",
				"https_enforced": true,
				"source":         map[string]any{"branch": "main", "path": "/docs"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/properties/values":
			return jsonResponse(http.StatusOK, []map[string]any{
				{"property_name": "service", "value": "old-service"},
				{"property_name": "legacy", "value": "remove-me"},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/repos/example-org/example-repo/properties/values":
			decodeJSON(t, r, &updatePropertiesRequest)
			return jsonResponse(http.StatusNoContent, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches" && strings.Contains(r.URL.RawQuery, "protected=true"):
			return jsonResponse(http.StatusOK, []map[string]any{
				{"name": "main"},
				{"name": "release"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches/main":
			return jsonResponse(http.StatusOK, map[string]any{"name": "main"})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches/main/protection":
			return jsonResponse(http.StatusNotFound, map[string]any{"message": "Not Found"})
		case r.Method == http.MethodPut && r.URL.Path == "/repos/example-org/example-repo/branches/main/protection":
			decodeJSON(t, r, &updateBranchProtectionRequest)
			return jsonResponse(http.StatusOK, map[string]any{
				"required_status_checks": map[string]any{
					"strict": true,
					"checks": []map[string]any{{"context": "ci/test"}},
				},
				"enforce_admins": map[string]any{"enabled": true},
				"required_pull_request_reviews": map[string]any{
					"dismissal_restrictions":          map[string]any{"users": []any{}, "teams": []any{}, "apps": []any{}},
					"dismiss_stale_reviews":           true,
					"require_code_owner_reviews":      true,
					"required_approving_review_count": 1,
					"require_last_push_approval":      false,
				},
				"bypass_pull_request_allowances": map[string]any{"users": []any{}, "teams": []any{}, "apps": []any{}},
				"required_linear_history":        map[string]any{"enabled": true},
			})
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/repos/example-org/example-repo/branches/") && strings.HasSuffix(r.URL.Path, "/protection"):
			deletedProtectionPaths = append(deletedProtectionPaths, r.URL.Path)
			return jsonResponse(http.StatusNoContent, nil)
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			return nil
		}
	})

	visibility := "private"
	description := "Managed description"
	homepage := "https://example.com"
	buildType := "legacy"
	squashTitle := "PR_TITLE"
	squashMessage := "PR_BODY"
	mergeTitle := "PR_TITLE"
	mergeMessage := "PR_BODY"

	spec := manifestv1alpha1.GitHubRepositorySpec{
		Owner:       "example-org",
		Name:        "example-repo",
		Visibility:  &visibility,
		Description: &description,
		Homepage:    &homepage,
		Topics:      []string{"anvil", "managed"},
		Features:    &manifestv1alpha1.GitHubRepositoryFeaturesSpec{HasIssues: boolPtr(true), HasProjects: boolPtr(false), HasWiki: boolPtr(true)},
		MergePolicy: &manifestv1alpha1.GitHubRepositoryMergePolicySpec{AllowSquashMerge: boolPtr(false), AllowMergeCommit: boolPtr(false), AllowRebaseMerge: boolPtr(true), AllowAutoMerge: boolPtr(true), AllowUpdateBranch: boolPtr(true), DeleteBranchOnMerge: boolPtr(true), SquashMergeCommitTitle: &squashTitle, SquashMergeCommitMessage: &squashMessage, MergeCommitTitle: &mergeTitle, MergeCommitMessage: &mergeMessage},
		SecurityAndAnalysis: &manifestv1alpha1.GitHubRepositorySecurityAndAnalysisSpec{
			AdvancedSecurity: &manifestv1alpha1.GitHubRepositorySecuritySettingSpec{Status: "enabled"},
		},
		Pages: &manifestv1alpha1.GitHubRepositoryPagesSpec{
			BuildType:     &buildType,
			HTTPSEnforced: boolPtr(true),
			Source:        &manifestv1alpha1.GitHubRepositoryPagesSourceSpec{Branch: "main", Path: "/docs"},
		},
		CustomProperties: []manifestv1alpha1.GitHubRepositoryCustomPropertySpec{
			{Name: "service", Value: "anvil"},
		},
		Branches: []manifestv1alpha1.GitHubRepositoryBranchSpec{
			{
				Name: "main",
				Protection: &manifestv1alpha1.GitHubRepositoryBranchProtectionSpec{
					RequiredStatusChecks: &manifestv1alpha1.GitHubRequiredStatusChecksSpec{
						Strict: true,
						Checks: []manifestv1alpha1.GitHubRequiredStatusCheckSpec{{Context: "ci/test"}},
					},
					EnforceAdmins:         boolPtr(true),
					RequiredLinearHistory: boolPtr(true),
					PullRequestReviews: &manifestv1alpha1.GitHubPullRequestReviewsSpec{
						DismissStaleReviews:          boolPtr(true),
						RequireCodeOwnerReviews:      boolPtr(true),
						RequiredApprovingReviewCount: intPtr(1),
					},
				},
			},
		},
	}

	messages, err := newTestPlan(spec).Run(context.Background(), client)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if patchRepositoryRequest["visibility"] != "private" {
		t.Fatalf("expected repository patch to include visibility, got %#v", patchRepositoryRequest)
	}
	if securityAndAnalysis, ok := patchRepositoryRequest["security_and_analysis"]; ok && securityAndAnalysis != nil {
		t.Fatalf("expected repository patch to omit security_and_analysis for public repos, got %#v", patchRepositoryRequest)
	}
	if patchRepositoryRequest["homepage"] != "https://example.com" {
		t.Fatalf("expected repository patch to include homepage, got %#v", patchRepositoryRequest)
	}
	if !reflect.DeepEqual(replaceTopicsRequest["names"], []any{"anvil", "managed"}) {
		t.Fatalf("expected topics replace request, got %#v", replaceTopicsRequest)
	}
	if createPagesRequest["build_type"] != "legacy" {
		t.Fatalf("expected pages create request to include build_type, got %#v", createPagesRequest)
	}

	properties := propertyValuesFromRequest(t, updatePropertiesRequest)
	expectedProperties := map[string]any{
		"service": "anvil",
		"legacy":  nil,
	}
	if !reflect.DeepEqual(properties, expectedProperties) {
		t.Fatalf("expected property update request %v, got %v", expectedProperties, properties)
	}

	if updateBranchProtectionRequest["required_linear_history"] != true {
		t.Fatalf("expected branch protection request to require linear history, got %#v", updateBranchProtectionRequest)
	}
	if !containsPath(deletedProtectionPaths, "/repos/example-org/example-repo/branches/release/protection") {
		t.Fatalf("expected release protection to be cleared, got %v", deletedProtectionPaths)
	}

	expectedMessages := []string{
		"Reconciling GitHubRepository example-org/example-repo",
		"Updated repository settings for example-org/example-repo",
		"Updated repository topics for example-org/example-repo",
		"Updated GitHub Pages settings for example-org/example-repo",
		"Updated custom properties for example-org/example-repo",
		"Updated branch protection for example-org/example-repo#main",
		"Cleared branch protection for example-org/example-repo#release",
	}
	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Fatalf("expected messages %v, got %v", expectedMessages, messages)
	}
}

func TestBuildSecurityAndAnalysisUpdateIncludesAdvancedSecurityForPrivateRepos(t *testing.T) {
	t.Parallel()

	repository := &ghapi.Repository{
		Visibility: "private",
		SecurityAndAnalysis: &ghapi.SecurityAndAnalysis{
			AdvancedSecurity: &ghapi.SecuritySetting{Status: "disabled"},
		},
	}

	spec := manifestv1alpha1.GitHubRepositorySpec{
		SecurityAndAnalysis: &manifestv1alpha1.GitHubRepositorySecurityAndAnalysisSpec{
			AdvancedSecurity: &manifestv1alpha1.GitHubRepositorySecuritySettingSpec{Status: "enabled"},
		},
	}

	update := buildSecurityAndAnalysisUpdate(repository, spec)
	if update == nil || update.AdvancedSecurity == nil {
		t.Fatalf("expected advanced security update for private repo, got %#v", update)
	}
	if update.AdvancedSecurity.Status != "enabled" {
		t.Fatalf("expected advanced security status enabled, got %#v", update.AdvancedSecurity)
	}
}

func TestBuildSecurityAndAnalysisUpdateOmitsAdvancedSecurityForPublicRepos(t *testing.T) {
	t.Parallel()

	repository := &ghapi.Repository{
		Visibility: "public",
		SecurityAndAnalysis: &ghapi.SecurityAndAnalysis{
			AdvancedSecurity: &ghapi.SecuritySetting{Status: "disabled"},
		},
	}

	spec := manifestv1alpha1.GitHubRepositorySpec{
		SecurityAndAnalysis: &manifestv1alpha1.GitHubRepositorySecurityAndAnalysisSpec{
			AdvancedSecurity: &manifestv1alpha1.GitHubRepositorySecuritySettingSpec{Status: "enabled"},
		},
	}

	update := buildSecurityAndAnalysisUpdate(repository, spec)
	if update != nil {
		t.Fatalf("expected no advanced security update for public repo, got %#v", update)
	}
}

func TestActorAllowanceStateToRequestUsesEmptyArrays(t *testing.T) {
	t.Parallel()

	request := actorAllowanceState{}.toRequest()

	users, ok := request["users"].([]string)
	if !ok {
		t.Fatalf("expected users to be []string, got %#v", request["users"])
	}
	if users == nil || len(users) != 0 {
		t.Fatalf("expected users to be an empty non-nil slice, got %#v", users)
	}

	teams, ok := request["teams"].([]string)
	if !ok {
		t.Fatalf("expected teams to be []string, got %#v", request["teams"])
	}
	if teams == nil || len(teams) != 0 {
		t.Fatalf("expected teams to be an empty non-nil slice, got %#v", teams)
	}

	apps, ok := request["apps"].([]string)
	if !ok {
		t.Fatalf("expected apps to be []string, got %#v", request["apps"])
	}
	if apps == nil || len(apps) != 0 {
		t.Fatalf("expected apps to be an empty non-nil slice, got %#v", apps)
	}
}

func TestPlanRunReportsUnsupportedBranchProtectionClearly(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusOK, map[string]any{
				"name":            "example-repo",
				"full_name":       "example-org/example-repo",
				"visibility":      "private",
				"default_branch":  "main",
				"owner":           map[string]any{"login": "example-org", "type": "Organization"},
				"topics":          []string{},
				"homepage":        "",
				"has_issues":      true,
				"has_projects":    true,
				"has_wiki":        true,
				"allow_squash_merge": true,
				"allow_merge_commit": true,
				"allow_rebase_merge": true,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches" && strings.Contains(r.URL.RawQuery, "protected=true"):
			return jsonResponse(http.StatusOK, []map[string]any{})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches/main":
			return jsonResponse(http.StatusOK, map[string]any{"name": "main"})
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo/branches/main/protection":
			return jsonResponse(http.StatusForbidden, map[string]any{"message": "Upgrade to GitHub Pro or make this repository public to enable this feature."})
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			return nil
		}
	})

	spec := manifestv1alpha1.GitHubRepositorySpec{
		Owner: "example-org",
		Name:  "example-repo",
		Branches: []manifestv1alpha1.GitHubRepositoryBranchSpec{
			{
				Name: "main",
				Protection: &manifestv1alpha1.GitHubRepositoryBranchProtectionSpec{
					EnforceAdmins: boolPtr(true),
				},
			},
		},
	}

	_, err := newTestPlan(spec).Run(context.Background(), client)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "branch protection is unavailable for example-org/example-repo") {
		t.Fatalf("expected clearer branch protection availability error, got %v", err)
	}
}

func TestPlanRunReportsUnsupportedAdvancedSecurityClearly(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) *http.Response {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusOK, map[string]any{
				"name":                    "example-repo",
				"full_name":               "example-org/example-repo",
				"visibility":              "private",
				"default_branch":          "main",
				"owner":                   map[string]any{"login": "example-org", "type": "Organization"},
				"topics":                  []string{},
				"homepage":                "",
				"has_issues":              true,
				"has_projects":            true,
				"has_wiki":                true,
				"allow_squash_merge":      true,
				"allow_merge_commit":      true,
				"allow_rebase_merge":      true,
				"allow_auto_merge":        false,
				"allow_update_branch":     false,
				"delete_branch_on_merge":  false,
				"security_and_analysis":   map[string]any{"advanced_security": map[string]any{"status": "disabled"}},
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/repos/example-org/example-repo":
			return jsonResponse(http.StatusUnprocessableEntity, map[string]any{"message": "Updating Advanced Security on this repository is not available, nor a pre-requisite for security features."})
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			return nil
		}
	})

	spec := manifestv1alpha1.GitHubRepositorySpec{
		Owner: "example-org",
		Name:  "example-repo",
		SecurityAndAnalysis: &manifestv1alpha1.GitHubRepositorySecurityAndAnalysisSpec{
			AdvancedSecurity: &manifestv1alpha1.GitHubRepositorySecuritySettingSpec{Status: "enabled"},
		},
	}

	_, err := newTestPlan(spec).Run(context.Background(), client)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "remove securityAndAnalysis.advancedSecurity from the manifest") {
		t.Fatalf("expected clearer advanced security availability error, got %v", err)
	}
}

func newTestPlan(spec manifestv1alpha1.GitHubRepositorySpec) Plan {
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

func decodeJSON(t *testing.T, r *http.Request, target any) {
	t.Helper()

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func propertyValuesFromRequest(t *testing.T, request map[string]any) map[string]any {
	t.Helper()

	rawProperties, ok := request["properties"].([]any)
	if !ok {
		t.Fatalf("expected properties request array, got %#v", request["properties"])
	}

	properties := make(map[string]any, len(rawProperties))
	for _, raw := range rawProperties {
		property, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("expected property entry object, got %#v", raw)
		}
		name, ok := property["property_name"].(string)
		if !ok {
			t.Fatalf("expected property_name string, got %#v", property["property_name"])
		}
		properties[strings.ToLower(name)] = property["value"]
	}

	return properties
}

func containsPath(values []string, path string) bool {
	for _, value := range values {
		if value == path {
			return true
		}
	}
	return false
}

func newTestClient(t *testing.T, responder func(*http.Request) *http.Response) *ghapi.Client {
	t.Helper()

	httpClient := &http.Client{
		Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			response := responder(r)
			if response == nil {
				t.Fatal("expected responder to return a response")
			}
			return response, nil
		}),
	}

	return ghapi.NewClient("https://api.github.test", "test-token", httpClient)
}

func jsonResponse(status int, payload any) *http.Response {
	response := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
	}

	if payload == nil {
		response.Body = io.NopCloser(bytes.NewReader(nil))
		return response
	}

	body, _ := json.Marshal(payload)
	response.Header.Set("Content-Type", "application/json")
	response.Body = io.NopCloser(bytes.NewReader(body))
	return response
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
