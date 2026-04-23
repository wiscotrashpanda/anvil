resource "tfe_workspace" "this" {
  name                = local.workspace_name
  project_id          = var.tfe_project_id
  description         = "Terraform workspace for ${var.github_repository} (${var.environment})."
  auto_apply          = var.tfe_workspace_auto_apply
  queue_all_runs      = var.tfe_workspace_queue_all_runs
  speculative_enabled = var.tfe_workspace_speculative_enabled
  terraform_version   = var.tfe_workspace_terraform_version
  working_directory   = var.tfe_workspace_working_directory
  tags = merge(var.tfe_workspace_tags, {
    repository  = local.repository_name
    environment = var.environment
  })

  trigger_patterns = local.use_tfe_vcs_repo ? try(var.tfe_vcs_repo.trigger_patterns, null) : null
  trigger_prefixes = local.use_tfe_vcs_repo ? try(var.tfe_vcs_repo.trigger_prefixes, null) : null

  dynamic "vcs_repo" {
    for_each = local.use_tfe_vcs_repo ? [var.tfe_vcs_repo] : []

    content {
      identifier                 = var.github_repository
      branch                     = try(vcs_repo.value.branch, null)
      oauth_token_id             = try(vcs_repo.value.oauth_token_id, null)
      github_app_installation_id = try(vcs_repo.value.github_app_installation_id, null)
      ingress_submodules         = try(vcs_repo.value.ingress_submodules, false)
    }
  }
}

resource "tfe_workspace_settings" "this" {
  workspace_id   = tfe_workspace.this.id
  execution_mode = var.tfe_workspace_execution_mode
  agent_pool_id  = var.tfe_workspace_agent_pool_id
}

resource "aws_cloudformation_stack_set" "provisioner_roles" {
  name                    = local.stack_set_name
  description             = "Provisioner IAM role for ${var.github_repository} (${var.environment})."
  permission_model        = var.stack_set_permission_model
  call_as                 = var.stack_set_call_as
  administration_role_arn = var.stack_set_permission_model == "SELF_MANAGED" ? var.stack_set_administration_role_arn : null
  execution_role_name     = var.stack_set_permission_model == "SELF_MANAGED" ? var.stack_set_execution_role_name : null
  capabilities            = ["CAPABILITY_NAMED_IAM"]
  template_body = templatefile("${path.module}/templates/provisioner-roles.yaml.tftpl", {
    github_oidc_provider_host = var.github_oidc_provider_host
    repository_full_name      = var.github_repository
    tfe_oidc_provider_host    = var.tfe_oidc_provider_host
  })

  parameters = {
    GitHubOIDCAudience        = var.github_oidc_audience
    GitHubOIDCSubject         = local.github_actions_subject
    HCPTerraformOIDCAudience  = var.tfe_oidc_audience
    HCPTerraformOIDCSubject   = local.tfe_subject
    ManagedPolicyArns         = join(",", var.managed_policy_arns)
    ProvisionerRoleName       = local.provisioner_role_name
    RepositoryFullName        = var.github_repository
    TerraformWorkspaceName    = local.workspace_name
    TerraformOrganizationName = tfe_workspace.this.organization
    EnvironmentName           = var.environment
  }

  tags = merge(local.common_tags, {
    Environment = var.environment
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
  stack_set_name            = aws_cloudformation_stack_set.provisioner_roles.name
  account_id                = var.account_id
  call_as                   = var.stack_set_call_as
  stack_set_instance_region = local.environment_region
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
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "account_id"
  value        = var.account_id
  category     = "terraform"
  description  = "AWS account ID for this workspace environment."
}

resource "tfe_variable" "aws_region" {
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "aws_region"
  value        = local.environment_region
  category     = "terraform"
  description  = "AWS region for this workspace environment."
}

resource "tfe_variable" "aws_region_env" {
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "AWS_REGION"
  value        = local.environment_region
  category     = "env"
  description  = "AWS provider region for dynamic credentials."
}

resource "tfe_variable" "tfc_aws_provider_auth" {
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "TFC_AWS_PROVIDER_AUTH"
  value        = "true"
  category     = "env"
  description  = "Enables HCP Terraform AWS dynamic provider credentials."
}

resource "tfe_variable" "tfc_aws_run_role_arn" {
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "TFC_AWS_RUN_ROLE_ARN"
  value        = local.provisioner_role_arn
  category     = "env"
  description  = "AWS IAM role assumed by HCP Terraform runs."

  depends_on = [aws_cloudformation_stack_set_instance.provisioner_roles]
}

resource "tfe_variable" "tfc_aws_workload_identity_audience" {
  count = var.manage_tfe_workspace_variables ? 1 : 0

  workspace_id = tfe_workspace.this.id
  key          = "TFC_AWS_WORKLOAD_IDENTITY_AUDIENCE"
  value        = var.tfe_oidc_audience
  category     = "env"
  description  = "OIDC audience expected by the AWS IAM trust policy."
}

resource "github_repository_environment" "this" {
  repository  = local.repository_name
  environment = var.environment
}

resource "github_actions_environment_variable" "aws_region" {
  repository    = local.repository_name
  environment   = github_repository_environment.this.environment
  variable_name = "AWS_REGION"
  value         = local.environment_region
}

resource "github_actions_environment_variable" "aws_account_id" {
  repository    = local.repository_name
  environment   = github_repository_environment.this.environment
  variable_name = "AWS_ACCOUNT_ID"
  value         = var.account_id
}

resource "github_actions_environment_variable" "aws_provisioner_role_arn" {
  repository    = local.repository_name
  environment   = github_repository_environment.this.environment
  variable_name = "AWS_PROVISIONER_ROLE_ARN"
  value         = local.provisioner_role_arn

  depends_on = [aws_cloudformation_stack_set_instance.provisioner_roles]
}
