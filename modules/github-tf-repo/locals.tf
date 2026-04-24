locals {
  environment_workspace_inputs = {
    for environment, config in var.environments :
    environment => {
      name                = try(coalesce(try(config.workspace.name, null), var.workspace.name), null)
      project_id          = try(coalesce(try(config.workspace.project_id, null), var.workspace.project_id), null)
      project_name        = try(coalesce(try(config.workspace.project_name, null), var.workspace.project_name), null)
      execution_mode      = try(coalesce(try(config.workspace.execution_mode, null), var.workspace.execution_mode), null)
      agent_pool_id       = try(coalesce(try(config.workspace.agent_pool_id, null), var.workspace.agent_pool_id), null)
      terraform_version   = try(coalesce(try(config.workspace.terraform_version, null), var.workspace.terraform_version), null)
      auto_apply          = try(coalesce(try(config.workspace.auto_apply, null), var.workspace.auto_apply), null)
      queue_all_runs      = try(coalesce(try(config.workspace.queue_all_runs, null), var.workspace.queue_all_runs), null)
      speculative_enabled = try(coalesce(try(config.workspace.speculative_enabled, null), var.workspace.speculative_enabled), null)
      working_directory   = try(coalesce(try(config.workspace.working_directory, null), var.workspace.working_directory), null)
      tags                = try(coalesce(try(config.workspace.tags, null), var.workspace.tags), null)
      manage_variables    = try(coalesce(try(config.workspace.manage_variables, null), var.workspace.manage_variables), null)

      hcp_terraform_subject = try(coalesce(
        try(config.workspace.hcp_terraform_subject, null),
        var.workspace.hcp_terraform_subject,
      ), null)

      vcs_repo = try(config.workspace.vcs_repo, null) == null && var.workspace.vcs_repo == null ? null : {
        branch = try(coalesce(
          try(config.workspace.vcs_repo.branch, null),
          try(var.workspace.vcs_repo.branch, null),
          module.github_repo.repository.default_branch,
        ), null)
        oauth_token_id = try(coalesce(
          try(config.workspace.vcs_repo.oauth_token_id, null),
          try(var.workspace.vcs_repo.oauth_token_id, null),
        ), null)
        github_app_installation_id = try(coalesce(
          try(config.workspace.vcs_repo.github_app_installation_id, null),
          try(var.workspace.vcs_repo.github_app_installation_id, null),
        ), null)
        ingress_submodules = try(coalesce(
          try(config.workspace.vcs_repo.ingress_submodules, null),
          try(var.workspace.vcs_repo.ingress_submodules, null),
        ), null)
        trigger_patterns = try(coalesce(
          try(config.workspace.vcs_repo.trigger_patterns, null),
          try(var.workspace.vcs_repo.trigger_patterns, null),
        ), null)
        trigger_prefixes = try(coalesce(
          try(config.workspace.vcs_repo.trigger_prefixes, null),
          try(var.workspace.vcs_repo.trigger_prefixes, null),
        ), null)
      }
    }
  }
}
