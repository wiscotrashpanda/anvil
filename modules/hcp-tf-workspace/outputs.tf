output "github_repository" {
  description = "GitHub repository path this workspace module targets."
  value = {
    owner     = local.repository_owner
    name      = local.repository_name
    full_name = var.github_repository
  }
}

output "workspace" {
  description = "Generated HCP Terraform workspace details."
  value = {
    id         = tfe_workspace.this.id
    name       = tfe_workspace.this.name
    html_url   = tfe_workspace.this.html_url
    account_id = var.aws.account_id
    region     = local.environment_region
  }
}

output "provisioner_role" {
  description = "Computed IAM provisioner role details and OIDC subjects."
  value = {
    name                   = local.provisioner_role_name
    arn                    = local.provisioner_role_arn
    github_actions_subject = local.github_actions_subject
    hcp_terraform_subject  = local.tfe_subject
  }
}

output "stack_set" {
  description = "CloudFormation StackSet details."
  value = {
    name         = aws_cloudformation_stack_set.provisioner_roles.name
    arn          = aws_cloudformation_stack_set.provisioner_roles.arn
    stack_set_id = aws_cloudformation_stack_set.provisioner_roles.stack_set_id
    instance_id  = aws_cloudformation_stack_set_instance.provisioner_roles.id
  }
}
