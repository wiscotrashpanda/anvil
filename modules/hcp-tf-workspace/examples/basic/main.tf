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
  region = var.stack_set_region
}

provider "github" {
  alias = "emkaytec"
  owner = var.github_owner
}

provider "tfe" {
  alias        = "emkaytec"
  organization = var.tfe_organization
}

module "workspace" {
  source = "../.."

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  github_repository = var.github_repository
  environment       = var.environment

  aws = {
    account_id                        = var.account_id
    stack_set_administration_role_arn = var.stack_set_administration_role_arn
  }

  workspace = {
    terraform_version = "1.10.0"
  }
}
