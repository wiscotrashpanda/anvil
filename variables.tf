variable "github_owner" {
  description = "GitHub user or organization configured on the root GitHub provider."
  type        = string

  validation {
    condition     = length(trimspace(var.github_owner)) > 0
    error_message = "github_owner must not be empty."
  }
}

variable "tfe_organization" {
  description = "HCP Terraform organization configured on the root TFE provider. Required when any GitHubRepository creates Terraform workspaces."
  type        = string
  default     = null

  validation {
    condition     = var.tfe_organization == null || length(trimspace(var.tfe_organization)) > 0
    error_message = "tfe_organization must not be blank when set."
  }
}

variable "stack_set_administration_role_arn" {
  description = "System-wide CloudFormation StackSet administration role ARN for self-managed StackSets."
  type        = string
  default     = null
}

variable "stack_set_execution_role_name" {
  description = "System-wide CloudFormation StackSet execution role name expected in target accounts."
  type        = string
  default     = "AWSCloudFormationStackSetExecutionRole"
}
