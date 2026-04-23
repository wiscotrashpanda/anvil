# github-repo module

Terraform module for one standalone GitHub repository.

The module manages only the GitHub repository itself: no HCP Terraform workspaces, no CloudFormation StackSets, and no AWS IAM provisioner roles.

## Usage

Configure the GitHub provider owner in the caller, map that provider into the module, and pass repository settings. The GitHub owner comes from the provider configuration rather than from the module input.

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
    name          = "docs-site"
    visibility    = "public"
    description   = "Public documentation site."
    autoInit      = true
    defaultBranch = "main"
    topics        = ["docs", "website"]
    features = {
      hasIssues   = true
      hasProjects = false
      hasWiki     = false
    }
    mergePolicy = {
      allowSquashMerge    = true
      allowMergeCommit    = false
      allowRebaseMerge    = true
      deleteBranchOnMerge = true
    }
  }
}
```

## Notes

- `visibility` defaults to `private` when omitted so newly created repositories stay on the safe side.
- `defaultBranch` is managed only when set. Renaming the default branch still requires the branch to exist.
- `name` is the only required repository setting. At the root composition layer, `metadata.name` can supply that value when the manifest omits `spec.repository.name`.
- Omitted `features`, `mergePolicy`, `homepage`, and `topics` fields are passed through as unset values, so GitHub/provider defaults remain in effect unless you declare them.
