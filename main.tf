module "github_repo" {
  for_each = local.github_repo_manifests

  source = "./modules/github-repo"

  providers = {
    github = github.default
  }

  repository = local.github_repo_module_inputs[each.key]
}

module "hcp_tf_workspace" {
  for_each = local.hcp_tf_workspace_module_inputs

  source = "./modules/hcp-tf-workspace"

  providers = {
    aws    = aws.default
    github = github.default
    tfe    = tfe.default
  }

  github_repository = each.value.github_repository
  environment       = each.value.environment

  aws = {
    account_id                        = each.value.aws.account_id
    region                            = each.value.aws.region
    partition                         = coalesce(each.value.aws.partition, "aws")
    managed_policy_arns               = coalesce(each.value.aws.managed_policy_arns, ["arn:aws:iam::aws:policy/ReadOnlyAccess"])
    github_actions_subject            = each.value.aws.github_actions_subject
    github_oidc_provider_host         = coalesce(each.value.aws.github_oidc_provider_host, "token.actions.githubusercontent.com")
    github_oidc_audience              = coalesce(each.value.aws.github_oidc_audience, "sts.amazonaws.com")
    tfe_oidc_provider_host            = coalesce(each.value.aws.tfe_oidc_provider_host, "app.terraform.io")
    tfe_oidc_audience                 = coalesce(each.value.aws.tfe_oidc_audience, "aws.workload.identity")
    stack_set_name_prefix             = each.value.aws.stack_set_name_prefix
    stack_set_permission_model        = coalesce(each.value.aws.stack_set_permission_model, "SELF_MANAGED")
    stack_set_administration_role_arn = var.stack_set_administration_role_arn
    stack_set_execution_role_name     = var.stack_set_execution_role_name
    stack_set_call_as                 = coalesce(each.value.aws.stack_set_call_as, "SELF")
    stack_set_operation_preferences   = each.value.aws.stack_set_operation_preferences
    retain_stack_instances_on_destroy = coalesce(each.value.aws.retain_stack_instances_on_destroy, false)
    tags                              = coalesce(each.value.aws.tags, {})
  }

  workspace = try(each.value.workspace, null) == null ? {} : each.value.workspace
}

module "github_tf_repo" {
  for_each = local.github_tf_repo_module_inputs

  source = "./modules/github-tf-repo"

  providers = {
    aws    = aws.default
    github = github.default
    tfe    = tfe.default
  }

  repository   = each.value.repository
  environments = each.value.environments

  aws = {
    region                            = coalesce(each.value.aws.region, "us-east-1")
    partition                         = coalesce(each.value.aws.partition, "aws")
    managed_policy_arns               = coalesce(each.value.aws.managed_policy_arns, ["arn:aws:iam::aws:policy/ReadOnlyAccess"])
    github_oidc_provider_host         = coalesce(each.value.aws.github_oidc_provider_host, "token.actions.githubusercontent.com")
    github_oidc_audience              = coalesce(each.value.aws.github_oidc_audience, "sts.amazonaws.com")
    tfe_oidc_provider_host            = coalesce(each.value.aws.tfe_oidc_provider_host, "app.terraform.io")
    tfe_oidc_audience                 = coalesce(each.value.aws.tfe_oidc_audience, "aws.workload.identity")
    stack_set_name_prefix             = each.value.aws.stack_set_name_prefix
    stack_set_permission_model        = coalesce(each.value.aws.stack_set_permission_model, "SELF_MANAGED")
    stack_set_administration_role_arn = var.stack_set_administration_role_arn
    stack_set_execution_role_name     = var.stack_set_execution_role_name
    stack_set_call_as                 = coalesce(each.value.aws.stack_set_call_as, "SELF")
    stack_set_operation_preferences   = each.value.aws.stack_set_operation_preferences
    retain_stack_instances_on_destroy = coalesce(each.value.aws.retain_stack_instances_on_destroy, false)
    tags                              = coalesce(each.value.aws.tags, {})
  }
  workspace = each.value.workspace
}
