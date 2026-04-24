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

  github_repository = module.github_repo.repository.full_name
  environment       = each.key

  aws = {
    account_id                        = each.value.aws.account_id
    region                            = coalesce(try(each.value.aws.region, null), var.aws.region)
    partition                         = coalesce(try(each.value.aws.partition, null), var.aws.partition)
    managed_policy_arns               = coalesce(try(each.value.aws.managed_policy_arns, null), var.aws.managed_policy_arns)
    github_actions_subject            = try(each.value.aws.github_actions_subject, null)
    github_oidc_provider_host         = var.aws.github_oidc_provider_host
    github_oidc_audience              = var.aws.github_oidc_audience
    tfe_oidc_provider_host            = var.aws.tfe_oidc_provider_host
    tfe_oidc_audience                 = var.aws.tfe_oidc_audience
    stack_set_name_prefix             = var.aws.stack_set_name_prefix
    stack_set_permission_model        = var.aws.stack_set_permission_model
    stack_set_administration_role_arn = var.aws.stack_set_administration_role_arn
    stack_set_execution_role_name     = var.aws.stack_set_execution_role_name
    stack_set_call_as                 = var.aws.stack_set_call_as
    stack_set_operation_preferences   = var.aws.stack_set_operation_preferences
    retain_stack_instances_on_destroy = var.aws.retain_stack_instances_on_destroy
    tags                              = var.aws.tags
  }

  workspace = local.environment_workspace_inputs[each.key]
}
