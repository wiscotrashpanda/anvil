variable "github_owner" {
  description = "GitHub owner for the repository."
  type        = string
}

variable "tfe_organization" {
  description = "HCP Terraform organization slug."
  type        = string
}

variable "repository_name" {
  description = "Repository name to create."
  type        = string
}

variable "stack_set_region" {
  description = "AWS region where StackSets are managed."
  type        = string
  default     = "us-east-1"
}

variable "target_account_id" {
  description = "AWS account ID targeted by the example environment."
  type        = string
}

variable "stack_set_administration_role_arn" {
  description = "Administration role ARN required for SELF_MANAGED StackSets."
  type        = string
}

variable "tfe_project_id" {
  description = "HCP Terraform project ID used by the environment-specific workspace override."
  type        = string
}

variable "tfe_agent_pool_id" {
  description = "HCP Terraform agent pool ID used by the environment-specific workspace override."
  type        = string
}

variable "github_app_installation_id" {
  description = "GitHub App installation ID for the HCP Terraform VCS connection."
  type        = number
}
