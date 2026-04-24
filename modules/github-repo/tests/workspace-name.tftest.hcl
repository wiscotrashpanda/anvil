mock_provider "aws" {}
mock_provider "github" {}
mock_provider "tfe" {}

run "explicit_workspace_name_drives_provisioner_role_name" {
  command = plan

  variables {
    create_terraform_workspaces = true

    repository = {
      name = "sample-service"
    }

    environments = {
      admin = {
        aws = {
          account_id = "111111111111"
        }
        workspace = {
          name = "custom-admin-workspace"
        }
      }
    }

    aws = {
      stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
    }
  }

  assert {
    condition     = output.workspaces.admin.name == "custom-admin-workspace"
    error_message = "Expected the explicit workspace.name to drive the HCP Terraform workspace name."
  }

  assert {
    condition     = output.provisioner_roles.admin.name == "custom-admin-workspace-provisioner-role"
    error_message = "Expected the explicit workspace.name to drive the provisioner role name."
  }

  assert {
    condition     = output.provisioner_roles.admin.arn == "arn:aws:iam::111111111111:role/custom-admin-workspace-provisioner-role"
    error_message = "Expected the provisioner role ARN to include the explicit workspace.name-derived role name."
  }

  assert {
    condition     = output.stack_sets.admin.name == "custom-admin-workspace-provisioner-roles"
    error_message = "Expected the explicit workspace.name to drive the StackSet name."
  }
}
