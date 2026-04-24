check "tfe_organization_configured_for_workspace_repositories" {
  assert {
    condition = (
      length(local.terraform_workspace_repo_inputs) == 0
      ) || (
      var.tfe_organization != null &&
      length(trimspace(var.tfe_organization)) > 0
    )
    error_message = "tfe_organization must be set when GitHubRepository manifests create Terraform workspaces."
  }
}

check "terraform_workspace_repositories_have_environments" {
  assert {
    condition = alltrue([
      for _, config in local.terraform_workspace_repo_inputs :
      length(config.environments) > 0
    ])
    error_message = "GitHubRepository manifests with createTerraformWorkspaces enabled must include at least one environment."
  }
}

check "terraform_workspace_aws_targets_configured" {
  assert {
    condition = alltrue(flatten([
      for _, config in local.terraform_workspace_repo_inputs : [
        for _, environment in config.environments :
        can(regex("^[0-9]{12}$", environment.aws.account_id))
      ]
    ]))
    error_message = "Terraform workspace environments must set aws.accountId to a 12-digit AWS account ID."
  }
}
