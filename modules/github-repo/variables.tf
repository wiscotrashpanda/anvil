variable "repository" {
  description = "Standalone GitHub repository settings. Ownership comes from the caller's GitHub provider."
  type = object({
    name          = string
    visibility    = optional(string)
    description   = optional(string)
    homepage      = optional(string)
    autoInit      = optional(bool, false)
    defaultBranch = optional(string)
    topics        = optional(list(string))
    features = optional(object({
      hasIssues   = optional(bool)
      hasProjects = optional(bool)
      hasWiki     = optional(bool)
    }))
    mergePolicy = optional(object({
      allowSquashMerge    = optional(bool)
      allowMergeCommit    = optional(bool)
      allowRebaseMerge    = optional(bool)
      allowAutoMerge      = optional(bool)
      allowUpdateBranch   = optional(bool)
      deleteBranchOnMerge = optional(bool)
    }))
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
    condition = try(var.repository.visibility, null) == null || contains(
      ["public", "private", "internal"],
      var.repository.visibility,
    )
    error_message = "repository.visibility must be public, private, or internal when set."
  }

  validation {
    condition     = try(var.repository.defaultBranch, null) == null || length(trimspace(var.repository.defaultBranch)) > 0
    error_message = "repository.defaultBranch must not be blank when set."
  }

  validation {
    condition = try(var.repository.topics, null) == null || (
      length([for topic in var.repository.topics : topic if trimspace(topic) == ""]) == 0 &&
      length(distinct([for topic in var.repository.topics : lower(trimspace(topic))])) == length(var.repository.topics)
    )
    error_message = "repository.topics must contain unique, non-blank values when set."
  }
}
