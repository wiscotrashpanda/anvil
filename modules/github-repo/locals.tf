locals {
  terraform_environments = var.create_terraform_workspaces ? var.environments : {}
  repository_default_branch = try(
    github_branch_default.this[0].branch,
    var.repository.default_branch,
  )

  environment_workspace_inputs = {
    for environment, config in local.terraform_environments :
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
          local.repository_default_branch,
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

  environment_regions = {
    for environment, config in local.terraform_environments :
    environment => coalesce(try(config.aws.region, null), var.aws.region)
  }

  environment_partitions = {
    for environment, config in local.terraform_environments :
    environment => coalesce(try(config.aws.partition, null), var.aws.partition)
  }

  environment_managed_policy_arns = {
    for environment, config in local.terraform_environments :
    environment => coalesce(try(config.aws.managed_policy_arns, null), var.aws.managed_policy_arns)
  }

  workspace_names = {
    for environment, _ in local.terraform_environments :
    environment => coalesce(local.environment_workspace_inputs[environment].name, "${github_repository.this.name}-${environment}")
  }

  provisioner_role_names = {
    for environment, _ in local.terraform_environments :
    environment => "${local.workspace_names[environment]}-provisioner-role"
  }

  provisioner_role_arns = {
    for environment, config in local.terraform_environments :
    environment => "arn:${local.environment_partitions[environment]}:iam::${config.aws.account_id}:role/${local.provisioner_role_names[environment]}"
  }

  github_actions_subjects = {
    for environment, config in local.terraform_environments :
    environment => coalesce(
      try(config.aws.github_actions_subject, null),
      "repo:${github_repository.this.full_name}:*",
    )
  }

  tfe_subjects = {
    for environment, _ in local.terraform_environments :
    environment => coalesce(
      local.environment_workspace_inputs[environment].hcp_terraform_subject,
      "organization:${tfe_workspace.this[environment].organization}:project:${local.environment_workspace_inputs[environment].project_name}:workspace:${local.workspace_names[environment]}:run_phase:*",
    )
  }

  stack_set_name_bases = {
    for environment, _ in local.terraform_environments :
    environment => var.aws.stack_set_name_prefix == null ? local.workspace_names[environment] : "${var.aws.stack_set_name_prefix}-${environment}"
  }

  stack_set_names = {
    for environment, _ in local.terraform_environments :
    environment => "${local.stack_set_name_bases[environment]}-provisioner-roles"
  }

  common_tags = {
    for environment, _ in local.terraform_environments :
    environment => merge(var.aws.tags, {
      ManagedBy  = "Terraform"
      Module     = "terraform-github-repository"
      Repository = github_repository.this.name
    })
  }

  use_tfe_vcs_repo = {
    for environment, workspace in local.environment_workspace_inputs :
    environment => workspace.vcs_repo != null
  }

  managed_variable_environments = {
    for environment, workspace in local.environment_workspace_inputs :
    environment => workspace
    if coalesce(workspace.manage_variables, true)
  }
}
