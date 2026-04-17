package hcpterraform

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
	"strings"
	"time"
)

const DefaultBaseURL = "https://app.terraform.io/api/v2"

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
		return fmt.Sprintf("hcp terraform api %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Message)
	}

	return fmt.Sprintf("hcp terraform api %s %s returned %d", e.Method, e.Path, e.StatusCode)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

func NewClientFromEnv() (*Client, error) {
	token := strings.TrimSpace(os.Getenv("TF_TOKEN_app_terraform_io"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("TFE_TOKEN"))
	}

	if token == "" {
		return nil, fmt.Errorf("missing HCP Terraform token: set TF_TOKEN_app_terraform_io or TFE_TOKEN")
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

type Workspace struct {
	ID                          string                      `json:"id"`
	Name                        string                      `json:"name"`
	Description                 string                      `json:"description"`
	TerraformVersion            string                      `json:"terraform-version"`
	WorkingDirectory            string                      `json:"working-directory"`
	ExecutionMode               string                      `json:"execution-mode"`
	AgentPoolID                 string                      `json:"agent-pool-id"`
	AllowDestroyPlan            bool                        `json:"allow-destroy-plan"`
	AssessmentsEnabled          bool                        `json:"assessments-enabled"`
	AutoApply                   bool                        `json:"auto-apply"`
	AutoApplyRunTrigger         bool                        `json:"auto-apply-run-trigger"`
	AutoDestroyAt               string                      `json:"auto-destroy-at"`
	AutoDestroyActivityDuration string                      `json:"auto-destroy-activity-duration"`
	FileTriggersEnabled         bool                        `json:"file-triggers-enabled"`
	GlobalRemoteState           bool                        `json:"global-remote-state"`
	ProjectRemoteState          bool                        `json:"project-remote-state"`
	QueueAllRuns                bool                        `json:"queue-all-runs"`
	SourceName                  string                      `json:"source-name"`
	SourceURL                   string                      `json:"source-url"`
	SpeculativeEnabled          bool                        `json:"speculative-enabled"`
	TriggerPatterns             []string                    `json:"trigger-patterns"`
	TriggerPrefixes             []string                    `json:"trigger-prefixes"`
	SettingOverwrites           *WorkspaceSettingOverwrites `json:"setting-overwrites"`
	VCSRepo                     *WorkspaceVCSRepo           `json:"vcs-repo"`
	ProjectID                   string
}

type WorkspaceSettingOverwrites struct {
	ExecutionMode       bool `json:"execution-mode"`
	AgentPoolID         bool `json:"agent-pool-id"`
	AutoApply           bool `json:"auto-apply"`
	FileTriggersEnabled bool `json:"file-triggers-enabled"`
	GlobalRemoteState   bool `json:"global-remote-state"`
	QueueAllRuns        bool `json:"queue-all-runs"`
	SpeculativeEnabled  bool `json:"speculative-enabled"`
	TerraformVersion    bool `json:"terraform-version"`
	TriggerPatterns     bool `json:"trigger-patterns"`
	TriggerPrefixes     bool `json:"trigger-prefixes"`
	WorkingDirectory    bool `json:"working-directory"`
}

type WorkspaceVCSRepo struct {
	Identifier        string `json:"identifier"`
	Branch            string `json:"branch"`
	IngressSubmodules bool   `json:"ingress-submodules"`
	OAuthTokenID      string `json:"oauth-token-id"`
	TagsRegex         string `json:"tags-regex"`
}

type WorkspaceTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkspaceTagBinding struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WorkspaceVariable struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Category    string `json:"category"`
	HCL         bool   `json:"hcl"`
	Sensitive   bool   `json:"sensitive"`
}

type VariableSet struct {
	ID string `json:"id"`
}

type WorkspaceRequest struct {
	ProjectID                   *string
	Description                 *string
	TerraformVersion            *string
	WorkingDirectory            *string
	ExecutionMode               *string
	AgentPoolID                 *string
	AllowDestroyPlan            *bool
	AssessmentsEnabled          *bool
	AutoApply                   *bool
	AutoApplyRunTrigger         *bool
	AutoDestroyAt               *string
	AutoDestroyActivityDuration *string
	FileTriggersEnabled         *bool
	GlobalRemoteState           *bool
	ProjectRemoteState          *bool
	QueueAllRuns                *bool
	SourceName                  *string
	SourceURL                   *string
	SpeculativeEnabled          *bool
	TriggerPatterns             []string
	TriggerPrefixes             []string
	SettingOverwrites           *WorkspaceSettingOverwrites
	VCSRepo                     *WorkspaceVCSRepo
	TagBindings                 []WorkspaceTagBinding
}

type workspaceResponseData struct {
	ID            string                      `json:"id"`
	Attributes    workspaceAttributesResponse `json:"attributes"`
	Relationships struct {
		Project relationshipData `json:"project"`
	} `json:"relationships"`
}

func (c *Client) GetWorkspace(ctx context.Context, organization string, name string) (*Workspace, error) {
	path := fmt.Sprintf("/organizations/%s/workspaces/%s", url.PathEscape(organization), url.PathEscape(name))
	var response struct {
		Data workspaceResponseData `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return response.Data.toWorkspace(), nil
}

func (c *Client) CreateWorkspace(ctx context.Context, organization string, name string, request WorkspaceRequest) (*Workspace, error) {
	path := fmt.Sprintf("/organizations/%s/workspaces", url.PathEscape(organization))
	var response struct {
		Data workspaceResponseData `json:"data"`
	}

	if err := c.request(ctx, http.MethodPost, path, buildWorkspacePayload(name, request), &response, http.StatusCreated, http.StatusOK); err != nil {
		return nil, err
	}

	return response.Data.toWorkspace(), nil
}

func (c *Client) UpdateWorkspace(ctx context.Context, workspaceID string, request WorkspaceRequest) (*Workspace, error) {
	path := fmt.Sprintf("/workspaces/%s", url.PathEscape(workspaceID))
	var response struct {
		Data workspaceResponseData `json:"data"`
	}

	if err := c.request(ctx, http.MethodPatch, path, buildWorkspacePayload("", request), &response, http.StatusOK); err != nil {
		return nil, err
	}

	return response.Data.toWorkspace(), nil
}

func (c *Client) ListWorkspaceTags(ctx context.Context, workspaceID string) ([]WorkspaceTag, error) {
	path := fmt.Sprintf("/workspaces/%s/relationships/tags", url.PathEscape(workspaceID))
	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	tags := make([]WorkspaceTag, 0, len(response.Data))
	for _, item := range response.Data {
		tags = append(tags, WorkspaceTag{ID: item.ID, Name: item.Attributes.Name})
	}

	return tags, nil
}

func (c *Client) AddWorkspaceTags(ctx context.Context, workspaceID string, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}

	data := make([]map[string]any, 0, len(tagNames))
	for _, tagName := range tagNames {
		data = append(data, map[string]any{
			"type": "tags",
			"attributes": map[string]any{
				"name": tagName,
			},
		})
	}

	return c.request(ctx, http.MethodPost, fmt.Sprintf("/workspaces/%s/relationships/tags", url.PathEscape(workspaceID)), map[string]any{"data": data}, nil, http.StatusNoContent)
}

func (c *Client) RemoveWorkspaceTags(ctx context.Context, workspaceID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	data := make([]map[string]any, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		data = append(data, map[string]any{
			"type": "tags",
			"id":   tagID,
		})
	}

	return c.request(ctx, http.MethodDelete, fmt.Sprintf("/workspaces/%s/relationships/tags", url.PathEscape(workspaceID)), map[string]any{"data": data}, nil, http.StatusNoContent)
}

func (c *Client) ListWorkspaceTagBindings(ctx context.Context, workspaceID string) ([]WorkspaceTagBinding, error) {
	path := fmt.Sprintf("/workspaces/%s/tag-bindings", url.PathEscape(workspaceID))
	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	bindings := make([]WorkspaceTagBinding, 0, len(response.Data))
	for _, item := range response.Data {
		bindings = append(bindings, WorkspaceTagBinding{
			ID:    item.ID,
			Key:   item.Attributes.Key,
			Value: item.Attributes.Value,
		})
	}

	return bindings, nil
}

func (c *Client) ReplaceWorkspaceTagBindings(ctx context.Context, workspaceID string, bindings []WorkspaceTagBinding) error {
	data := make([]map[string]any, 0, len(bindings))
	for _, binding := range bindings {
		data = append(data, map[string]any{
			"type": "tag-bindings",
			"attributes": map[string]any{
				"key":   binding.Key,
				"value": binding.Value,
			},
		})
	}

	return c.request(ctx, http.MethodPatch, fmt.Sprintf("/workspaces/%s/tag-bindings", url.PathEscape(workspaceID)), map[string]any{"data": data}, nil, http.StatusOK)
}

func (c *Client) ListRemoteStateConsumers(ctx context.Context, workspaceID string) ([]string, error) {
	path := fmt.Sprintf("/workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID))
	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(response.Data))
	for _, item := range response.Data {
		ids = append(ids, item.ID)
	}

	return ids, nil
}

func (c *Client) AddRemoteStateConsumers(ctx context.Context, workspaceID string, consumerIDs []string) error {
	if len(consumerIDs) == 0 {
		return nil
	}

	return c.request(ctx, http.MethodPost, fmt.Sprintf("/workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID)), workspaceIDListPayload(consumerIDs), nil, http.StatusNoContent)
}

func (c *Client) RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, consumerIDs []string) error {
	if len(consumerIDs) == 0 {
		return nil
	}

	return c.request(ctx, http.MethodDelete, fmt.Sprintf("/workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID)), workspaceIDListPayload(consumerIDs), nil, http.StatusNoContent)
}

func (c *Client) ListVariables(ctx context.Context, workspaceID string) ([]WorkspaceVariable, error) {
	path := fmt.Sprintf("/workspaces/%s/vars", url.PathEscape(workspaceID))
	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Key         string `json:"key"`
				Value       string `json:"value"`
				Description string `json:"description"`
				Category    string `json:"category"`
				HCL         bool   `json:"hcl"`
				Sensitive   bool   `json:"sensitive"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	variables := make([]WorkspaceVariable, 0, len(response.Data))
	for _, item := range response.Data {
		variables = append(variables, WorkspaceVariable{
			ID:          item.ID,
			Key:         item.Attributes.Key,
			Value:       item.Attributes.Value,
			Description: item.Attributes.Description,
			Category:    item.Attributes.Category,
			HCL:         item.Attributes.HCL,
			Sensitive:   item.Attributes.Sensitive,
		})
	}

	return variables, nil
}

func (c *Client) CreateVariable(ctx context.Context, workspaceID string, variable WorkspaceVariable) error {
	path := fmt.Sprintf("/workspaces/%s/vars", url.PathEscape(workspaceID))
	request := map[string]any{
		"data": map[string]any{
			"type": "vars",
			"attributes": map[string]any{
				"key":         variable.Key,
				"value":       variable.Value,
				"description": variable.Description,
				"category":    variable.Category,
				"hcl":         variable.HCL,
				"sensitive":   variable.Sensitive,
			},
		},
	}

	return c.request(ctx, http.MethodPost, path, request, nil, http.StatusCreated, http.StatusOK)
}

func (c *Client) UpdateVariable(ctx context.Context, workspaceID string, variableID string, variable WorkspaceVariable) error {
	path := fmt.Sprintf("/workspaces/%s/vars/%s", url.PathEscape(workspaceID), url.PathEscape(variableID))
	request := map[string]any{
		"data": map[string]any{
			"id":   variableID,
			"type": "vars",
			"attributes": map[string]any{
				"key":         variable.Key,
				"value":       variable.Value,
				"description": variable.Description,
				"category":    variable.Category,
				"hcl":         variable.HCL,
				"sensitive":   variable.Sensitive,
			},
		},
	}

	return c.request(ctx, http.MethodPatch, path, request, nil, http.StatusOK)
}

func (c *Client) DeleteVariable(ctx context.Context, workspaceID string, variableID string) error {
	path := fmt.Sprintf("/workspaces/%s/vars/%s", url.PathEscape(workspaceID), url.PathEscape(variableID))
	return c.request(ctx, http.MethodDelete, path, nil, nil, http.StatusNoContent)
}

func (c *Client) ListWorkspaceVariableSets(ctx context.Context, workspaceID string) ([]VariableSet, error) {
	path := fmt.Sprintf("/workspaces/%s/varsets", url.PathEscape(workspaceID))
	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := c.request(ctx, http.MethodGet, path, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	varsets := make([]VariableSet, 0, len(response.Data))
	for _, item := range response.Data {
		varsets = append(varsets, VariableSet{ID: item.ID})
	}

	return varsets, nil
}

func (c *Client) AddVariableSetToWorkspace(ctx context.Context, variableSetID string, workspaceID string) error {
	path := fmt.Sprintf("/varsets/%s/relationships/workspaces", url.PathEscape(variableSetID))
	return c.request(ctx, http.MethodPost, path, workspaceIDListPayload([]string{workspaceID}), nil, http.StatusNoContent)
}

func (c *Client) RemoveVariableSetFromWorkspace(ctx context.Context, variableSetID string, workspaceID string) error {
	path := fmt.Sprintf("/varsets/%s/relationships/workspaces", url.PathEscape(variableSetID))
	return c.request(ctx, http.MethodDelete, path, workspaceIDListPayload([]string{workspaceID}), nil, http.StatusNoContent)
}

type workspaceAttributesResponse struct {
	Name                        string                      `json:"name"`
	Description                 string                      `json:"description"`
	TerraformVersion            string                      `json:"terraform-version"`
	WorkingDirectory            string                      `json:"working-directory"`
	ExecutionMode               string                      `json:"execution-mode"`
	AgentPoolID                 string                      `json:"agent-pool-id"`
	AllowDestroyPlan            bool                        `json:"allow-destroy-plan"`
	AssessmentsEnabled          bool                        `json:"assessments-enabled"`
	AutoApply                   bool                        `json:"auto-apply"`
	AutoApplyRunTrigger         bool                        `json:"auto-apply-run-trigger"`
	AutoDestroyAt               string                      `json:"auto-destroy-at"`
	AutoDestroyActivityDuration string                      `json:"auto-destroy-activity-duration"`
	FileTriggersEnabled         bool                        `json:"file-triggers-enabled"`
	GlobalRemoteState           bool                        `json:"global-remote-state"`
	ProjectRemoteState          bool                        `json:"project-remote-state"`
	QueueAllRuns                bool                        `json:"queue-all-runs"`
	SourceName                  string                      `json:"source-name"`
	SourceURL                   string                      `json:"source-url"`
	SpeculativeEnabled          bool                        `json:"speculative-enabled"`
	TriggerPatterns             []string                    `json:"trigger-patterns"`
	TriggerPrefixes             []string                    `json:"trigger-prefixes"`
	SettingOverwrites           *WorkspaceSettingOverwrites `json:"setting-overwrites"`
	VCSRepo                     *WorkspaceVCSRepo           `json:"vcs-repo"`
}

type relationshipData struct {
	Data *struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (d workspaceResponseData) toWorkspace() *Workspace {
	workspace := &Workspace{
		ID:                          d.ID,
		Name:                        d.Attributes.Name,
		Description:                 d.Attributes.Description,
		TerraformVersion:            d.Attributes.TerraformVersion,
		WorkingDirectory:            d.Attributes.WorkingDirectory,
		ExecutionMode:               d.Attributes.ExecutionMode,
		AgentPoolID:                 d.Attributes.AgentPoolID,
		AllowDestroyPlan:            d.Attributes.AllowDestroyPlan,
		AssessmentsEnabled:          d.Attributes.AssessmentsEnabled,
		AutoApply:                   d.Attributes.AutoApply,
		AutoApplyRunTrigger:         d.Attributes.AutoApplyRunTrigger,
		AutoDestroyAt:               d.Attributes.AutoDestroyAt,
		AutoDestroyActivityDuration: d.Attributes.AutoDestroyActivityDuration,
		FileTriggersEnabled:         d.Attributes.FileTriggersEnabled,
		GlobalRemoteState:           d.Attributes.GlobalRemoteState,
		ProjectRemoteState:          d.Attributes.ProjectRemoteState,
		QueueAllRuns:                d.Attributes.QueueAllRuns,
		SourceName:                  d.Attributes.SourceName,
		SourceURL:                   d.Attributes.SourceURL,
		SpeculativeEnabled:          d.Attributes.SpeculativeEnabled,
		TriggerPatterns:             d.Attributes.TriggerPatterns,
		TriggerPrefixes:             d.Attributes.TriggerPrefixes,
		SettingOverwrites:           d.Attributes.SettingOverwrites,
		VCSRepo:                     d.Attributes.VCSRepo,
	}

	if d.Relationships.Project.Data != nil {
		workspace.ProjectID = d.Relationships.Project.Data.ID
	}

	return workspace
}

func buildWorkspacePayload(name string, request WorkspaceRequest) map[string]any {
	attributes := make(map[string]any)
	if name != "" {
		attributes["name"] = name
	}
	if request.Description != nil {
		attributes["description"] = *request.Description
	}
	if request.TerraformVersion != nil {
		attributes["terraform-version"] = *request.TerraformVersion
	}
	if request.WorkingDirectory != nil {
		attributes["working-directory"] = *request.WorkingDirectory
	}
	if request.ExecutionMode != nil {
		attributes["execution-mode"] = *request.ExecutionMode
	}
	if request.AgentPoolID != nil {
		attributes["agent-pool-id"] = *request.AgentPoolID
	}
	if request.AllowDestroyPlan != nil {
		attributes["allow-destroy-plan"] = *request.AllowDestroyPlan
	}
	if request.AssessmentsEnabled != nil {
		attributes["assessments-enabled"] = *request.AssessmentsEnabled
	}
	if request.AutoApply != nil {
		attributes["auto-apply"] = *request.AutoApply
	}
	if request.AutoApplyRunTrigger != nil {
		attributes["auto-apply-run-trigger"] = *request.AutoApplyRunTrigger
	}
	if request.AutoDestroyAt != nil {
		attributes["auto-destroy-at"] = *request.AutoDestroyAt
	}
	if request.AutoDestroyActivityDuration != nil {
		attributes["auto-destroy-activity-duration"] = *request.AutoDestroyActivityDuration
	}
	if request.FileTriggersEnabled != nil {
		attributes["file-triggers-enabled"] = *request.FileTriggersEnabled
	}
	if request.GlobalRemoteState != nil {
		attributes["global-remote-state"] = *request.GlobalRemoteState
	}
	if request.ProjectRemoteState != nil {
		attributes["project-remote-state"] = *request.ProjectRemoteState
	}
	if request.QueueAllRuns != nil {
		attributes["queue-all-runs"] = *request.QueueAllRuns
	}
	if request.SourceName != nil {
		attributes["source-name"] = *request.SourceName
	}
	if request.SourceURL != nil {
		attributes["source-url"] = *request.SourceURL
	}
	if request.SpeculativeEnabled != nil {
		attributes["speculative-enabled"] = *request.SpeculativeEnabled
	}
	if request.TriggerPatterns != nil {
		attributes["trigger-patterns"] = request.TriggerPatterns
	}
	if request.TriggerPrefixes != nil {
		attributes["trigger-prefixes"] = request.TriggerPrefixes
	}
	if request.SettingOverwrites != nil {
		attributes["setting-overwrites"] = request.SettingOverwrites
	}
	if request.VCSRepo != nil {
		attributes["vcs-repo"] = request.VCSRepo
	}

	data := map[string]any{
		"type": "workspaces",
	}
	if len(attributes) > 0 {
		data["attributes"] = attributes
	}

	relationships := make(map[string]any)
	if request.ProjectID != nil {
		relationships["project"] = map[string]any{
			"data": map[string]any{
				"type": "projects",
				"id":   *request.ProjectID,
			},
		}
	}
	if request.TagBindings != nil {
		items := make([]map[string]any, 0, len(request.TagBindings))
		for _, binding := range request.TagBindings {
			items = append(items, map[string]any{
				"type": "tag-bindings",
				"attributes": map[string]any{
					"key":   binding.Key,
					"value": binding.Value,
				},
			})
		}
		relationships["tag-bindings"] = map[string]any{"data": items}
	}
	if len(relationships) > 0 {
		data["relationships"] = relationships
	}

	return map[string]any{"data": data}
}

func workspaceIDListPayload(workspaceIDs []string) map[string]any {
	data := make([]map[string]any, 0, len(workspaceIDs))
	for _, workspaceID := range workspaceIDs {
		data = append(data, map[string]any{
			"type": "workspaces",
			"id":   workspaceID,
		})
	}

	return map[string]any{"data": data}
}

func (c *Client) request(ctx context.Context, method string, path string, requestBody any, responseBody any, expectedStatusCodes ...int) error {
	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal hcp terraform request %s %s: %w", method, path, err)
		}

		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build hcp terraform request %s %s: %w", method, path, err)
	}

	req.Header.Set("Accept", "application/vnd.api+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute hcp terraform request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	for _, expectedStatus := range expectedStatusCodes {
		if resp.StatusCode == expectedStatus {
			if responseBody == nil || resp.StatusCode == http.StatusNoContent {
				io.Copy(io.Discard, resp.Body)
				return nil
			}

			if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
				return fmt.Errorf("decode hcp terraform response %s %s: %w", method, path, err)
			}

			return nil
		}
	}

	var apiErrBody struct {
		Errors []struct {
			Detail string `json:"detail"`
			Title  string `json:"title"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiErrBody); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode hcp terraform error response %s %s: %w", method, path, err)
	}

	message := ""
	if len(apiErrBody.Errors) > 0 {
		message = apiErrBody.Errors[0].Detail
		if message == "" {
			message = apiErrBody.Errors[0].Title
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Method:     method,
		Path:       path,
		Message:    message,
	}
}
