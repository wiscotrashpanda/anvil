output "github_repositories" {
  description = "GitHub repository outputs keyed by manifest metadata.name."
  value = {
    for name, repo in module.github_repo :
    name => repo.repository
  }
}

output "github_repository_workspaces" {
  description = "Generated HCP Terraform workspace outputs keyed by GitHub repository manifest metadata.name."
  value = {
    for name, repo in module.github_repo :
    name => {
      workspaces        = repo.workspaces
      provisioner_roles = repo.provisioner_roles
      stack_sets        = repo.stack_sets
    }
    if length(repo.workspaces) > 0
  }
}
