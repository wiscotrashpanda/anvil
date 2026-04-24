# github-tf-repo module

Terraform module for one repository-shaped unit of infrastructure:

- one GitHub repository
- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- one AWS IAM provisioner role per environment trusted by GitHub Actions and HCP Terraform

This is the Terraform-module version of the Forge `manifest compose terraform-github-repo` idea: a single declaration fans out into the primitive resources needed for a repo-backed Terraform workload.

The GitHub repository itself is created through the shared sibling `modules/github-repo` module, and each environment's HCP Terraform workspace plus AWS role wiring is created through the shared sibling `modules/hcp-tf-workspace` module. That keeps both parts reusable on their own while preserving the one-declaration repo-backed workflow here.

The module currently lives inside `emkaytec/anvil` while the baseline architecture direction settles. It is intentionally isolated under `modules/github-tf-repo` so it can later move to a standalone `terraform-aws-github-tf-repo` repository with minimal path churn.

## Usage

Configure provider ownership in the root module, explicitly map those providers into the module, then pass only repository-specific workload inputs.

```hcl
provider "github" {
  alias = "emkaytec"
  owner = "emkaytec"
}

provider "tfe" {
  alias        = "emkaytec"
  organization = "emkaytec"
}

module "repo" {
  source = "./modules/github-tf-repo"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  repository = {
    name        = "sample-service"
    description = "Terraform-managed sample service."
    topics      = ["aws", "terraform"]
  }

  environments = {
    dev = {
      account_id = "111111111111"
    }
    prod = {
      account_id          = "222222222222"
      region              = "us-east-2"
      managed_policy_arns = ["arn:aws:iam::aws:policy/PowerUserAccess"]
    }
  }

  stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
}
```

By default, the module keeps HCP Terraform workspaces API/CLI-driven. Set `tfe_vcs_repo` when the workspace should be connected to the GitHub repository through an existing HCP Terraform VCS connection.

```hcl
tfe_vcs_repo = {
  oauth_token_id = "ot-..."
  branch         = "main"
}
```

## AWS Prerequisites

The module creates the shared provisioner role through CloudFormation StackSets, so the StackSet prerequisites must already exist:

- the AWS provider must be authenticated in the StackSet administrator account
- for `SELF_MANAGED` StackSets, each target account must already have the configured StackSet execution role
- each target account must already have OIDC providers for `token.actions.githubusercontent.com` and `app.terraform.io`

The CloudFormation template creates only the IAM role per environment. OIDC provider bootstrap stays separate so existing accounts are not forced through a fragile create-or-conflict path.

## Defaults

The repo defaults to private, `auto_init = true`, squash merges enabled, merge/rebase commits disabled, branch deletion on merge enabled, and `main` as the default branch.

Each environment defaults to:

- workspace name: `<repository-name>-<environment>`
- region: `us-east-1`
- role name: `<workspace-name>-provisioner-role`
- managed policy: `arn:aws:iam::aws:policy/ReadOnlyAccess`

HCP Terraform workspace variables are created by default. `TFC_AWS_RUN_ROLE_ARN` points at the shared provisioner role.

- `account_id`
- `aws_region`
- `AWS_REGION`
- `TFC_AWS_PROVIDER_AUTH`
- `TFC_AWS_RUN_ROLE_ARN`
- `TFC_AWS_WORKLOAD_IDENTITY_AUDIENCE`

Set `manage_tfe_workspace_variables = false` to leave workspace variables unmanaged.

## OIDC Subject Defaults

GitHub Actions defaults to the created repository's full name:

```text
repo:<repository.full_name>:*
```

HCP Terraform defaults to the organization resolved by the TFE provider/workspace:

```text
organization:<workspace.organization>:project:<tfe_project_name>:workspace:<workspace-name>:run_phase:*
```

Use `github_actions_subject` or `tfe_subject` on an environment to narrow the corresponding trust statement.
