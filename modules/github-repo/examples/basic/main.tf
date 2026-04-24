terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.0, < 7.0"
    }

    github = {
      source  = "integrations/github"
      version = ">= 6.0, < 7.0"
    }

    tfe = {
      source  = "hashicorp/tfe"
      version = ">= 0.76.0, < 1.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

provider "github" {
  alias = "emkaytec"
  owner = var.github_owner
}

provider "tfe" {}

module "repo" {
  source = "../.."

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe
  }

  create_terraform_workspaces = false

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
