output "github_repositories" {
  description = "Standalone GitHub repository outputs keyed by manifest metadata.name."
  value = {
    for name, repo in module.github_repo :
    name => repo.repository
  }
}

output "github_tf_repositories" {
  description = "Repo-backed Terraform workload outputs keyed by manifest metadata.name."
  value = {
    for name, repo in module.github_tf_repo :
    name => {
      repository        = repo.repository
      workspaces        = repo.workspaces
      provisioner_roles = repo.provisioner_roles
      stack_sets        = repo.stack_sets
    }
  }
}
