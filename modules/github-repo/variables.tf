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
