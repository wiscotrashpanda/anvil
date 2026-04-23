output "repository" {
  description = "Created GitHub repository details."
  value = {
    owner          = split("/", github_repository.this.full_name)[0]
    name           = github_repository.this.name
    full_name      = github_repository.this.full_name
    html_url       = github_repository.this.html_url
    ssh_clone_url  = github_repository.this.ssh_clone_url
    default_branch = try(github_branch_default.this[0].branch, var.repository.defaultBranch, null)
  }
}
