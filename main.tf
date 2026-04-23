module "github_repo" {
  for_each = local.github_repo_manifests

  source = "./modules/github-repo"

  providers = {
    github = github.emkaytec
  }

  repository = merge(
    { name = each.key },
    try(each.value.spec.repository, each.value.spec),
  )
}

module "github_tf_repo" {
  for_each = local.github_tf_repo_manifests

  source = "./modules/github-tf-repo"

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  repository   = merge({ name = each.key }, try(each.value.spec.repository, {}))
  environments = each.value.spec.environments

  default_region                  = try(each.value.spec.default_region, "us-east-1")
  aws_partition                   = try(each.value.spec.aws_partition, "aws")
  managed_policy_arns             = try(each.value.spec.managed_policy_arns, ["arn:aws:iam::aws:policy/ReadOnlyAccess"])
  github_oidc_provider_host       = try(each.value.spec.github_oidc_provider_host, "token.actions.githubusercontent.com")
  github_oidc_audience            = try(each.value.spec.github_oidc_audience, "sts.amazonaws.com")
  tfe_oidc_provider_host          = try(each.value.spec.tfe_oidc_provider_host, "app.terraform.io")
  tfe_oidc_audience               = try(each.value.spec.tfe_oidc_audience, "aws.workload.identity")
  tfe_project_id                  = try(each.value.spec.tfe_project_id, null)
  tfe_project_name                = try(each.value.spec.tfe_project_name, "*")
  tfe_workspace_execution_mode    = try(each.value.spec.tfe_workspace_execution_mode, "remote")
  tfe_workspace_agent_pool_id     = try(each.value.spec.tfe_workspace_agent_pool_id, null)
  tfe_workspace_terraform_version = try(each.value.spec.tfe_workspace_terraform_version, null)
  tfe_workspace_auto_apply        = try(each.value.spec.tfe_workspace_auto_apply, false)
  tfe_workspace_queue_all_runs    = try(each.value.spec.tfe_workspace_queue_all_runs, true)
  tfe_workspace_speculative_enabled = try(
    each.value.spec.tfe_workspace_speculative_enabled,
    true,
  )
  tfe_workspace_working_directory   = try(each.value.spec.tfe_workspace_working_directory, null)
  tfe_workspace_tags                = try(each.value.spec.tfe_workspace_tags, {})
  tfe_vcs_repo                      = try(each.value.spec.tfe_vcs_repo, null)
  manage_tfe_workspace_variables    = try(each.value.spec.manage_tfe_workspace_variables, true)
  stack_set_name_prefix             = try(each.value.spec.stack_set_name_prefix, null)
  stack_set_permission_model        = try(each.value.spec.stack_set_permission_model, "SELF_MANAGED")
  stack_set_administration_role_arn = var.stack_set_administration_role_arn
  stack_set_execution_role_name     = var.stack_set_execution_role_name
  stack_set_call_as                 = try(each.value.spec.stack_set_call_as, "SELF")
  stack_set_operation_preferences   = try(each.value.spec.stack_set_operation_preferences, null)
  retain_stack_instances_on_destroy = try(each.value.spec.retain_stack_instances_on_destroy, false)
  tags                              = try(each.value.spec.tags, {})
}
