# github-repo module

Terraform module for one GitHub repository. This is the shared repository-creation module used directly for standalone repositories and indirectly by `github-tf-repo`.

The module manages only the GitHub repository itself: no HCP Terraform workspaces, no CloudFormation StackSets, and no AWS IAM provisioner roles.

## Usage

Configure the GitHub provider owner in the caller, map that provider into the module, and pass repository settings. The input object intentionally matches the `repository` variable used by `modules/github-tf-repo` so both paths create the repo from the same settings.

```hcl
provider "github" {
  alias = "emkaytec"
  owner = "emkaytec"
}

module "repo" {
  source = "./modules/github-repo"

  providers = {
    github = github.emkaytec
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

## Notes

- `visibility` defaults to `private`, `auto_init` defaults to `true`, and `default_branch` defaults to `main`.
- `manage_default_branch` and `rename_default_branch` let callers control whether Terraform actively manages or renames the default branch.
- At the root composition layer, `GitHubRepository` manifests are translated into this shared repository input shape before reaching the module.
