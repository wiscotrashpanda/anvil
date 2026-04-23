variable "stack_set_region" {
  description = "AWS region where StackSets are managed."
  type        = string
  default     = "us-east-1"
}

variable "github_owner" {
  description = "GitHub owner for the repository."
  type        = string
}

variable "tfe_organization" {
  description = "HCP Terraform organization slug."
  type        = string
}

variable "github_repository" {
  description = "Repository path in owner/repository form."
  type        = string
}

variable "environment" {
  description = "Workspace environment name."
  type        = string
}

variable "account_id" {
  description = "Target AWS account ID."
  type        = string
}

variable "stack_set_administration_role_arn" {
  description = "Administration role ARN required for SELF_MANAGED StackSets."
  type        = string
}
