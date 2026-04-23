resource "github_repository" "this" {
  name                   = var.repository.name
  description            = var.repository.description
  visibility             = var.repository.visibility
  homepage_url           = var.repository.homepage_url
  topics                 = var.repository.topics
  auto_init              = var.repository.auto_init
  archive_on_destroy     = var.repository.archive_on_destroy
  has_issues             = var.repository.has_issues
  has_projects           = var.repository.has_projects
  has_wiki               = var.repository.has_wiki
  has_discussions        = var.repository.has_discussions
  allow_merge_commit     = var.repository.allow_merge_commit
  allow_squash_merge     = var.repository.allow_squash_merge
  allow_rebase_merge     = var.repository.allow_rebase_merge
  delete_branch_on_merge = var.repository.delete_branch_on_merge
}

resource "github_branch_default" "this" {
  count = var.repository.manage_default_branch ? 1 : 0

  repository = github_repository.this.name
  branch     = var.repository.default_branch
  rename     = var.repository.rename_default_branch
}
