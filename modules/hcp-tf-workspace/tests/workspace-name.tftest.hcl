mock_provider "aws" {}
mock_provider "github" {}
mock_provider "tfe" {}

run "explicit_workspace_name_drives_provisioner_role_name" {
  command = plan

  variables {
    github_repository = "emkaytec/sample-service"
    environment       = "admin"
    aws = {
      account_id                        = "111111111111"
      stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
    }

    workspace = {
      name = "custom-admin-workspace"
    }
  }

  assert {
    condition     = output.workspace.name == "custom-admin-workspace"
    error_message = "Expected the explicit workspace.name to drive the HCP Terraform workspace name."
  }

  assert {
    condition     = output.provisioner_role.name == "custom-admin-workspace-provisioner-role"
    error_message = "Expected the explicit workspace.name to drive the provisioner role name."
  }

  assert {
    condition     = output.provisioner_role.arn == "arn:aws:iam::111111111111:role/custom-admin-workspace-provisioner-role"
    error_message = "Expected the provisioner role ARN to include the explicit workspace.name-derived role name."
  }

  assert {
    condition     = output.stack_set.name == "custom-admin-workspace-provisioner-roles"
    error_message = "Expected the explicit workspace.name to drive the StackSet name."
  }
}
