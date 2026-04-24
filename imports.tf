moved {
  from = module.github_tf_repo["anvil"].module.github_repo.github_repository.this
  to   = module.github_repo["anvil"].github_repository.this
}

moved {
  from = module.github_tf_repo["anvil"].module.github_repo.github_branch_default.this[0]
  to   = module.github_repo["anvil"].github_branch_default.this[0]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_workspace.this
  to   = module.github_repo["anvil"].tfe_workspace.this["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_workspace_settings.this
  to   = module.github_repo["anvil"].tfe_workspace_settings.this["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].aws_cloudformation_stack_set.provisioner_roles
  to   = module.github_repo["anvil"].aws_cloudformation_stack_set.provisioner_roles["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].aws_cloudformation_stack_set_instance.provisioner_roles
  to   = module.github_repo["anvil"].aws_cloudformation_stack_set_instance.provisioner_roles["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.account_id[0]
  to   = module.github_repo["anvil"].tfe_variable.account_id["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.aws_region[0]
  to   = module.github_repo["anvil"].tfe_variable.aws_region["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.aws_region_env[0]
  to   = module.github_repo["anvil"].tfe_variable.aws_region_env["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.tfc_aws_provider_auth[0]
  to   = module.github_repo["anvil"].tfe_variable.tfc_aws_provider_auth["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.tfc_aws_run_role_arn[0]
  to   = module.github_repo["anvil"].tfe_variable.tfc_aws_run_role_arn["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].tfe_variable.tfc_aws_workload_identity_audience[0]
  to   = module.github_repo["anvil"].tfe_variable.tfc_aws_workload_identity_audience["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].github_repository_environment.this
  to   = module.github_repo["anvil"].github_repository_environment.this["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].github_actions_environment_variable.aws_region
  to   = module.github_repo["anvil"].github_actions_environment_variable.aws_region["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].github_actions_environment_variable.aws_account_id
  to   = module.github_repo["anvil"].github_actions_environment_variable.aws_account_id["admin"]
}

moved {
  from = module.github_tf_repo["anvil"].module.hcp_tf_workspace["admin"].github_actions_environment_variable.aws_provisioner_role_arn
  to   = module.github_repo["anvil"].github_actions_environment_variable.aws_provisioner_role_arn["admin"]
}
