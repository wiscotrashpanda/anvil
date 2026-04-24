# GitHub Repository Module

Terraform module for one GitHub repository. The module can manage a standalone GitHub repository, or it can also create repo-backed HCP Terraform workspaces, AWS provisioner IAM roles, CloudFormation StackSets, GitHub environments, and workspace variables when `create_terraform_workspaces` is enabled.

This module currently lives inside `emkaytec/anvil`. It is intentionally shaped like an extractable Terraform module and is expected to move to a standalone `terraform-github-repository` repository later.

## Usage

Standalone GitHub repository:

```hcl
module "repo" {
  source = "./modules/github-repo"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe
  }

  repository = {
    name        = "docs-site"
    description = "Public documentation site."
    visibility  = "public"
    topics      = ["docs", "website"]
  }
}
```

Repo-backed Terraform workspaces:

```hcl
module "repo" {
  source = "./modules/github-repo"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  create_terraform_workspaces = true

  repository = {
    name        = "sample-service"
    description = "Terraform-managed sample service."
    topics      = ["aws", "terraform"]
  }

  environments = {
    dev = {
      aws = {
        account_id = "111111111111"
      }
    }
  }

  aws = {
    stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
  }
}
```

## Examples

- `examples/basic` creates a standalone GitHub repository.
- `examples/complete` declares every available module input and creates one Terraform workspace environment.

The `.forge.example/example` directory mirrors these examples as `GitHubRepository` manifests for root-level desired-state composition.

## Requirements

| Name | Version |
| --- | --- |
| terraform | `>= 1.6.0` |
| aws | `>= 6.0, < 7.0` |
| github | `>= 6.0, < 7.0` |
| tfe | `>= 0.76.0, < 1.0` |

## Providers

| Name | Purpose |
| --- | --- |
| `aws` | Creates CloudFormation StackSets and StackSet instances when Terraform workspaces are enabled. |
| `github` | Creates the GitHub repository, default branch, GitHub environments, and environment variables. |
| `tfe` | Creates HCP Terraform workspaces, workspace settings, and workspace variables when Terraform workspaces are enabled. |

The caller owns provider configuration. The GitHub provider owner determines the repository owner, and the TFE provider organization determines the HCP Terraform organization.

## Resources

Always managed:

- `github_repository.this`
- `github_branch_default.this`, when `repository.manage_default_branch` is true

Managed only when `create_terraform_workspaces = true`:

- `tfe_workspace.this`
- `tfe_workspace_settings.this`
- `aws_cloudformation_stack_set.provisioner_roles`
- `aws_cloudformation_stack_set_instance.provisioner_roles`
- `tfe_variable.account_id`
- `tfe_variable.aws_region`
- `tfe_variable.aws_region_env`
- `tfe_variable.tfc_aws_provider_auth`
- `tfe_variable.tfc_aws_run_role_arn`
- `tfe_variable.tfc_aws_workload_identity_audience`
- `github_repository_environment.this`
- `github_actions_environment_variable.aws_region`
- `github_actions_environment_variable.aws_account_id`
- `github_actions_environment_variable.aws_provisioner_role_arn`

## Inputs

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| `repository` | object | n/a | GitHub repository settings. `name` is required. |
| `create_terraform_workspaces` | bool | `false` | Creates HCP Terraform, AWS provisioner role, and GitHub environment resources for each environment. |
| `environments` | map(object) | `{}` | Environment-specific AWS and workspace settings. Required when `create_terraform_workspaces` is true. |
| `aws` | object | `{}` | AWS defaults for generated Terraform workspace environments. |
| `workspace` | object | `{}` | HCP Terraform workspace defaults for generated environments. |

### `repository`

| Field | Default | Description |
| --- | --- | --- |
| `name` | n/a | Repository name. |
| `description` | `""` | Repository description. |
| `visibility` | `"private"` | Repository visibility: `public`, `private`, or `internal`. |
| `topics` | `[]` | Repository topics. Topics must be unique and non-blank. |
| `homepage_url` | `null` | Repository homepage URL. |
| `auto_init` | `true` | Initialize the repository with a README. |
| `archive_on_destroy` | `true` | Archive instead of deleting the repository on destroy. |
| `has_issues` | `true` | Enable GitHub Issues. |
| `has_projects` | `false` | Enable GitHub Projects. |
| `has_wiki` | `false` | Enable the repository wiki. |
| `has_discussions` | `false` | Enable GitHub Discussions. |
| `allow_merge_commit` | `false` | Allow merge commits. |
| `allow_squash_merge` | `true` | Allow squash merges. |
| `allow_rebase_merge` | `false` | Allow rebase merges. |
| `delete_branch_on_merge` | `true` | Delete branches after pull requests merge. |
| `default_branch` | `"main"` | Default branch name. |
| `manage_default_branch` | `true` | Manage the default branch with `github_branch_default`. |
| `rename_default_branch` | `false` | Rename the default branch when managing it. |

### `environments`

`environments` is keyed by environment name. Environment keys must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens.

Each environment requires `aws.account_id`. Other environment-level AWS and workspace fields override the module-level `aws` and `workspace` defaults for that environment.

### `aws`

| Field | Default | Description |
| --- | --- | --- |
| `region` | `"us-east-1"` | Default target AWS region. |
| `partition` | `"aws"` | AWS partition used when computing provisioner role ARNs. |
| `managed_policy_arns` | `["arn:aws:iam::aws:policy/ReadOnlyAccess"]` | Managed policies attached to provisioner roles. |
| `github_oidc_provider_host` | `"token.actions.githubusercontent.com"` | GitHub Actions OIDC provider host in target accounts. |
| `github_oidc_audience` | `"sts.amazonaws.com"` | GitHub Actions OIDC audience. |
| `tfe_oidc_provider_host` | `"app.terraform.io"` | HCP Terraform OIDC provider host in target accounts. |
| `tfe_oidc_audience` | `"aws.workload.identity"` | HCP Terraform OIDC audience. |
| `stack_set_name_prefix` | `null` | Optional StackSet name prefix. Defaults to the workspace name. |
| `stack_set_permission_model` | `"SELF_MANAGED"` | CloudFormation StackSet permission model: `SELF_MANAGED` or `SERVICE_MANAGED`. |
| `stack_set_administration_role_arn` | `null` | Required when using `SELF_MANAGED` StackSets with Terraform workspaces enabled. |
| `stack_set_execution_role_name` | `"AWSCloudFormationStackSetExecutionRole"` | Execution role name in target accounts for self-managed StackSets. |
| `stack_set_call_as` | `"SELF"` | StackSet caller mode: `SELF` or `DELEGATED_ADMIN`. |
| `stack_set_operation_preferences` | `null` | Optional StackSet operation preferences. |
| `retain_stack_instances_on_destroy` | `false` | Retain StackSet instances on destroy. |
| `tags` | `{}` | Tags applied to generated AWS resources. |

### `workspace`

| Field | Default | Description |
| --- | --- | --- |
| `name` | `null` | HCP Terraform workspace name. Defaults to `<repository-name>-<environment>`. |
| `project_id` | `null` | HCP Terraform project ID. |
| `project_name` | `"*"` | HCP Terraform project name used in the default OIDC subject. |
| `execution_mode` | `"remote"` | Workspace execution mode: `remote`, `local`, or `agent`. |
| `agent_pool_id` | `null` | Agent pool ID when using agent execution. |
| `terraform_version` | `null` | Workspace Terraform version. |
| `auto_apply` | `false` | Enable workspace auto-apply. |
| `queue_all_runs` | `true` | Queue all runs for the workspace. |
| `speculative_enabled` | `true` | Enable speculative plans. |
| `working_directory` | `null` | Workspace working directory. |
| `tags` | `{}` | Workspace tags. |
| `vcs_repo` | `null` | Optional HCP Terraform VCS connection. Workspaces are API/CLI-driven when omitted. |
| `manage_variables` | `true` | Manage generated HCP Terraform variables. |
| `hcp_terraform_subject` | `null` | Override the default HCP Terraform OIDC subject. |

`workspace.vcs_repo` must set exactly one of `oauth_token_id` or `github_app_installation_id`.

## Outputs

| Name | Description |
| --- | --- |
| `repository` | Created GitHub repository details. |
| `workspaces` | Generated HCP Terraform workspaces keyed by environment. |
| `provisioner_roles` | Computed IAM role names, ARNs, and OIDC subjects keyed by environment. |
| `stack_sets` | CloudFormation StackSet details keyed by environment. |

## AWS Prerequisites

The module creates provisioner roles through CloudFormation StackSets, so the StackSet prerequisites must already exist before enabling Terraform workspaces:

- the AWS provider must be authenticated in the StackSet administrator account
- for `SELF_MANAGED` StackSets, each target account must already have the configured StackSet execution role
- each target account must already have OIDC providers for `token.actions.githubusercontent.com` and `app.terraform.io`

The CloudFormation template creates only the IAM role per environment. OIDC provider bootstrap stays separate so existing accounts are not forced through a fragile create-or-conflict path.
