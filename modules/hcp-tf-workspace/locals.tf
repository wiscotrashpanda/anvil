locals {
  repository_parts = split("/", var.github_repository)
  repository_owner = local.repository_parts[0]
  repository_name  = local.repository_parts[1]

  environment_region = coalesce(var.aws.region, "us-east-1")
  workspace_name     = coalesce(var.workspace.name, "${local.repository_name}-${var.environment}")

  provisioner_role_name = "${local.workspace_name}-provisioner-role"
  provisioner_role_arn  = "arn:${var.aws.partition}:iam::${var.aws.account_id}:role/${local.provisioner_role_name}"

  github_actions_subject = coalesce(
    var.aws.github_actions_subject,
    "repo:${var.github_repository}:*",
  )

  tfe_subject = coalesce(
    var.workspace.hcp_terraform_subject,
    "organization:${tfe_workspace.this.organization}:project:${var.workspace.project_name}:workspace:${local.workspace_name}:run_phase:*",
  )

  stack_set_name_base = var.aws.stack_set_name_prefix == null ? local.workspace_name : "${var.aws.stack_set_name_prefix}-${var.environment}"
  stack_set_name      = "${local.stack_set_name_base}-provisioner-roles"

  common_tags = merge(var.aws.tags, {
    ManagedBy  = "Terraform"
    Module     = "terraform-aws-hcp-tf-workspace"
    Repository = local.repository_name
  })

  use_tfe_vcs_repo = var.workspace.vcs_repo != null
}
