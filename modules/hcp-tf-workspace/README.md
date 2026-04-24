# hcp-tf-workspace module

Terraform module for one repo-backed HCP Terraform workspace environment plus its shared AWS provisioner role and GitHub environment wiring.

This module assumes the GitHub repository already exists. Callers provide the GitHub repository path in `owner/repository` form, and the module creates:

- one HCP Terraform workspace
- one CloudFormation StackSet
- one AWS IAM provisioner role
- one GitHub environment with AWS variables for that environment

## Usage

```hcl
provider "github" {
  alias = "emkaytec"
  owner = "emkaytec"
}

provider "tfe" {
  alias        = "emkaytec"
  organization = "emkaytec"
}

module "workspace" {
  source = "./modules/hcp-tf-workspace"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  github_repository = "emkaytec/sample-service"
  environment       = "dev"

  aws = {
    account_id                         = "111111111111"
    stack_set_administration_role_arn  = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
  }

  workspace = {
    terraform_version = "1.10.0"
    working_directory = "terraform"
  }
}
```

## Notes

- `workspace.name` defaults to `<repository>-<environment>` and is also the default naming base for the provisioner role and StackSet name.
- `aws.region` defaults to `us-east-1`.
- When `workspace.vcs_repo` is set, the module uses `github_repository` as the HCP Terraform VCS identifier.
- The GitHub provider owner and TFE provider organization still come from the caller-level provider configuration.
