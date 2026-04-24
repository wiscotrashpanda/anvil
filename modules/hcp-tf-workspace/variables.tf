variable "github_repository" {
  description = "GitHub repository path in owner/repository form."
  type        = string

  validation {
    condition     = can(regex("^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$", var.github_repository))
    error_message = "github_repository must be in owner/repository form."
  }
}

variable "environment" {
  description = "Environment name used for workspace, StackSet, IAM role, and GitHub environment naming."
  type        = string

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]*$", var.environment))
    error_message = "environment must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens."
  }
}

variable "aws" {
  description = "AWS target for this workspace environment."
  type = object({
    account_id                        = string
    region                            = optional(string)
    partition                         = optional(string, "aws")
    managed_policy_arns               = optional(list(string), ["arn:aws:iam::aws:policy/ReadOnlyAccess"])
    github_actions_subject            = optional(string)
    github_oidc_provider_host         = optional(string, "token.actions.githubusercontent.com")
    github_oidc_audience              = optional(string, "sts.amazonaws.com")
    tfe_oidc_provider_host            = optional(string, "app.terraform.io")
    tfe_oidc_audience                 = optional(string, "aws.workload.identity")
    stack_set_name_prefix             = optional(string)
    stack_set_permission_model        = optional(string, "SELF_MANAGED")
    stack_set_administration_role_arn = optional(string)
    stack_set_execution_role_name     = optional(string, "AWSCloudFormationStackSetExecutionRole")
    stack_set_call_as                 = optional(string, "SELF")
    stack_set_operation_preferences = optional(object({
      failure_tolerance_count      = optional(number)
      failure_tolerance_percentage = optional(number)
      max_concurrent_count         = optional(number)
      max_concurrent_percentage    = optional(number)
      region_concurrency_type      = optional(string)
      region_order                 = optional(list(string))
    }))
    retain_stack_instances_on_destroy = optional(bool, false)
    tags                              = optional(map(string), {})
  })

  validation {
    condition     = can(regex("^[0-9]{12}$", var.aws.account_id))
    error_message = "aws.account_id must be a 12-digit AWS account ID."
  }

  validation {
    condition     = contains(["SELF_MANAGED", "SERVICE_MANAGED"], var.aws.stack_set_permission_model)
    error_message = "aws.stack_set_permission_model must be SELF_MANAGED or SERVICE_MANAGED."
  }

  validation {
    condition     = contains(["SELF", "DELEGATED_ADMIN"], var.aws.stack_set_call_as)
    error_message = "aws.stack_set_call_as must be SELF or DELEGATED_ADMIN."
  }
}

variable "workspace" {
  description = "HCP Terraform workspace settings."
  type = object({
    name                = optional(string)
    project_id          = optional(string)
    project_name        = optional(string, "*")
    execution_mode      = optional(string, "remote")
    agent_pool_id       = optional(string)
    terraform_version   = optional(string)
    auto_apply          = optional(bool, false)
    queue_all_runs      = optional(bool, true)
    speculative_enabled = optional(bool, true)
    working_directory   = optional(string)
    tags                = optional(map(string), {})
    vcs_repo = optional(object({
      branch                     = optional(string)
      oauth_token_id             = optional(string)
      github_app_installation_id = optional(number)
      ingress_submodules         = optional(bool, false)
      trigger_patterns           = optional(list(string))
      trigger_prefixes           = optional(list(string))
    }))
    manage_variables      = optional(bool, true)
    hcp_terraform_subject = optional(string)
  })
  default  = {}
  nullable = false

  validation {
    condition     = contains(["remote", "local", "agent"], var.workspace.execution_mode)
    error_message = "workspace.execution_mode must be remote, local, or agent."
  }

  validation {
    condition = var.workspace.vcs_repo == null || (
      (try(var.workspace.vcs_repo.oauth_token_id, null) != null) !=
      (try(var.workspace.vcs_repo.github_app_installation_id, null) != null)
    )
    error_message = "workspace.vcs_repo must set exactly one of oauth_token_id or github_app_installation_id."
  }
}
