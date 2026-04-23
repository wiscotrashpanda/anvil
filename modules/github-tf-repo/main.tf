resource "github_repository" "this" {
  name                   = local.repository_name
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
}

resource "github_branch_default" "this" {
  count = var.repository.manage_default_branch ? 1 : 0

  repository = github_repository.this.name
  branch     = local.default_branch
  rename     = var.repository.rename_default_branch
}

resource "tfe_workspace" "this" {
  for_each = var.environments

  name                = local.workspace_names[each.key]
  organization        = var.tfe_organization
  project_id          = var.tfe_project_id
  description         = "Terraform workspace for ${github_repository.this.full_name} (${each.key})."
  auto_apply          = var.tfe_workspace_auto_apply
  queue_all_runs      = var.tfe_workspace_queue_all_runs
  speculative_enabled = var.tfe_workspace_speculative_enabled
  terraform_version   = var.tfe_workspace_terraform_version
  working_directory   = var.tfe_workspace_working_directory
  tags = merge(var.tfe_workspace_tags, {
    repository  = local.repository_name
    environment = each.key
  })

  trigger_patterns = local.use_tfe_vcs_repo ? try(var.tfe_vcs_repo.trigger_patterns, null) : null
  trigger_prefixes = local.use_tfe_vcs_repo ? try(var.tfe_vcs_repo.trigger_prefixes, null) : null

  dynamic "vcs_repo" {
    for_each = local.use_tfe_vcs_repo ? [var.tfe_vcs_repo] : []

    content {
      identifier                 = github_repository.this.full_name
      branch                     = coalesce(try(vcs_repo.value.branch, null), local.default_branch)
      oauth_token_id             = try(vcs_repo.value.oauth_token_id, null)
      github_app_installation_id = try(vcs_repo.value.github_app_installation_id, null)
      ingress_submodules         = try(vcs_repo.value.ingress_submodules, false)
    }
  }
}

resource "tfe_workspace_settings" "this" {
  for_each = var.environments

  workspace_id   = tfe_workspace.this[each.key].id
  execution_mode = var.tfe_workspace_execution_mode
  agent_pool_id  = var.tfe_workspace_agent_pool_id
}

resource "aws_cloudformation_stack_set" "provisioner_roles" {
  for_each = var.environments

  name                    = local.stack_set_names[each.key]
  description             = "Provisioner IAM roles for ${github_repository.this.full_name} (${each.key})."
  permission_model        = var.stack_set_permission_model
  call_as                 = var.stack_set_call_as
  administration_role_arn = var.stack_set_permission_model == "SELF_MANAGED" ? var.stack_set_administration_role_arn : null
  execution_role_name     = var.stack_set_permission_model == "SELF_MANAGED" ? var.stack_set_execution_role_name : null
  capabilities            = ["CAPABILITY_NAMED_IAM"]
  template_body = templatefile("${path.module}/templates/provisioner-roles.yaml.tftpl", {
    github_oidc_provider_host = var.github_oidc_provider_host
    tfe_oidc_provider_host    = var.tfe_oidc_provider_host
  })

  parameters = {
    GitHubOIDCAudience        = var.github_oidc_audience
    GitHubOIDCSubject         = local.github_actions_subjects[each.key]
    GitHubActionsRoleName     = local.github_actions_role_names[each.key]
    HCPTerraformOIDCAudience  = var.tfe_oidc_audience
    HCPTerraformOIDCSubject   = local.tfe_subjects[each.key]
    HCPTerraformRoleName      = local.tfe_role_names[each.key]
    ManagedPolicyArns         = join(",", local.managed_policy_arns_by_environment[each.key])
    RepositoryFullName        = github_repository.this.full_name
    TerraformWorkspaceName    = local.workspace_names[each.key]
    TerraformOrganizationName = var.tfe_organization
    EnvironmentName           = each.key
  }

  tags = merge(local.common_tags, {
    Environment = each.key
  })

  dynamic "operation_preferences" {
    for_each = var.stack_set_operation_preferences == null ? [] : [var.stack_set_operation_preferences]

    content {
      failure_tolerance_count      = try(operation_preferences.value.failure_tolerance_count, null)
      failure_tolerance_percentage = try(operation_preferences.value.failure_tolerance_percentage, null)
      max_concurrent_count         = try(operation_preferences.value.max_concurrent_count, null)
      max_concurrent_percentage    = try(operation_preferences.value.max_concurrent_percentage, null)
      region_concurrency_type      = try(operation_preferences.value.region_concurrency_type, null)
      region_order                 = try(operation_preferences.value.region_order, null)
    }
  }

  lifecycle {
    precondition {
      condition     = var.stack_set_permission_model != "SELF_MANAGED" || var.stack_set_administration_role_arn != null
      error_message = "stack_set_administration_role_arn is required when stack_set_permission_model is SELF_MANAGED."
    }
  }
}

resource "aws_cloudformation_stack_set_instance" "provisioner_roles" {
  for_each = var.environments

  stack_set_name            = aws_cloudformation_stack_set.provisioner_roles[each.key].name
  account_id                = each.value.account_id
  call_as                   = var.stack_set_call_as
  stack_set_instance_region = local.environment_regions[each.key]
  retain_stack              = var.retain_stack_instances_on_destroy

  dynamic "operation_preferences" {
    for_each = var.stack_set_operation_preferences == null ? [] : [var.stack_set_operation_preferences]

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
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "account_id"
  value        = each.value.account_id
  category     = "terraform"
  description  = "AWS account ID for this workspace environment."
}

resource "tfe_variable" "aws_region" {
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "aws_region"
  value        = local.environment_regions[each.key]
  category     = "terraform"
  description  = "AWS region for this workspace environment."
}

resource "tfe_variable" "aws_region_env" {
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "AWS_REGION"
  value        = local.environment_regions[each.key]
  category     = "env"
  description  = "AWS provider region for dynamic credentials."
}

resource "tfe_variable" "tfc_aws_provider_auth" {
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_PROVIDER_AUTH"
  value        = "true"
  category     = "env"
  description  = "Enables HCP Terraform AWS dynamic provider credentials."
}

resource "tfe_variable" "tfc_aws_run_role_arn" {
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_RUN_ROLE_ARN"
  value        = local.tfe_role_arns[each.key]
  category     = "env"
  description  = "AWS IAM role assumed by HCP Terraform runs."

  depends_on = [aws_cloudformation_stack_set_instance.provisioner_roles]
}

resource "tfe_variable" "tfc_aws_workload_identity_audience" {
  for_each = var.manage_tfe_workspace_variables ? var.environments : {}

  workspace_id = tfe_workspace.this[each.key].id
  key          = "TFC_AWS_WORKLOAD_IDENTITY_AUDIENCE"
  value        = var.tfe_oidc_audience
  category     = "env"
  description  = "OIDC audience expected by the AWS IAM trust policy."
}
