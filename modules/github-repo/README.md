# github-repo module

Terraform module for one GitHub repository. By default it manages only the repository. When `create_terraform_workspaces = true`, it also creates the repo-backed HCP Terraform workspaces, AWS provisioner roles, CloudFormation StackSets, GitHub environments, and workspace variables for each configured environment.

## Usage

Configure provider ownership in the caller, map providers into the module, and pass repository settings.

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
  source = "./modules/github-repo"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  repository = {
    name                   = "docs-site"
    description            = "Public documentation site."
    visibility             = "public"
    topics                 = ["docs", "website"]
    auto_init              = true
    has_issues             = true
    has_projects           = false
    has_wiki               = false
    allow_squash_merge     = true
    allow_merge_commit     = false
    allow_rebase_merge     = true
    delete_branch_on_merge = true
    default_branch         = "main"
    manage_default_branch  = true
  }
}
```

Enable Terraform workspace creation for infrastructure repositories:

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
    prod = {
      aws = {
        account_id          = "222222222222"
        region              = "us-east-2"
        managed_policy_arns = ["arn:aws:iam::aws:policy/PowerUserAccess"]
      }
    }
  }

  aws = {
    stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
  }
}
```

By default, HCP Terraform workspaces are API/CLI-driven. Set `workspace.vcs_repo` when the workspace should be connected to the GitHub repository through an existing HCP Terraform VCS connection.

## AWS Prerequisites

The module creates provisioner roles through CloudFormation StackSets, so the StackSet prerequisites must already exist before enabling Terraform workspaces:

- the AWS provider must be authenticated in the StackSet administrator account
- for `SELF_MANAGED` StackSets, each target account must already have the configured StackSet execution role
- each target account must already have OIDC providers for `token.actions.githubusercontent.com` and `app.terraform.io`

The CloudFormation template creates only the IAM role per environment. OIDC provider bootstrap stays separate so existing accounts are not forced through a fragile create-or-conflict path.

## Notes

- `visibility` defaults to `private`, `auto_init` defaults to `true`, and `default_branch` defaults to `main`.
- `manage_default_branch` and `rename_default_branch` control whether Terraform actively manages or renames the default branch.
- `create_terraform_workspaces` defaults to `false`. When it is `true`, `environments` must include at least one environment and each environment must provide `aws.account_id`.
- Each Terraform environment defaults to workspace name `<repository-name>-<environment>`, region `us-east-1`, role name `<workspace-name>-provisioner-role`, and managed policy `arn:aws:iam::aws:policy/ReadOnlyAccess`.
- HCP Terraform workspace variables are created by default. Set `workspace.manage_variables = false` to leave them unmanaged.
