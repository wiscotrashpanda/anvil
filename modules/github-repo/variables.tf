variable "repository" {
  description = "GitHub repository settings. Ownership comes from the caller's GitHub provider."
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
    condition     = length(trimspace(var.repository.name)) > 0
    error_message = "repository.name must not be empty."
  }

  validation {
    condition     = can(regex("^[A-Za-z0-9_.-]+$", var.repository.name))
    error_message = "repository.name must contain only letters, numbers, dots, underscores, and hyphens."
  }

  validation {
    condition     = contains(["public", "private", "internal"], var.repository.visibility)
    error_message = "repository.visibility must be public, private, or internal."
  }

  validation {
    condition = (
      length([for topic in var.repository.topics : topic if trimspace(topic) == ""]) == 0 &&
      length(distinct([for topic in var.repository.topics : lower(trimspace(topic))])) == length(var.repository.topics)
    )
    error_message = "repository.topics must contain unique, non-blank values."
  }
}

variable "create_terraform_workspaces" {
  description = "Whether to create HCP Terraform workspaces, AWS provisioner roles, and GitHub environments for this repository."
  type        = bool
  default     = false
  nullable    = false
}

variable "environments" {
  description = "Environment map. Required when create_terraform_workspaces is true. Each entry creates one HCP Terraform workspace, one StackSet, and one IAM role."
  type = map(object({
    aws = object({
      account_id             = string
      region                 = optional(string)
      partition              = optional(string)
      managed_policy_arns    = optional(list(string))
      github_actions_subject = optional(string)
    })
    workspace = optional(object({
      name                = optional(string)
      project_id          = optional(string)
      project_name        = optional(string)
      execution_mode      = optional(string)
      agent_pool_id       = optional(string)
      terraform_version   = optional(string)
      auto_apply          = optional(bool)
      queue_all_runs      = optional(bool)
      speculative_enabled = optional(bool)
      working_directory   = optional(string)
      tags                = optional(map(string))
      vcs_repo = optional(object({
        branch                     = optional(string)
        oauth_token_id             = optional(string)
        github_app_installation_id = optional(number)
        ingress_submodules         = optional(bool)
        trigger_patterns           = optional(list(string))
        trigger_prefixes           = optional(list(string))
      }))
      manage_variables      = optional(bool)
      hcp_terraform_subject = optional(string)
    }))
  }))
  default  = {}
  nullable = false

  validation {
    condition = alltrue([
      for environment, _ in var.environments :
      can(regex("^[a-z][a-z0-9-]*$", environment))
    ])
    error_message = "Environment keys must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens."
  }

  validation {
    condition = alltrue([
      for _, config in var.environments :
      can(regex("^[0-9]{12}$", config.aws.account_id))
    ])
    error_message = "Each environment aws.account_id must be a 12-digit AWS account ID."
  }

  validation {
    condition = alltrue([
      for _, config in var.environments :
      config.workspace == null ||
      try(config.workspace.execution_mode, null) == null ||
      contains(["remote", "local", "agent"], config.workspace.execution_mode)
    ])
    error_message = "Each environment workspace.execution_mode must be remote, local, or agent when set."
  }

  validation {
    condition = alltrue([
      for _, config in var.environments :
      try(config.workspace.vcs_repo, null) == null || (
        (try(config.workspace.vcs_repo.oauth_token_id, null) != null) !=
        (try(config.workspace.vcs_repo.github_app_installation_id, null) != null)
      )
    ])
    error_message = "Each environment workspace.vcs_repo must set exactly one of oauth_token_id or github_app_installation_id when set."
  }
}

variable "aws" {
  description = "AWS defaults for generated Terraform workspace environments."
  type = object({
    region                            = optional(string, "us-east-1")
    partition                         = optional(string, "aws")
    managed_policy_arns               = optional(list(string), ["arn:aws:iam::aws:policy/ReadOnlyAccess"])
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
  default  = {}
  nullable = false

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
  description = "Default HCP Terraform workspace settings for generated environments."
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
