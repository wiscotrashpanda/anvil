package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	DefaultBaseURL    = "https://api.github.com"
	defaultAPIVersion = "2022-11-28"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type APIError struct {
	StatusCode int
	Method     string
	Path       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("github api %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Message)
	}

	return fmt.Sprintf("github api %s %s returned %d", e.Method, e.Path, e.StatusCode)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

func NewClientFromEnv() (*Client, error) {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GH_TOKEN"))
	}

	if token == "" {
		return nil, fmt.Errorf("missing GitHub token: set GITHUB_TOKEN or GH_TOKEN")
	}

	return NewClient(DefaultBaseURL, token, nil), nil
}

func NewClient(baseURL string, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
		token:      token,
	}
}

type Account struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Repository struct {
	Name                     string               `json:"name"`
	FullName                 string               `json:"full_name"`
	Visibility               string               `json:"visibility"`
	Description              *string              `json:"description"`
	Homepage                 string               `json:"homepage"`
	Archived                 bool                 `json:"archived"`
	DefaultBranch            string               `json:"default_branch"`
	Topics                   []string             `json:"topics"`
	Owner                    Account              `json:"owner"`
	HasIssues                bool                 `json:"has_issues"`
	HasProjects              bool                 `json:"has_projects"`
	HasWiki                  bool                 `json:"has_wiki"`
	AllowSquashMerge         bool                 `json:"allow_squash_merge"`
	AllowMergeCommit         bool                 `json:"allow_merge_commit"`
	AllowRebaseMerge         bool                 `json:"allow_rebase_merge"`
	AllowAutoMerge           bool                 `json:"allow_auto_merge"`
	AllowUpdateBranch        bool                 `json:"allow_update_branch"`
	DeleteBranchOnMerge      bool                 `json:"delete_branch_on_merge"`
	SquashMergeCommitTitle   string               `json:"squash_merge_commit_title"`
	SquashMergeCommitMessage string               `json:"squash_merge_commit_message"`
	MergeCommitTitle         string               `json:"merge_commit_title"`
	MergeCommitMessage       string               `json:"merge_commit_message"`
	SecurityAndAnalysis      *SecurityAndAnalysis `json:"security_and_analysis"`
}

type SecurityAndAnalysis struct {
	AdvancedSecurity                      *SecuritySetting `json:"advanced_security,omitempty"`
	CodeSecurity                          *SecuritySetting `json:"code_security,omitempty"`
	SecretScanning                        *SecuritySetting `json:"secret_scanning,omitempty"`
	SecretScanningPushProtection          *SecuritySetting `json:"secret_scanning_push_protection,omitempty"`
	SecretScanningAIDetection             *SecuritySetting `json:"secret_scanning_ai_detection,omitempty"`
	SecretScanningNonProviderPatterns     *SecuritySetting `json:"secret_scanning_non_provider_patterns,omitempty"`
	SecretScanningDelegatedAlertDismissal *SecuritySetting `json:"secret_scanning_delegated_alert_dismissal,omitempty"`
	SecretScanningDelegatedBypass         *SecuritySetting `json:"secret_scanning_delegated_bypass,omitempty"`
}

type SecuritySetting struct {
	Status string `json:"status"`
}

type CreateRepositoryRequest struct {
	Name                string  `json:"name"`
	Visibility          *string `json:"visibility,omitempty"`
	Description         *string `json:"description,omitempty"`
	Homepage            *string `json:"homepage,omitempty"`
	AutoInit            bool    `json:"auto_init,omitempty"`
	GitignoreTemplate   string  `json:"gitignore_template,omitempty"`
	LicenseTemplate     string  `json:"license_template,omitempty"`
	IsTemplate          *bool   `json:"is_template,omitempty"`
	HasIssues           *bool   `json:"has_issues,omitempty"`
	HasProjects         *bool   `json:"has_projects,omitempty"`
	HasWiki             *bool   `json:"has_wiki,omitempty"`
	AllowSquashMerge    *bool   `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit    *bool   `json:"allow_merge_commit,omitempty"`
	AllowRebaseMerge    *bool   `json:"allow_rebase_merge,omitempty"`
	AllowAutoMerge      *bool   `json:"allow_auto_merge,omitempty"`
	DeleteBranchOnMerge *bool   `json:"delete_branch_on_merge,omitempty"`
}

type UpdateRepositoryRequest struct {
	Visibility               *string              `json:"visibility,omitempty"`
	Description              *string              `json:"description,omitempty"`
	Homepage                 *string              `json:"homepage,omitempty"`
	Archived                 *bool                `json:"archived,omitempty"`
	DefaultBranch            *string              `json:"default_branch,omitempty"`
	HasIssues                *bool                `json:"has_issues,omitempty"`
	HasProjects              *bool                `json:"has_projects,omitempty"`
	HasWiki                  *bool                `json:"has_wiki,omitempty"`
	AllowSquashMerge         *bool                `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit         *bool                `json:"allow_merge_commit,omitempty"`
	AllowRebaseMerge         *bool                `json:"allow_rebase_merge,omitempty"`
	AllowAutoMerge           *bool                `json:"allow_auto_merge,omitempty"`
	AllowUpdateBranch        *bool                `json:"allow_update_branch,omitempty"`
	DeleteBranchOnMerge      *bool                `json:"delete_branch_on_merge,omitempty"`
	SquashMergeCommitTitle   *string              `json:"squash_merge_commit_title,omitempty"`
	SquashMergeCommitMessage *string              `json:"squash_merge_commit_message,omitempty"`
	MergeCommitTitle         *string              `json:"merge_commit_title,omitempty"`
	MergeCommitMessage       *string              `json:"merge_commit_message,omitempty"`
	SecurityAndAnalysis      *SecurityAndAnalysis `json:"security_and_analysis,omitempty"`
}

func (r UpdateRepositoryRequest) IsZero() bool {
	return r.Visibility == nil &&
		r.Description == nil &&
		r.Homepage == nil &&
		r.Archived == nil &&
		r.DefaultBranch == nil &&
		r.HasIssues == nil &&
		r.HasProjects == nil &&
		r.HasWiki == nil &&
		r.AllowSquashMerge == nil &&
		r.AllowMergeCommit == nil &&
		r.AllowRebaseMerge == nil &&
		r.AllowAutoMerge == nil &&
		r.AllowUpdateBranch == nil &&
		r.DeleteBranchOnMerge == nil &&
		r.SquashMergeCommitTitle == nil &&
		r.SquashMergeCommitMessage == nil &&
		r.MergeCommitTitle == nil &&
		r.MergeCommitMessage == nil &&
		r.SecurityAndAnalysis == nil
}

type TopicsResponse struct {
	Names []string `json:"names"`
}

type PagesSite struct {
	BuildType     string       `json:"build_type"`
	CNAME         *string      `json:"cname"`
	HTTPSEnforced bool         `json:"https_enforced"`
	Source        *PagesSource `json:"source"`
}

type PagesSource struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

type UpdatePagesRequest struct {
	BuildType     *string      `json:"build_type,omitempty"`
	CNAME         any          `json:"cname,omitempty"`
	HTTPSEnforced *bool        `json:"https_enforced,omitempty"`
	Source        *PagesSource `json:"source,omitempty"`
}

type CustomPropertyValue struct {
	PropertyName string `json:"property_name"`
	Value        any    `json:"value"`
}

type Branch struct {
	Name string `json:"name"`
}

type BranchProtection struct {
	RequiredStatusChecks        *RequiredStatusChecks `json:"required_status_checks"`
	EnforceAdmins               *EnabledSetting       `json:"enforce_admins"`
	RequiredPullRequestReviews  *PullRequestReviews   `json:"required_pull_request_reviews"`
	Restrictions                *ActorAllowance       `json:"restrictions"`
	BypassPullRequestAllowances *ActorAllowance       `json:"bypass_pull_request_allowances"`
	RequiredLinearHistory       *EnabledSetting       `json:"required_linear_history"`
}

type RequiredStatusChecks struct {
	Strict bool                  `json:"strict"`
	Checks []RequiredStatusCheck `json:"checks"`
}

type RequiredStatusCheck struct {
	Context string `json:"context"`
	Name    string `json:"name"`
	AppID   *int64 `json:"app_id"`
}

func (c RequiredStatusCheck) CheckContext() string {
	if c.Context != "" {
		return c.Context
	}

	return c.Name
}

type EnabledSetting struct {
	Enabled bool `json:"enabled"`
}

type PullRequestReviews struct {
	DismissalRestrictions        *ActorAllowance `json:"dismissal_restrictions"`
	DismissStaleReviews          bool            `json:"dismiss_stale_reviews"`
	RequireCodeOwnerReviews      bool            `json:"require_code_owner_reviews"`
	RequiredApprovingReviewCount int             `json:"required_approving_review_count"`
	RequireLastPushApproval      bool            `json:"require_last_push_approval"`
	BypassPullRequestAllowances  *ActorAllowance `json:"bypass_pull_request_allowances"`
}

type ActorAllowance struct {
	Users []UserActor `json:"users"`
	Teams []TeamActor `json:"teams"`
	Apps  []AppActor  `json:"apps"`
}

type UserActor struct {
	Login string `json:"login"`
}

type TeamActor struct {
	Slug string `json:"slug"`
}

type AppActor struct {
	Slug string `json:"slug"`
}

func (c *Client) GetRepository(ctx context.Context, owner string, repo string) (*Repository, error) {
	var repository Repository
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s", url.PathEscape(owner), url.PathEscape(repo)), nil, &repository, http.StatusOK); err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) GetAccount(ctx context.Context, owner string) (*Account, error) {
	var account Account
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/users/%s", url.PathEscape(owner)), nil, &account, http.StatusOK); err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) CreateOrganizationRepository(ctx context.Context, org string, request CreateRepositoryRequest) (*Repository, error) {
	var repository Repository
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/orgs/%s/repos", url.PathEscape(org)), request, &repository, http.StatusCreated); err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) CreateUserRepository(ctx context.Context, request CreateRepositoryRequest) (*Repository, error) {
	var repository Repository
	if err := c.request(ctx, http.MethodPost, "/user/repos", request, &repository, http.StatusCreated); err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) UpdateRepository(ctx context.Context, owner string, repo string, request UpdateRepositoryRequest) (*Repository, error) {
	var repository Repository
	if err := c.request(ctx, http.MethodPatch, fmt.Sprintf("/repos/%s/%s", url.PathEscape(owner), url.PathEscape(repo)), request, &repository, http.StatusOK); err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) ReplaceTopics(ctx context.Context, owner string, repo string, topics []string) error {
	request := TopicsResponse{Names: topics}
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/repos/%s/%s/topics", url.PathEscape(owner), url.PathEscape(repo)), request, nil, http.StatusOK)
}

func (c *Client) GetPages(ctx context.Context, owner string, repo string) (*PagesSite, error) {
	var pages PagesSite
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/pages", url.PathEscape(owner), url.PathEscape(repo)), nil, &pages, http.StatusOK); err != nil {
		return nil, err
	}

	return &pages, nil
}

func (c *Client) CreatePages(ctx context.Context, owner string, repo string, request UpdatePagesRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/pages", url.PathEscape(owner), url.PathEscape(repo)), request, nil, http.StatusCreated, http.StatusOK)
}

func (c *Client) UpdatePages(ctx context.Context, owner string, repo string, request UpdatePagesRequest) error {
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/repos/%s/%s/pages", url.PathEscape(owner), url.PathEscape(repo)), request, nil, http.StatusNoContent, http.StatusOK)
}

func (c *Client) GetCustomPropertyValues(ctx context.Context, owner string, repo string) ([]CustomPropertyValue, error) {
	var properties []CustomPropertyValue
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/properties/values", url.PathEscape(owner), url.PathEscape(repo)), nil, &properties, http.StatusOK); err != nil {
		return nil, err
	}

	return properties, nil
}

func (c *Client) UpdateCustomPropertyValues(ctx context.Context, owner string, repo string, properties []CustomPropertyValue) error {
	request := map[string]any{"properties": properties}
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/repos/%s/%s/properties/values", url.PathEscape(owner), url.PathEscape(repo)), request, nil, http.StatusNoContent)
}

func (c *Client) GetBranch(ctx context.Context, owner string, repo string, branch string) (*Branch, error) {
	var current Branch
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/branches/%s", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(branch)), nil, &current, http.StatusOK); err != nil {
		return nil, err
	}

	return &current, nil
}

func (c *Client) ListProtectedBranches(ctx context.Context, owner string, repo string) ([]Branch, error) {
	var branches []Branch
	path := fmt.Sprintf("/repos/%s/%s/branches?protected=true&per_page=100", url.PathEscape(owner), url.PathEscape(repo))
	if err := c.request(ctx, http.MethodGet, path, nil, &branches, http.StatusOK); err != nil {
		return nil, err
	}

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].Name < branches[j].Name
	})

	return branches, nil
}

func (c *Client) GetBranchProtection(ctx context.Context, owner string, repo string, branch string) (*BranchProtection, error) {
	var protection BranchProtection
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/repos/%s/%s/branches/%s/protection", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(branch)), nil, &protection, http.StatusOK); err != nil {
		return nil, err
	}

	return &protection, nil
}

func (c *Client) UpdateBranchProtection(ctx context.Context, owner string, repo string, branch string, request map[string]any) error {
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/repos/%s/%s/branches/%s/protection", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(branch)), request, nil, http.StatusOK)
}

func (c *Client) DeleteBranchProtection(ctx context.Context, owner string, repo string, branch string) error {
	return c.request(ctx, http.MethodDelete, fmt.Sprintf("/repos/%s/%s/branches/%s/protection", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(branch)), nil, nil, http.StatusNoContent)
}

func (c *Client) request(ctx context.Context, method string, path string, requestBody any, responseBody any, expectedStatusCodes ...int) error {
	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal github request %s %s: %w", method, path, err)
		}

		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build github request %s %s: %w", method, path, err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", defaultAPIVersion)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute github request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	for _, expectedStatus := range expectedStatusCodes {
		if resp.StatusCode == expectedStatus {
			if responseBody == nil || resp.StatusCode == http.StatusNoContent {
				io.Copy(io.Discard, resp.Body)
				return nil
			}

			if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
				return fmt.Errorf("decode github response %s %s: %w", method, path, err)
			}

			return nil
		}
	}

	var apiErrBody struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiErrBody); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode github error response %s %s: %w", method, path, err)
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Method:     method,
		Path:       path,
		Message:    apiErrBody.Message,
	}
}
