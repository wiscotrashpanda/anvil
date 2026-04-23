locals {
  repository_name = var.repository.name
  default_branch  = var.repository.default_branch

  stack_set_name_prefix = coalesce(var.stack_set_name_prefix, local.repository_name)

  environment_regions = {
    for environment, config in var.environments :
    environment => coalesce(try(config.region, null), var.default_region)
  }

  workspace_names = {
    for environment, config in var.environments :
    environment => coalesce(try(config.workspace_name, null), "${local.repository_name}-${environment}")
  }

  managed_policy_arns_by_environment = {
    for environment, config in var.environments :
    environment => try(config.managed_policy_arns, null) != null ? config.managed_policy_arns : var.managed_policy_arns
  }

  provisioner_role_names = {
    for environment in keys(var.environments) :
    environment => "${local.repository_name}-${environment}-provisioner-role"
  }

  provisioner_role_arns = {
    for environment, config in var.environments :
    environment => "arn:${var.aws_partition}:iam::${config.account_id}:role/${local.provisioner_role_names[environment]}"
  }

  github_actions_subjects = {
    for environment, config in var.environments :
    environment => coalesce(
      try(config.github_actions_subject, null),
      "repo:${module.github_repo.repository.full_name}:*"
    )
  }

  tfe_subjects = {
    for environment, config in var.environments :
    environment => coalesce(
      try(config.tfe_subject, null),
      "organization:${tfe_workspace.this[environment].organization}:project:${var.tfe_project_name}:workspace:${local.workspace_names[environment]}:run_phase:*"
    )
  }

  stack_set_names = {
    for environment in keys(var.environments) :
    environment => "${local.stack_set_name_prefix}-${environment}-provisioner-roles"
  }

  common_tags = merge(var.tags, {
    ManagedBy  = "Terraform"
    Module     = "terraform-aws-github-tf-repo"
    Repository = local.repository_name
  })

  use_tfe_vcs_repo = var.tfe_vcs_repo != null
}
