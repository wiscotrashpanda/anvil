check "tfe_organization_configured_for_github_tf_repo_manifests" {
  assert {
    condition = length(local.github_tf_repo_manifests) == 0 || (
      var.tfe_organization != null &&
      length(trimspace(var.tfe_organization)) > 0
    )
    error_message = "tfe_organization must be set when GitHubTerraformRepository manifests are present."
  }
}
