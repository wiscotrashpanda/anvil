provider "aws" {
  alias = "default"
}

provider "github" {
  alias = "default"
  owner = var.github_owner
}

provider "tfe" {
  alias        = "default"
  organization = var.tfe_organization
}
