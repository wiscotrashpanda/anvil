output "repository" {
  description = "Created GitHub repository details."
  value       = module.github_repo.repository
}

output "workspaces" {
  description = "Generated HCP Terraform workspaces keyed by environment."
  value = {
    for environment, workspace in module.hcp_tf_workspace :
    environment => workspace.workspace
  }
}

output "provisioner_roles" {
  description = "Computed IAM provisioner role names, ARNs, and OIDC subjects keyed by environment."
  value = {
    for environment, workspace in module.hcp_tf_workspace :
    environment => workspace.provisioner_role
  }
}

output "stack_sets" {
  description = "CloudFormation StackSet details keyed by environment."
  value = {
    for environment, workspace in module.hcp_tf_workspace :
    environment => workspace.stack_set
  }
}
