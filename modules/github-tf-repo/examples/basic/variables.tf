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

variable "repository_name" {
  description = "Repository name to create."
  type        = string
}

variable "environments" {
  description = "Environment map."
  type = map(object({
    aws = object({
      account_id = string
      region     = optional(string)
    })
  }))
}

variable "stack_set_administration_role_arn" {
  description = "Administration role ARN required for SELF_MANAGED StackSets."
  type        = string
}
