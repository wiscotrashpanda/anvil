module "github_repo" {
  for_each = local.managed_github_repo_module_inputs

  source  = "app.terraform.io/emkaytec/repository/github"
  version = "0.0.1"

  providers = {
    aws    = aws.default
    github = github.default
    tfe    = tfe.default
  }

  repository                  = each.value.repository
  create_terraform_workspaces = each.value.create_terraform_workspaces
  environments                = each.value.environments

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
