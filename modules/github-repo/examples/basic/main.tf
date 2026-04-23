terraform {
  required_version = ">= 1.6.0"

  required_providers {
    github = {
      source  = "integrations/github"
      version = ">= 6.0, < 7.0"
    }
  }
}

provider "github" {
  alias = "emkaytec"
  owner = var.github_owner
}

module "repo" {
  source = "../.."

  providers = {
    github = github.emkaytec
  }

  repository = {
    name                   = var.repository_name
    description            = "Example standalone GitHub repository."
    visibility             = "private"
    topics                 = ["github", "standalone"]
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
