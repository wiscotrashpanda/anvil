output "repository" {
  description = "Created GitHub repository details."
  value = {
    name           = github_repository.this.name
    full_name      = github_repository.this.full_name
    html_url       = github_repository.this.html_url
    ssh_clone_url  = github_repository.this.ssh_clone_url
    default_branch = try(github_branch_default.this[0].branch, var.repository.default_branch)
  }
}

output "workspaces" {
  description = "Generated HCP Terraform workspaces keyed by environment."
  value = {
    for environment, workspace in tfe_workspace.this :
    environment => {
      id         = workspace.id
      name       = workspace.name
      html_url   = workspace.html_url
      account_id = local.terraform_environments[environment].aws.account_id
      region     = local.environment_regions[environment]
    }
  }
}

output "provisioner_roles" {
  description = "Computed IAM provisioner role names, ARNs, and OIDC subjects keyed by environment."
  value = {
    for environment, _ in local.terraform_environments :
    environment => {
      name                   = local.provisioner_role_names[environment]
      arn                    = local.provisioner_role_arns[environment]
      github_actions_subject = local.github_actions_subjects[environment]
      hcp_terraform_subject  = local.tfe_subjects[environment]
    }
  }
}

output "stack_sets" {
  description = "CloudFormation StackSet details keyed by environment."
  value = {
    for environment, stack_set in aws_cloudformation_stack_set.provisioner_roles :
    environment => {
      name         = stack_set.name
      arn          = stack_set.arn
      stack_set_id = stack_set.stack_set_id
      instance_id  = aws_cloudformation_stack_set_instance.provisioner_roles[environment].id
    }
  }
}
