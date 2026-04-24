output "repository" {
  description = "Created GitHub repository details."
  value       = module.repo.repository
}

output "workspaces" {
  description = "Generated HCP Terraform workspaces keyed by environment."
  value       = module.repo.workspaces
}

output "provisioner_roles" {
  description = "Computed IAM provisioner roles keyed by environment."
  value       = module.repo.provisioner_roles
}

output "stack_sets" {
  description = "CloudFormation StackSet details keyed by environment."
  value       = module.repo.stack_sets
}
