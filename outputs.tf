output "github_repositories" {
  description = "Standalone GitHub repository outputs keyed by manifest metadata.name."
  value = {
    for name, repo in module.github_repo :
    name => repo.repository
  }
}

output "hcp_tf_workspaces" {
  description = "Standalone HCP Terraform workspace environment outputs keyed by manifest metadata.name."
  value = {
    for name, workspace in module.hcp_tf_workspace :
    name => {
      github_repository = workspace.github_repository
      workspace         = workspace.workspace
      provisioner_role  = workspace.provisioner_role
      stack_set         = workspace.stack_set
    }
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
