module "github_repo" {
  source = "../github-repo"

  providers = {
    github = github
  }

  repository = var.repository
}

module "hcp_tf_workspace" {
  for_each = var.environments

  source = "../hcp-tf-workspace"

  providers = {
    aws    = aws
    github = github
    tfe    = tfe
  }

  github_repository      = module.github_repo.repository.full_name
  environment            = each.key
  account_id             = each.value.account_id
  region                 = coalesce(try(each.value.region, null), var.default_region)
  workspace_name         = try(each.value.workspace_name, null)
  managed_policy_arns    = coalesce(try(each.value.managed_policy_arns, null), var.managed_policy_arns)
  github_actions_subject = try(each.value.github_actions_subject, null)
  tfe_subject            = try(each.value.tfe_subject, null)

  aws_partition             = var.aws_partition
  github_oidc_provider_host = var.github_oidc_provider_host
  github_oidc_audience      = var.github_oidc_audience
  tfe_oidc_provider_host    = var.tfe_oidc_provider_host
  tfe_oidc_audience         = var.tfe_oidc_audience

  tfe_project_id                    = var.tfe_project_id
  tfe_project_name                  = var.tfe_project_name
  tfe_workspace_execution_mode      = var.tfe_workspace_execution_mode
  tfe_workspace_agent_pool_id       = var.tfe_workspace_agent_pool_id
  tfe_workspace_terraform_version   = var.tfe_workspace_terraform_version
  tfe_workspace_auto_apply          = var.tfe_workspace_auto_apply
  tfe_workspace_queue_all_runs      = var.tfe_workspace_queue_all_runs
  tfe_workspace_speculative_enabled = var.tfe_workspace_speculative_enabled
  tfe_workspace_working_directory   = var.tfe_workspace_working_directory
  tfe_workspace_tags                = var.tfe_workspace_tags
  tfe_vcs_repo = var.tfe_vcs_repo == null ? null : merge(
    var.tfe_vcs_repo,
    { branch = coalesce(try(var.tfe_vcs_repo.branch, null), module.github_repo.repository.default_branch) },
  )
  manage_tfe_workspace_variables = var.manage_tfe_workspace_variables

  stack_set_name_prefix             = var.stack_set_name_prefix
  stack_set_permission_model        = var.stack_set_permission_model
  stack_set_administration_role_arn = var.stack_set_administration_role_arn
  stack_set_execution_role_name     = var.stack_set_execution_role_name
  stack_set_call_as                 = var.stack_set_call_as
  stack_set_operation_preferences   = var.stack_set_operation_preferences
  retain_stack_instances_on_destroy = var.retain_stack_instances_on_destroy
  tags                              = var.tags
}
