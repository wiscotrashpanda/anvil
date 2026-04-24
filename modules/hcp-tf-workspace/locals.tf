locals {
  repository_parts = split("/", var.github_repository)
  repository_owner = local.repository_parts[0]
  repository_name  = local.repository_parts[1]

  environment_region = coalesce(var.region, "us-east-1")
  workspace_name     = coalesce(var.workspace_name, "${local.repository_name}-${var.environment}")

  provisioner_role_name = "${local.workspace_name}-provisioner-role"
  provisioner_role_arn  = "arn:${var.aws_partition}:iam::${var.account_id}:role/${local.provisioner_role_name}"

  github_actions_subject = coalesce(
    var.github_actions_subject,
    "repo:${var.github_repository}:*",
  )

  tfe_subject = coalesce(
    var.tfe_subject,
    "organization:${tfe_workspace.this.organization}:project:${var.tfe_project_name}:workspace:${local.workspace_name}:run_phase:*",
  )

  stack_set_name = "${coalesce(var.stack_set_name_prefix, local.repository_name)}-${var.environment}-provisioner-roles"

  common_tags = merge(var.tags, {
    ManagedBy  = "Terraform"
    Module     = "terraform-aws-hcp-tf-workspace"
    Repository = local.repository_name
  })

  use_tfe_vcs_repo = var.tfe_vcs_repo != null
}
