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

variable "account_id" {
  description = "AWS account ID targeted by the workspace environment."
  type        = string

  validation {
    condition     = can(regex("^[0-9]{12}$", var.account_id))
    error_message = "account_id must be a 12-digit AWS account ID."
  }
}

variable "region" {
  description = "AWS region for the workspace environment. Defaults to us-east-1 when omitted."
  type        = string
  default     = null
}

variable "workspace_name" {
  description = "Optional explicit HCP Terraform workspace name. Defaults to <repository>-<environment>."
  type        = string
  default     = null
}

variable "managed_policy_arns" {
  description = "Managed policy ARNs attached to the provisioner role."
  type        = list(string)
  default     = ["arn:aws:iam::aws:policy/ReadOnlyAccess"]
}

variable "github_actions_subject" {
  description = "Optional GitHub Actions OIDC subject. Defaults to repo:<github_repository>:*."
  type        = string
  default     = null
}

variable "tfe_subject" {
  description = "Optional HCP Terraform OIDC subject. Defaults to the generated workspace run-phase subject."
  type        = string
  default     = null
}

variable "aws_partition" {
  description = "AWS partition used when computing IAM role ARNs."
  type        = string
  default     = "aws"
}

variable "github_oidc_provider_host" {
  description = "OIDC provider host for GitHub Actions. The provider must already exist in the target AWS account."
  type        = string
  default     = "token.actions.githubusercontent.com"
}

variable "github_oidc_audience" {
  description = "Expected GitHub Actions OIDC audience claim."
  type        = string
  default     = "sts.amazonaws.com"
}

variable "tfe_oidc_provider_host" {
  description = "OIDC provider host for HCP Terraform. The provider must already exist in the target AWS account."
  type        = string
  default     = "app.terraform.io"
}

variable "tfe_oidc_audience" {
  description = "Expected HCP Terraform OIDC audience claim."
  type        = string
  default     = "aws.workload.identity"
}

variable "tfe_project_id" {
  description = "Optional HCP Terraform project ID used to place the workspace."
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
  description = "Optional Terraform version or version constraint for the workspace."
  type        = string
  default     = null
}

variable "tfe_workspace_auto_apply" {
  description = "Whether the workspace auto-applies successful plans."
  type        = bool
  default     = false
}

variable "tfe_workspace_queue_all_runs" {
  description = "Whether the workspace queues runs immediately after creation."
  type        = bool
  default     = true
}

variable "tfe_workspace_speculative_enabled" {
  description = "Whether the workspace allows speculative plans."
  type        = bool
  default     = true
}

variable "tfe_workspace_working_directory" {
  description = "Optional working directory for the workspace."
  type        = string
  default     = null
}

variable "tfe_workspace_tags" {
  description = "Tags applied to the HCP Terraform workspace."
  type        = map(string)
  default     = {}
}

variable "tfe_vcs_repo" {
  description = "Optional HCP Terraform VCS connection. When null, the workspace is API/CLI-driven."
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
  description = "Whether to populate the workspace with account, region, and AWS dynamic credential variables."
  type        = bool
  default     = true
}

variable "stack_set_name_prefix" {
  description = "Optional prefix for the CloudFormation StackSet name. Defaults to the repository name."
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
  description = "Execution role name that must already exist in the target account for SELF_MANAGED StackSets."
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
  description = "Tags applied to the AWS StackSet."
  type        = map(string)
  default     = {}
}
