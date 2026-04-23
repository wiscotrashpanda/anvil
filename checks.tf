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
    error_message = "HCPTerraformWorkspace spec.github_repository must use the same owner configured by github_owner."
  }
}
