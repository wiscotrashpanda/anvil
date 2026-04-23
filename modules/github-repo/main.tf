resource "github_repository" "this" {
  name                   = var.repository.name
  description            = try(var.repository.description, null)
  visibility             = coalesce(try(var.repository.visibility, null), "private")
  homepage_url           = try(var.repository.homepage, null)
  auto_init              = try(var.repository.autoInit, false)
  topics                 = try(var.repository.topics, null)
  has_issues             = try(var.repository.features.hasIssues, null)
  has_projects           = try(var.repository.features.hasProjects, null)
  has_wiki               = try(var.repository.features.hasWiki, null)
  allow_squash_merge     = try(var.repository.mergePolicy.allowSquashMerge, null)
  allow_merge_commit     = try(var.repository.mergePolicy.allowMergeCommit, null)
  allow_rebase_merge     = try(var.repository.mergePolicy.allowRebaseMerge, null)
  allow_auto_merge       = try(var.repository.mergePolicy.allowAutoMerge, null)
  allow_update_branch    = try(var.repository.mergePolicy.allowUpdateBranch, null)
  delete_branch_on_merge = try(var.repository.mergePolicy.deleteBranchOnMerge, null)
}

resource "github_branch_default" "this" {
  count = try(var.repository.defaultBranch, null) == null ? 0 : 1

  repository = github_repository.this.name
  branch     = var.repository.defaultBranch
  rename     = true
}
