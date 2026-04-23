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
    name          = var.repository_name
    description   = "Example standalone GitHub repository."
    visibility    = "private"
    autoInit      = true
    defaultBranch = "main"
    topics        = ["github", "standalone"]
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
