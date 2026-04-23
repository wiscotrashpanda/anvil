variable "repository" {
  description = "GitHub repository settings. Defaults intentionally cover the common Terraform repo shape."
  type = object({
    name                   = string
    description            = optional(string, "")
    visibility             = optional(string, "private")
    topics                 = optional(list(string), [])
    homepage_url           = optional(string)
    auto_init              = optional(bool, true)
    archive_on_destroy     = optional(bool, true)
    has_issues             = optional(bool, true)
    has_projects           = optional(bool, false)
    has_wiki               = optional(bool, false)
    has_discussions        = optional(bool, false)
    allow_merge_commit     = optional(bool, false)
    allow_squash_merge     = optional(bool, true)
    allow_rebase_merge     = optional(bool, false)
    delete_branch_on_merge = optional(bool, true)
    default_branch         = optional(string, "main")
    manage_default_branch  = optional(bool, true)
    rename_default_branch  = optional(bool, false)
  })

  validation {
    condition     = can(regex("^[A-Za-z0-9_.-]+$", var.repository.name))
    error_message = "repository.name must contain only letters, numbers, dots, underscores, and hyphens."
  }

  validation {
    condition     = contains(["public", "private", "internal"], var.repository.visibility)
    error_message = "repository.visibility must be public, private, or internal."
  }
}

variable "environments" {
  description = "Environment-to-account map. Each entry creates one TFE workspace, one StackSet, and two IAM roles."
  type = map(object({
    account_id             = string
    region                 = optional(string)
    workspace_name         = optional(string)
    managed_policy_arns    = optional(list(string))
    github_actions_subject = optional(string)
    tfe_subject            = optional(string)
  }))

  validation {
    condition     = length(var.environments) > 0
    error_message = "environments must include at least one environment."
  }

  validation {
    condition = alltrue([
      for environment, config in var.environments :
      can(regex("^[a-z][a-z0-9-]*$", environment))
    ])
    error_message = "Environment keys must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens."
  }

  validation {
    condition = alltrue([
      for _, config in var.environments :
      can(regex("^[0-9]{12}$", config.account_id))
    ])
    error_message = "Each environment account_id must be a 12-digit AWS account ID."
  }
}

variable "default_region" {
  description = "Default AWS region for StackSet instances and workspace AWS region variables."
  type        = string
  default     = "us-east-1"
}

variable "aws_partition" {
  description = "AWS partition used when computing IAM role ARNs for outputs and HCP Terraform variables."
  type        = string
  default     = "aws"
}

variable "managed_policy_arns" {
  description = "Managed policy ARNs attached to both provisioner roles unless an environment overrides them."
  type        = list(string)
  default     = ["arn:aws:iam::aws:policy/ReadOnlyAccess"]
}

variable "github_oidc_provider_host" {
  description = "OIDC provider host for GitHub Actions. The provider must already exist in each target AWS account."
  type        = string
  default     = "token.actions.githubusercontent.com"
}

variable "github_oidc_audience" {
  description = "Expected GitHub Actions OIDC audience claim."
  type        = string
  default     = "sts.amazonaws.com"
}

variable "tfe_oidc_provider_host" {
  description = "OIDC provider host for HCP Terraform. The provider must already exist in each target AWS account."
  type        = string
  default     = "app.terraform.io"
}

variable "tfe_oidc_audience" {
  description = "Expected HCP Terraform OIDC audience claim."
  type        = string
  default     = "aws.workload.identity"
}

variable "tfe_project_id" {
  description = "Optional HCP Terraform project ID used to place the generated workspaces."
  type        = string
  default     = null
}

variable "tfe_project_name" {
  description = "Project name used in HCP Terraform OIDC subject conditions. Use * to wildcard the project segment."
  type        = string
  default     = "*"
}

variable "tfe_workspace_execution_mode" {
  description = "Execution mode applied through tfe_workspace_settings."
  type        = string
  default     = "remote"

  validation {
    condition     = contains(["remote", "local", "agent"], var.tfe_workspace_execution_mode)
    error_message = "tfe_workspace_execution_mode must be remote, local, or agent."
  }
}

variable "tfe_workspace_agent_pool_id" {
  description = "Optional agent pool ID for agent execution mode."
  type        = string
  default     = null
}

variable "tfe_workspace_terraform_version" {
  description = "Optional Terraform version or version constraint for generated workspaces."
  type        = string
  default     = null
}

variable "tfe_workspace_auto_apply" {
  description = "Whether generated workspaces auto-apply successful plans."
  type        = bool
  default     = false
}

variable "tfe_workspace_queue_all_runs" {
  description = "Whether generated workspaces queue runs immediately after creation."
  type        = bool
  default     = true
}

variable "tfe_workspace_speculative_enabled" {
  description = "Whether generated workspaces allow speculative plans."
  type        = bool
  default     = true
}

variable "tfe_workspace_working_directory" {
  description = "Optional working directory for all generated workspaces."
  type        = string
  default     = null
}

variable "tfe_workspace_tags" {
  description = "Tags applied to each generated HCP Terraform workspace."
  type        = map(string)
  default     = {}
}

variable "tfe_vcs_repo" {
  description = "Optional HCP Terraform VCS connection. When null, workspaces are API/CLI-driven."
  type = object({
    branch                     = optional(string)
    oauth_token_id             = optional(string)
    github_app_installation_id = optional(number)
    ingress_submodules         = optional(bool, false)
    trigger_patterns           = optional(list(string))
    trigger_prefixes           = optional(list(string))
  })
  default = null

  validation {
    condition = var.tfe_vcs_repo == null || (
      (try(var.tfe_vcs_repo.oauth_token_id, null) != null) !=
      (try(var.tfe_vcs_repo.github_app_installation_id, null) != null)
    )
    error_message = "tfe_vcs_repo must set exactly one of oauth_token_id or github_app_installation_id."
  }
}

variable "manage_tfe_workspace_variables" {
  description = "Whether to populate generated workspaces with account, region, and AWS dynamic credential variables."
  type        = bool
  default     = true
}

variable "stack_set_name_prefix" {
  description = "Optional prefix for generated CloudFormation StackSet names. Defaults to repository.name."
  type        = string
  default     = null
}

variable "stack_set_permission_model" {
  description = "CloudFormation StackSet permission model."
  type        = string
  default     = "SELF_MANAGED"

  validation {
    condition     = contains(["SELF_MANAGED", "SERVICE_MANAGED"], var.stack_set_permission_model)
    error_message = "stack_set_permission_model must be SELF_MANAGED or SERVICE_MANAGED."
  }
}

variable "stack_set_administration_role_arn" {
  description = "Administration role ARN required for SELF_MANAGED StackSets."
  type        = string
  default     = null
}

variable "stack_set_execution_role_name" {
  description = "Execution role name that must already exist in each target account for SELF_MANAGED StackSets."
  type        = string
  default     = "AWSCloudFormationStackSetExecutionRole"
}

variable "stack_set_call_as" {
  description = "Whether StackSet operations run as the caller account or a delegated administrator."
  type        = string
  default     = "SELF"

  validation {
    condition     = contains(["SELF", "DELEGATED_ADMIN"], var.stack_set_call_as)
    error_message = "stack_set_call_as must be SELF or DELEGATED_ADMIN."
  }
}

variable "stack_set_operation_preferences" {
  description = "Optional CloudFormation StackSet operation preferences."
  type = object({
    failure_tolerance_count      = optional(number)
    failure_tolerance_percentage = optional(number)
    max_concurrent_count         = optional(number)
    max_concurrent_percentage    = optional(number)
    region_concurrency_type      = optional(string)
    region_order                 = optional(list(string))
  })
  default = null
}

variable "retain_stack_instances_on_destroy" {
  description = "Whether Terraform should remove StackSet instances while retaining the underlying stacks and IAM roles."
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags applied to AWS StackSets."
  type        = map(string)
  default     = {}
}
