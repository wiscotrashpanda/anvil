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
  owner = var.github_owner
}

provider "tfe" {}

module "repo" {
  source = "../.."

  github_owner     = var.github_owner
  tfe_organization = var.tfe_organization

  repository = {
    name        = var.repository_name
    description = "Example Terraform-backed GitHub repository."
    topics      = ["aws", "terraform"]
  }

  environments = var.environments

  stack_set_administration_role_arn = var.stack_set_administration_role_arn
}
