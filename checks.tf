check "tfe_organization_configured_for_workspace_manifests" {
  assert {
    condition = (
      length(local.github_tf_repo_manifests) == 0 &&
      length(local.hcp_tf_workspace_manifests) == 0
      ) || (
      var.tfe_organization != null &&
      length(trimspace(var.tfe_organization)) > 0
    )
    error_message = "tfe_organization must be set when GitHubTerraformRepository or HCPTerraformWorkspace manifests are present."
  }
}

check "github_owner_matches_hcp_tf_workspace_repo_paths" {
  assert {
    condition = alltrue([
      for _, config in local.hcp_tf_workspace_module_inputs :
      length(split("/", trimspace(try(config.github_repository, "")))) == 2 &&
      split("/", trimspace(config.github_repository))[0] == var.github_owner
    ])
    error_message = "HCPTerraformWorkspace spec.githubRepository must use the same owner configured by github_owner."
  }
}

check "hcp_tf_workspace_aws_targets_configured" {
  assert {
    condition = alltrue([
      for _, config in local.hcp_tf_workspace_module_inputs :
      can(regex("^[0-9]{12}$", config.aws.account_id))
    ])
    error_message = "HCPTerraformWorkspace spec.aws.accountId must be a 12-digit AWS account ID."
  }
}
