provider "aws" {
}

provider "github" {
  owner = var.github_owner
}

provider "tfe" {
  organization = var.tfe_organization
}
