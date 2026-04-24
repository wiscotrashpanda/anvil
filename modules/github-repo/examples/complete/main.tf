terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.0, < 7.0"
    }

    github = {
      source  = "integrations/github"
      version = ">= 6.0, < 7.0"
    }

    tfe = {
      source  = "hashicorp/tfe"
      version = ">= 0.76.0, < 1.0"
    }
  }
}

provider "aws" {
  region = var.stack_set_region
}

provider "github" {
  alias = "emkaytec"
  owner = var.github_owner
}

provider "tfe" {
  alias        = "emkaytec"
  organization = var.tfe_organization
}

module "repo" {
  source = "../.."

  providers = {
    aws    = aws
    github = github.emkaytec
    tfe    = tfe.emkaytec
  }

  create_terraform_workspaces = true

  repository = {
    name                   = var.repository_name
    description            = "Complete example Terraform-backed GitHub repository."
    visibility             = "private"
    topics                 = ["aws", "terraform", "platform"]
    homepage_url           = "https://example.com/platform/complete-service"
    auto_init              = true
    archive_on_destroy     = true
    has_issues             = true
    has_projects           = false
    has_wiki               = false
    has_discussions        = true
    allow_merge_commit     = true
    allow_squash_merge     = true
    allow_rebase_merge     = false
    delete_branch_on_merge = true
    default_branch         = "main"
    manage_default_branch  = true
    rename_default_branch  = false
  }

  environments = {
    admin = {
      aws = {
        account_id             = var.target_account_id
        region                 = "us-east-2"
        partition              = "aws"
        managed_policy_arns    = ["arn:aws:iam::aws:policy/PowerUserAccess"]
        github_actions_subject = "repo:${var.github_owner}/${var.repository_name}:environment:admin"
      }

      workspace = {
        name                = "complete-service-admin"
        project_id          = var.tfe_project_id
        project_name        = "platform"
        execution_mode      = "agent"
        agent_pool_id       = var.tfe_agent_pool_id
        terraform_version   = "1.14.8"
        auto_apply          = false
        queue_all_runs      = true
        speculative_enabled = true
        working_directory   = "terraform"
        tags = {
          environment = "admin"
          owner       = "platform"
          service     = "complete-service"
        }
        vcs_repo = {
          branch                     = "main"
          oauth_token_id             = null
          github_app_installation_id = var.github_app_installation_id
          ingress_submodules         = false
          trigger_patterns           = ["terraform/**/*.tf"]
          trigger_prefixes           = ["terraform/"]
        }
        manage_variables      = true
        hcp_terraform_subject = "organization:${var.tfe_organization}:project:platform:workspace:complete-service-admin:run_phase:*"
      }
    }
  }

  aws = {
    region                    = "us-east-1"
    partition                 = "aws"
    managed_policy_arns       = ["arn:aws:iam::aws:policy/ReadOnlyAccess"]
    github_oidc_provider_host = "token.actions.githubusercontent.com"
    github_oidc_audience      = "sts.amazonaws.com"
    tfe_oidc_provider_host    = "app.terraform.io"
    tfe_oidc_audience         = "aws.workload.identity"

    stack_set_name_prefix             = "complete-service"
    stack_set_permission_model        = "SELF_MANAGED"
    stack_set_administration_role_arn = var.stack_set_administration_role_arn
    stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
    stack_set_call_as                 = "SELF"
    stack_set_operation_preferences = {
      failure_tolerance_count      = 0
      failure_tolerance_percentage = null
      max_concurrent_count         = 1
      max_concurrent_percentage    = null
      region_concurrency_type      = "SEQUENTIAL"
      region_order                 = ["us-east-2"]
    }
    retain_stack_instances_on_destroy = false
    tags = {
      managed-by = "terraform"
      owner      = "platform"
      service    = "complete-service"
    }
  }

  workspace = {
    name                = null
    project_id          = null
    project_name        = "platform"
    execution_mode      = "remote"
    agent_pool_id       = null
    terraform_version   = "1.14.8"
    auto_apply          = false
    queue_all_runs      = true
    speculative_enabled = true
    working_directory   = "terraform"
    tags = {
      owner   = "platform"
      service = "complete-service"
    }
    vcs_repo = {
      branch                     = "main"
      oauth_token_id             = null
      github_app_installation_id = var.github_app_installation_id
      ingress_submodules         = false
      trigger_patterns           = ["terraform/**/*.tf"]
      trigger_prefixes           = ["terraform/"]
    }
    manage_variables      = true
    hcp_terraform_subject = null
  }
}
