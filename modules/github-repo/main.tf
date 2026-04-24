resource "github_repository" "this" {
  name                   = var.repository.name
  description            = var.repository.description
  visibility             = var.repository.visibility
  homepage_url           = var.repository.homepage_url
  topics                 = var.repository.topics
  auto_init              = var.repository.auto_init
  archive_on_destroy     = var.repository.archive_on_destroy
  has_issues             = var.repository.has_issues
  has_projects           = var.repository.has_projects
  has_wiki               = var.repository.has_wiki
  has_discussions        = var.repository.has_discussions
  allow_merge_commit     = var.repository.allow_merge_commit
  allow_squash_merge     = var.repository.allow_squash_merge
  allow_rebase_merge     = var.repository.allow_rebase_merge
  delete_branch_on_merge = var.repository.delete_branch_on_merge

  lifecycle {
    precondition {
      condition     = !var.create_terraform_workspaces || length(var.environments) > 0
      error_message = "environments must include at least one environment when create_terraform_workspaces is true."
    }

    precondition {
      condition     = !var.create_terraform_workspaces || var.aws.stack_set_permission_model != "SELF_MANAGED" || var.aws.stack_set_administration_role_arn != null
      error_message = "aws.stack_set_administration_role_arn is required when create_terraform_workspaces is true and aws.stack_set_permission_model is SELF_MANAGED."
    }
  }
}

resource "github_branch_default" "this" {
  count = var.repository.manage_default_branch ? 1 : 0

  repository = github_repository.this.name
  branch     = var.repository.default_branch
  rename     = var.repository.rename_default_branch
}

resource "tfe_workspace" "this" {
  for_each = local.terraform_environments

  name                = local.workspace_names[each.key]
  project_id          = local.environment_workspace_inputs[each.key].project_id
  description         = "Terraform workspace for ${github_repository.this.full_name} (${each.key})."
  auto_apply          = local.environment_workspace_inputs[each.key].auto_apply
  queue_all_runs      = local.environment_workspace_inputs[each.key].queue_all_runs
  speculative_enabled = local.environment_workspace_inputs[each.key].speculative_enabled
  terraform_version   = local.environment_workspace_inputs[each.key].terraform_version
  working_directory   = local.environment_workspace_inputs[each.key].working_directory
  tags = merge(local.environment_workspace_inputs[each.key].tags, {
    repository  = github_repository.this.name
    environment = each.key
  })

  trigger_patterns = local.use_tfe_vcs_repo[each.key] ? try(local.environment_workspace_inputs[each.key].vcs_repo.trigger_patterns, null) : null
  trigger_prefixes = local.use_tfe_vcs_repo[each.key] ? try(local.environment_workspace_inputs[each.key].vcs_repo.trigger_prefixes, null) : null

  dynamic "vcs_repo" {
    for_each = local.use_tfe_vcs_repo[each.key] ? [local.environment_workspace_inputs[each.key].vcs_repo] : []

    content {
      identifier                 = github_repository.this.full_name
      branch                     = try(vcs_repo.value.branch, null)
      oauth_token_id             = try(vcs_repo.value.oauth_token_id, null)
      github_app_installation_id = try(vcs_repo.value.github_app_installation_id, null)
      ingress_submodules         = try(vcs_repo.value.ingress_submodules, false)
    }
  }
}

resource "tfe_workspace_settings" "this" {
  for_each = local.terraform_environments

  workspace_id   = tfe_workspace.this[each.key].id
  execution_mode = local.environment_workspace_inputs[each.key].execution_mode
  agent_pool_id  = local.environment_workspace_inputs[each.key].agent_pool_id
}

resource "aws_cloudformation_stack_set" "provisioner_roles" {
  for_each = local.terraform_environments

  name                    = local.stack_set_names[each.key]
  description             = "Provisioner IAM role for ${github_repository.this.full_name} (${each.key})."
  permission_model        = var.aws.stack_set_permission_model
  call_as                 = var.aws.stack_set_call_as
  administration_role_arn = var.aws.stack_set_permission_model == "SELF_MANAGED" ? var.aws.stack_set_administration_role_arn : null
  execution_role_name     = var.aws.stack_set_permission_model == "SELF_MANAGED" ? var.aws.stack_set_execution_role_name : null
  capabilities            = ["CAPABILITY_NAMED_IAM"]
  template_body = templatefile("${path.module}/templates/provisioner-roles.yaml.tftpl", {
    github_oidc_provider_host = var.aws.github_oidc_provider_host
    repository_full_name      = github_repository.this.full_name
    tfe_oidc_provider_host    = var.aws.tfe_oidc_provider_host
  })

  parameters = {
    GitHubOIDCAudience        = var.aws.github_oidc_audience
    GitHubOIDCSubject         = local.github_actions_subjects[each.key]
    HCPTerraformOIDCAudience  = var.aws.tfe_oidc_audience
    HCPTerraformOIDCSubject   = local.tfe_subjects[each.key]
    ManagedPolicyArns         = join(",", local.environment_managed_policy_arns[each.key])
    ProvisionerRoleName       = local.provisioner_role_names[each.key]
    RepositoryFullName        = github_repository.this.full_name
    TerraformWorkspaceName    = local.workspace_names[each.key]
    TerraformOrganizationName = tfe_workspace.this[each.key].organization
    EnvironmentName           = each.key
  }

  tags = merge(local.common_tags[each.key], {
    Environment = each.key
  })

  dynamic "operation_preferences" {
    for_each = var.aws.stack_set_operation_preferences == null ? [] : [var.aws.stack_set_operation_preferences]

    content {
      failure_tolerance_count      = try(operation_preferences.value.failure_tolerance_count, null)
      failure_tolerance_percentage = try(operation_preferences.value.failure_tolerance_percentage, null)
      max_concurrent_count         = try(operation_preferences.value.max_concurrent_count, null)
      max_concurrent_percentage    = try(operation_preferences.value.max_concurrent_percentage, null)
      region_concurrency_type      = try(operation_preferences.value.region_concurrency_type, null)
      region_order                 = try(operation_preferences.value.region_order, null)
    }
  }
}

resource "aws_cloudformation_stack_set_instance" "provisioner_roles" {
  for_each = local.terraform_environments

  stack_set_name            = aws_cloudformation_stack_set.provisioner_roles[each.key].name
  account_id                = each.value.aws.account_id
  call_as                   = var.aws.stack_set_call_as
  stack_set_instance_region = local.environment_regions[each.key]
  retain_stack              = var.aws.retain_stack_instances_on_destroy

  dynamic "operation_preferences" {
    for_each = var.aws.stack_set_operation_preferences == null ? [] : [var.aws.stack_set_operation_preferences]

    content {
      failure_tolerance_count      = try(operation_preferences.value.failure_tolerance_count, null)
      failure_tolerance_percentage = try(operation_preferences.value.failure_tolerance_percentage, null)
      max_concurrent_count         = try(operation_preferences.value.max_concurrent_count, null)
      max_concurrent_percentage    = try(operation_preferences.value.max_concurrent_percentage, null)
      region_concurrency_type      = try(operation_preferences.value.region_concurrency_type, null)
      region_order                 = try(operation_preferences.value.region_order, null)
    }
  }
}

resource "tfe_variable" "account_id" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "account_id"
  value        = local.terraform_environments[each.key].aws.account_id
  category     = "terraform"
  description  = "AWS account ID for this workspace environment."
}

resource "tfe_variable" "aws_region" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "aws_region"
  value        = local.environment_regions[each.key]
  category     = "terraform"
  description  = "AWS region for this workspace environment."
}

resource "tfe_variable" "aws_region_env" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "AWS_REGION"
  value        = local.environment_regions[each.key]
  category     = "env"
  description  = "AWS provider region for dynamic credentials."
}

resource "tfe_variable" "tfc_aws_provider_auth" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_PROVIDER_AUTH"
  value        = "true"
  category     = "env"
  description  = "Enables HCP Terraform AWS dynamic provider credentials."
}

resource "tfe_variable" "tfc_aws_run_role_arn" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_RUN_ROLE_ARN"
  value        = local.provisioner_role_arns[each.key]
  category     = "env"
  description  = "AWS IAM role assumed by HCP Terraform runs."

  depends_on = [aws_cloudformation_stack_set_instance.provisioner_roles]
}

resource "tfe_variable" "tfc_aws_workload_identity_audience" {
  for_each = local.managed_variable_environments

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_WORKLOAD_IDENTITY_AUDIENCE"
  value        = var.aws.tfe_oidc_audience
  category     = "env"
  description  = "OIDC audience expected by the AWS IAM trust policy."
}

resource "github_repository_environment" "this" {
  for_each = local.terraform_environments

  repository  = github_repository.this.name
  environment = each.key
}

resource "github_actions_environment_variable" "aws_region" {
  for_each = local.terraform_environments

  repository    = github_repository.this.name
  environment   = github_repository_environment.this[each.key].environment
  variable_name = "AWS_REGION"
  value         = local.environment_regions[each.key]
}

resource "github_actions_environment_variable" "aws_account_id" {
  for_each = local.terraform_environments

  repository    = github_repository.this.name
  environment   = github_repository_environment.this[each.key].environment
  variable_name = "AWS_ACCOUNT_ID"
  value         = each.value.aws.account_id
}

resource "github_actions_environment_variable" "aws_provisioner_role_arn" {
  for_each = local.terraform_environments

  repository    = github_repository.this.name
  environment   = github_repository_environment.this[each.key].environment
  variable_name = "AWS_PROVISIONER_ROLE_ARN"
  value         = local.provisioner_role_arns[each.key]

  depends_on = [aws_cloudformation_stack_set_instance.provisioner_roles]
}
