locals {
  manifest_directory = "${path.module}/.forge"

  manifest_files = sort(tolist(setunion(
    fileset(local.manifest_directory, "*.yaml"),
    fileset(local.manifest_directory, "*.yml"),
  )))

  manifests_by_file = {
    for file_name in local.manifest_files :
    file_name => yamldecode(file("${local.manifest_directory}/${file_name}"))
  }

  github_repo_manifests = {
    for file_name, manifest in local.manifests_by_file :
    try(manifest.metadata.name, trimsuffix(trimsuffix(basename(file_name), ".yaml"), ".yml")) => manifest
    if try(manifest.kind, "") == "GitHubRepository"
  }

  github_repo_module_inputs = {
    for name, manifest in local.github_repo_manifests :
    name => merge(
      { name = name },
      try(manifest.spec, {}),
      try(manifest.spec.repository, {}),
      try(manifest.spec.homepage, null) != null || try(manifest.spec.homepage_url, null) != null || try(manifest.spec.repository.homepage, null) != null || try(manifest.spec.repository.homepage_url, null) != null ? {
        homepage_url = coalesce(
          try(manifest.spec.repository.homepage_url, null),
          try(manifest.spec.repository.homepage, null),
          try(manifest.spec.homepage_url, null),
          try(manifest.spec.homepage, null),
        )
      } : {},
      try(manifest.spec.autoInit, null) != null || try(manifest.spec.auto_init, null) != null || try(manifest.spec.repository.autoInit, null) != null || try(manifest.spec.repository.auto_init, null) != null ? {
        auto_init = coalesce(
          try(manifest.spec.repository.auto_init, null),
          try(manifest.spec.repository.autoInit, null),
          try(manifest.spec.auto_init, null),
          try(manifest.spec.autoInit, null),
        )
      } : {},
      try(manifest.spec.defaultBranch, null) != null || try(manifest.spec.default_branch, null) != null || try(manifest.spec.repository.defaultBranch, null) != null || try(manifest.spec.repository.default_branch, null) != null ? {
        default_branch = coalesce(
          try(manifest.spec.repository.default_branch, null),
          try(manifest.spec.repository.defaultBranch, null),
          try(manifest.spec.default_branch, null),
          try(manifest.spec.defaultBranch, null),
        )
      } : {},
      try(manifest.spec.features.hasIssues, null) != null || try(manifest.spec.has_issues, null) != null || try(manifest.spec.repository.features.hasIssues, null) != null || try(manifest.spec.repository.has_issues, null) != null ? {
        has_issues = coalesce(
          try(manifest.spec.repository.has_issues, null),
          try(manifest.spec.repository.features.hasIssues, null),
          try(manifest.spec.has_issues, null),
          try(manifest.spec.features.hasIssues, null),
        )
      } : {},
      try(manifest.spec.features.hasProjects, null) != null || try(manifest.spec.has_projects, null) != null || try(manifest.spec.repository.features.hasProjects, null) != null || try(manifest.spec.repository.has_projects, null) != null ? {
        has_projects = coalesce(
          try(manifest.spec.repository.has_projects, null),
          try(manifest.spec.repository.features.hasProjects, null),
          try(manifest.spec.has_projects, null),
          try(manifest.spec.features.hasProjects, null),
        )
      } : {},
      try(manifest.spec.features.hasWiki, null) != null || try(manifest.spec.has_wiki, null) != null || try(manifest.spec.repository.features.hasWiki, null) != null || try(manifest.spec.repository.has_wiki, null) != null ? {
        has_wiki = coalesce(
          try(manifest.spec.repository.has_wiki, null),
          try(manifest.spec.repository.features.hasWiki, null),
          try(manifest.spec.has_wiki, null),
          try(manifest.spec.features.hasWiki, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.allow_merge_commit, null) != null || try(manifest.spec.repository.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.repository.allow_merge_commit, null) != null ? {
        allow_merge_commit = coalesce(
          try(manifest.spec.repository.allow_merge_commit, null),
          try(manifest.spec.repository.mergePolicy.allowMergeCommit, null),
          try(manifest.spec.allow_merge_commit, null),
          try(manifest.spec.mergePolicy.allowMergeCommit, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.allowSquashMerge, null) != null || try(manifest.spec.allow_squash_merge, null) != null || try(manifest.spec.repository.mergePolicy.allowSquashMerge, null) != null || try(manifest.spec.repository.allow_squash_merge, null) != null ? {
        allow_squash_merge = coalesce(
          try(manifest.spec.repository.allow_squash_merge, null),
          try(manifest.spec.repository.mergePolicy.allowSquashMerge, null),
          try(manifest.spec.allow_squash_merge, null),
          try(manifest.spec.mergePolicy.allowSquashMerge, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.allowRebaseMerge, null) != null || try(manifest.spec.allow_rebase_merge, null) != null || try(manifest.spec.repository.mergePolicy.allowRebaseMerge, null) != null || try(manifest.spec.repository.allow_rebase_merge, null) != null ? {
        allow_rebase_merge = coalesce(
          try(manifest.spec.repository.allow_rebase_merge, null),
          try(manifest.spec.repository.mergePolicy.allowRebaseMerge, null),
          try(manifest.spec.allow_rebase_merge, null),
          try(manifest.spec.mergePolicy.allowRebaseMerge, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.allowAutoMerge, null) != null || try(manifest.spec.allow_auto_merge, null) != null || try(manifest.spec.repository.mergePolicy.allowAutoMerge, null) != null || try(manifest.spec.repository.allow_auto_merge, null) != null ? {
        allow_auto_merge = coalesce(
          try(manifest.spec.repository.allow_auto_merge, null),
          try(manifest.spec.repository.mergePolicy.allowAutoMerge, null),
          try(manifest.spec.allow_auto_merge, null),
          try(manifest.spec.mergePolicy.allowAutoMerge, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.allowUpdateBranch, null) != null || try(manifest.spec.allow_update_branch, null) != null || try(manifest.spec.repository.mergePolicy.allowUpdateBranch, null) != null || try(manifest.spec.repository.allow_update_branch, null) != null ? {
        allow_update_branch = coalesce(
          try(manifest.spec.repository.allow_update_branch, null),
          try(manifest.spec.repository.mergePolicy.allowUpdateBranch, null),
          try(manifest.spec.allow_update_branch, null),
          try(manifest.spec.mergePolicy.allowUpdateBranch, null),
        )
      } : {},
      try(manifest.spec.mergePolicy.deleteBranchOnMerge, null) != null || try(manifest.spec.delete_branch_on_merge, null) != null || try(manifest.spec.repository.mergePolicy.deleteBranchOnMerge, null) != null || try(manifest.spec.repository.delete_branch_on_merge, null) != null ? {
        delete_branch_on_merge = coalesce(
          try(manifest.spec.repository.delete_branch_on_merge, null),
          try(manifest.spec.repository.mergePolicy.deleteBranchOnMerge, null),
          try(manifest.spec.delete_branch_on_merge, null),
          try(manifest.spec.mergePolicy.deleteBranchOnMerge, null),
        )
      } : {},
    )
  }

  hcp_tf_workspace_manifests = {
    for file_name, manifest in local.manifests_by_file :
    try(manifest.metadata.name, trimsuffix(trimsuffix(basename(file_name), ".yaml"), ".yml")) => manifest
    if try(manifest.kind, "") == "HCPTerraformWorkspace"
  }

  hcp_tf_workspace_module_inputs = {
    for name, manifest in local.hcp_tf_workspace_manifests :
    name => merge(
      {
        github_repository = try(coalesce(manifest.spec.githubRepository, manifest.spec.github_repository), null)
        environment       = try(manifest.spec.environment, null)
        aws = {
          account_id                        = try(coalesce(manifest.spec.aws.accountId, manifest.spec.aws.account_id), null)
          region                            = try(manifest.spec.aws.region, null)
          partition                         = try(manifest.spec.aws.partition, null)
          managed_policy_arns               = try(coalesce(manifest.spec.aws.managedPolicyArns, manifest.spec.aws.managed_policy_arns, manifest.spec.managedPolicyArns, manifest.spec.managed_policy_arns), null)
          github_actions_subject            = try(coalesce(manifest.spec.aws.githubActionsSubject, manifest.spec.aws.github_actions_subject, manifest.spec.githubActionsSubject, manifest.spec.github_actions_subject), null)
          github_oidc_provider_host         = try(coalesce(manifest.spec.aws.githubOidcProviderHost, manifest.spec.aws.github_oidc_provider_host, manifest.spec.githubOidcProviderHost, manifest.spec.github_oidc_provider_host), null)
          github_oidc_audience              = try(coalesce(manifest.spec.aws.githubOidcAudience, manifest.spec.aws.github_oidc_audience, manifest.spec.githubOidcAudience, manifest.spec.github_oidc_audience), null)
          tfe_oidc_provider_host            = try(coalesce(manifest.spec.aws.tfeOidcProviderHost, manifest.spec.aws.tfe_oidc_provider_host, manifest.spec.tfeOidcProviderHost, manifest.spec.tfe_oidc_provider_host), null)
          tfe_oidc_audience                 = try(coalesce(manifest.spec.aws.tfeOidcAudience, manifest.spec.aws.tfe_oidc_audience, manifest.spec.tfeOidcAudience, manifest.spec.tfe_oidc_audience), null)
          stack_set_name_prefix             = try(coalesce(manifest.spec.aws.stackSetNamePrefix, manifest.spec.aws.stack_set_name_prefix, manifest.spec.stackSetNamePrefix, manifest.spec.stack_set_name_prefix), null)
          stack_set_permission_model        = try(coalesce(manifest.spec.aws.stackSetPermissionModel, manifest.spec.aws.stack_set_permission_model, manifest.spec.stackSetPermissionModel, manifest.spec.stack_set_permission_model), null)
          stack_set_call_as                 = try(coalesce(manifest.spec.aws.stackSetCallAs, manifest.spec.aws.stack_set_call_as, manifest.spec.stackSetCallAs, manifest.spec.stack_set_call_as), null)
          stack_set_operation_preferences   = try(coalesce(manifest.spec.aws.stackSetOperationPreferences, manifest.spec.aws.stack_set_operation_preferences, manifest.spec.stackSetOperationPreferences, manifest.spec.stack_set_operation_preferences), null)
          retain_stack_instances_on_destroy = try(coalesce(manifest.spec.aws.retainStackInstancesOnDestroy, manifest.spec.aws.retain_stack_instances_on_destroy, manifest.spec.retainStackInstancesOnDestroy, manifest.spec.retain_stack_instances_on_destroy), null)
          tags                              = try(coalesce(manifest.spec.aws.tags, manifest.spec.tags), null)
        }
        workspace = {
          name                  = try(manifest.spec.workspace.name, null)
          project_id            = try(coalesce(manifest.spec.workspace.projectId, manifest.spec.workspace.project_id), null)
          project_name          = try(coalesce(manifest.spec.workspace.projectName, manifest.spec.workspace.project_name), null)
          execution_mode        = try(coalesce(manifest.spec.workspace.executionMode, manifest.spec.workspace.execution_mode), null)
          agent_pool_id         = try(coalesce(manifest.spec.workspace.agentPoolId, manifest.spec.workspace.agent_pool_id), null)
          terraform_version     = try(coalesce(manifest.spec.workspace.terraformVersion, manifest.spec.workspace.terraform_version), null)
          auto_apply            = try(coalesce(manifest.spec.workspace.autoApply, manifest.spec.workspace.auto_apply), null)
          queue_all_runs        = try(coalesce(manifest.spec.workspace.queueAllRuns, manifest.spec.workspace.queue_all_runs), null)
          speculative_enabled   = try(coalesce(manifest.spec.workspace.speculativeEnabled, manifest.spec.workspace.speculative_enabled), null)
          working_directory     = try(coalesce(manifest.spec.workspace.workingDirectory, manifest.spec.workspace.working_directory), null)
          tags                  = try(manifest.spec.workspace.tags, null)
          manage_variables      = try(coalesce(manifest.spec.workspace.manageVariables, manifest.spec.workspace.manage_variables), null)
          hcp_terraform_subject = try(coalesce(manifest.spec.workspace.hcpTerraformSubject, manifest.spec.workspace.hcp_terraform_subject), null)
          vcs_repo = try(coalesce(manifest.spec.workspace.vcsRepo, manifest.spec.workspace.vcs_repo), null) == null ? null : {
            branch                     = try(coalesce(manifest.spec.workspace.vcsRepo.branch, manifest.spec.workspace.vcs_repo.branch), null)
            oauth_token_id             = try(coalesce(manifest.spec.workspace.vcsRepo.oauthTokenId, manifest.spec.workspace.vcs_repo.oauth_token_id), null)
            github_app_installation_id = try(coalesce(manifest.spec.workspace.vcsRepo.githubAppInstallationId, manifest.spec.workspace.vcs_repo.github_app_installation_id), null)
            ingress_submodules         = try(coalesce(manifest.spec.workspace.vcsRepo.ingressSubmodules, manifest.spec.workspace.vcs_repo.ingress_submodules), null)
            trigger_patterns           = try(coalesce(manifest.spec.workspace.vcsRepo.triggerPatterns, manifest.spec.workspace.vcs_repo.trigger_patterns), null)
            trigger_prefixes           = try(coalesce(manifest.spec.workspace.vcsRepo.triggerPrefixes, manifest.spec.workspace.vcs_repo.trigger_prefixes), null)
          }
        }
      },
      try(coalesce(manifest.spec.managedPolicyArns, manifest.spec.managed_policy_arns), null) == null ? {} : {
        managed_policy_arns = try(coalesce(manifest.spec.managedPolicyArns, manifest.spec.managed_policy_arns), null)
      },
      try(coalesce(manifest.spec.githubActionsSubject, manifest.spec.github_actions_subject), null) == null ? {} : {
        github_actions_subject = try(coalesce(manifest.spec.githubActionsSubject, manifest.spec.github_actions_subject), null)
      },
      try(coalesce(manifest.spec.githubOidcProviderHost, manifest.spec.github_oidc_provider_host), null) == null ? {} : {
        github_oidc_provider_host = try(coalesce(manifest.spec.githubOidcProviderHost, manifest.spec.github_oidc_provider_host), null)
      },
      try(coalesce(manifest.spec.githubOidcAudience, manifest.spec.github_oidc_audience), null) == null ? {} : {
        github_oidc_audience = try(coalesce(manifest.spec.githubOidcAudience, manifest.spec.github_oidc_audience), null)
      },
      try(coalesce(manifest.spec.tfeOidcProviderHost, manifest.spec.tfe_oidc_provider_host), null) == null ? {} : {
        tfe_oidc_provider_host = try(coalesce(manifest.spec.tfeOidcProviderHost, manifest.spec.tfe_oidc_provider_host), null)
      },
      try(coalesce(manifest.spec.tfeOidcAudience, manifest.spec.tfe_oidc_audience), null) == null ? {} : {
        tfe_oidc_audience = try(coalesce(manifest.spec.tfeOidcAudience, manifest.spec.tfe_oidc_audience), null)
      },
      try(coalesce(manifest.spec.stackSetNamePrefix, manifest.spec.stack_set_name_prefix), null) == null ? {} : {
        stack_set_name_prefix = try(coalesce(manifest.spec.stackSetNamePrefix, manifest.spec.stack_set_name_prefix), null)
      },
      try(coalesce(manifest.spec.stackSetPermissionModel, manifest.spec.stack_set_permission_model), null) == null ? {} : {
        stack_set_permission_model = try(coalesce(manifest.spec.stackSetPermissionModel, manifest.spec.stack_set_permission_model), null)
      },
      try(coalesce(manifest.spec.stackSetCallAs, manifest.spec.stack_set_call_as), null) == null ? {} : {
        stack_set_call_as = try(coalesce(manifest.spec.stackSetCallAs, manifest.spec.stack_set_call_as), null)
      },
      try(coalesce(manifest.spec.stackSetOperationPreferences, manifest.spec.stack_set_operation_preferences), null) == null ? {} : {
        stack_set_operation_preferences = try(coalesce(manifest.spec.stackSetOperationPreferences, manifest.spec.stack_set_operation_preferences), null)
      },
      try(coalesce(manifest.spec.retainStackInstancesOnDestroy, manifest.spec.retain_stack_instances_on_destroy), null) == null ? {} : {
        retain_stack_instances_on_destroy = try(coalesce(manifest.spec.retainStackInstancesOnDestroy, manifest.spec.retain_stack_instances_on_destroy), null)
      },
      try(manifest.spec.tags, null) == null ? {} : {
        tags = try(manifest.spec.tags, null)
      },
    )
  }

  github_tf_repo_manifests = {
    for file_name, manifest in local.manifests_by_file :
    try(manifest.metadata.name, trimsuffix(trimsuffix(basename(file_name), ".yaml"), ".yml")) => manifest
    if try(manifest.kind, "") == "GitHubTerraformRepository"
  }

  github_tf_repo_module_inputs = {
    for name, manifest in local.github_tf_repo_manifests :
    name => {
      repository = merge(
        { name = name },
        try(manifest.spec.repository, {}),
        try(manifest.spec.repository.homepage, null) != null || try(manifest.spec.repository.homepage_url, null) != null ? {
          homepage_url = coalesce(
            try(manifest.spec.repository.homepage_url, null),
            try(manifest.spec.repository.homepage, null),
          )
        } : {},
        try(manifest.spec.repository.autoInit, null) != null || try(manifest.spec.repository.auto_init, null) != null ? {
          auto_init = coalesce(
            try(manifest.spec.repository.auto_init, null),
            try(manifest.spec.repository.autoInit, null),
          )
        } : {},
        try(manifest.spec.repository.defaultBranch, null) != null || try(manifest.spec.repository.default_branch, null) != null ? {
          default_branch = coalesce(
            try(manifest.spec.repository.default_branch, null),
            try(manifest.spec.repository.defaultBranch, null),
          )
        } : {},
        try(manifest.spec.repository.features.hasIssues, null) != null || try(manifest.spec.repository.has_issues, null) != null ? {
          has_issues = coalesce(
            try(manifest.spec.repository.has_issues, null),
            try(manifest.spec.repository.features.hasIssues, null),
          )
        } : {},
        try(manifest.spec.repository.features.hasProjects, null) != null || try(manifest.spec.repository.has_projects, null) != null ? {
          has_projects = coalesce(
            try(manifest.spec.repository.has_projects, null),
            try(manifest.spec.repository.features.hasProjects, null),
          )
        } : {},
        try(manifest.spec.repository.features.hasWiki, null) != null || try(manifest.spec.repository.has_wiki, null) != null ? {
          has_wiki = coalesce(
            try(manifest.spec.repository.has_wiki, null),
            try(manifest.spec.repository.features.hasWiki, null),
          )
        } : {},
        try(manifest.spec.repository.features.hasDiscussions, null) != null || try(manifest.spec.repository.has_discussions, null) != null ? {
          has_discussions = coalesce(
            try(manifest.spec.repository.has_discussions, null),
            try(manifest.spec.repository.features.hasDiscussions, null),
          )
        } : {},
        try(manifest.spec.repository.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.repository.allow_merge_commit, null) != null ? {
          allow_merge_commit = coalesce(
            try(manifest.spec.repository.allow_merge_commit, null),
            try(manifest.spec.repository.mergePolicy.allowMergeCommit, null),
          )
        } : {},
        try(manifest.spec.repository.mergePolicy.allowSquashMerge, null) != null || try(manifest.spec.repository.allow_squash_merge, null) != null ? {
          allow_squash_merge = coalesce(
            try(manifest.spec.repository.allow_squash_merge, null),
            try(manifest.spec.repository.mergePolicy.allowSquashMerge, null),
          )
        } : {},
        try(manifest.spec.repository.mergePolicy.allowRebaseMerge, null) != null || try(manifest.spec.repository.allow_rebase_merge, null) != null ? {
          allow_rebase_merge = coalesce(
            try(manifest.spec.repository.allow_rebase_merge, null),
            try(manifest.spec.repository.mergePolicy.allowRebaseMerge, null),
          )
        } : {},
        try(manifest.spec.repository.mergePolicy.deleteBranchOnMerge, null) != null || try(manifest.spec.repository.delete_branch_on_merge, null) != null ? {
          delete_branch_on_merge = coalesce(
            try(manifest.spec.repository.delete_branch_on_merge, null),
            try(manifest.spec.repository.mergePolicy.deleteBranchOnMerge, null),
          )
        } : {},
      )
      environments = {
        for environment, config in manifest.spec.environments :
        environment => merge(
          {
            aws = {
              account_id             = try(coalesce(config.aws.accountId, config.aws.account_id), null)
              region                 = try(config.aws.region, null)
              partition              = try(config.aws.partition, null)
              managed_policy_arns    = try(coalesce(config.aws.managedPolicyArns, config.aws.managed_policy_arns, config.managedPolicyArns, config.managed_policy_arns), null)
              github_actions_subject = try(coalesce(config.aws.githubActionsSubject, config.aws.github_actions_subject, config.githubActionsSubject, config.github_actions_subject), null)
            }
            workspace = try(coalesce(config.workspace, null), null) == null ? null : {
              name                  = try(config.workspace.name, null)
              project_id            = try(coalesce(config.workspace.projectId, config.workspace.project_id), null)
              project_name          = try(coalesce(config.workspace.projectName, config.workspace.project_name), null)
              execution_mode        = try(coalesce(config.workspace.executionMode, config.workspace.execution_mode), null)
              agent_pool_id         = try(coalesce(config.workspace.agentPoolId, config.workspace.agent_pool_id), null)
              terraform_version     = try(coalesce(config.workspace.terraformVersion, config.workspace.terraform_version), null)
              auto_apply            = try(coalesce(config.workspace.autoApply, config.workspace.auto_apply), null)
              queue_all_runs        = try(coalesce(config.workspace.queueAllRuns, config.workspace.queue_all_runs), null)
              speculative_enabled   = try(coalesce(config.workspace.speculativeEnabled, config.workspace.speculative_enabled), null)
              working_directory     = try(coalesce(config.workspace.workingDirectory, config.workspace.working_directory), null)
              tags                  = try(config.workspace.tags, null)
              manage_variables      = try(coalesce(config.workspace.manageVariables, config.workspace.manage_variables), null)
              hcp_terraform_subject = try(coalesce(config.workspace.hcpTerraformSubject, config.workspace.hcp_terraform_subject), null)
              vcs_repo = try(coalesce(config.workspace.vcsRepo, config.workspace.vcs_repo), null) == null ? null : {
                branch                     = try(coalesce(config.workspace.vcsRepo.branch, config.workspace.vcs_repo.branch), null)
                oauth_token_id             = try(coalesce(config.workspace.vcsRepo.oauthTokenId, config.workspace.vcs_repo.oauth_token_id), null)
                github_app_installation_id = try(coalesce(config.workspace.vcsRepo.githubAppInstallationId, config.workspace.vcs_repo.github_app_installation_id), null)
                ingress_submodules         = try(coalesce(config.workspace.vcsRepo.ingressSubmodules, config.workspace.vcs_repo.ingress_submodules), null)
                trigger_patterns           = try(coalesce(config.workspace.vcsRepo.triggerPatterns, config.workspace.vcs_repo.trigger_patterns), null)
                trigger_prefixes           = try(coalesce(config.workspace.vcsRepo.triggerPrefixes, config.workspace.vcs_repo.trigger_prefixes), null)
              }
            }
          },
          try(coalesce(config.managedPolicyArns, config.managed_policy_arns), null) == null ? {} : {
            managed_policy_arns = try(coalesce(config.managedPolicyArns, config.managed_policy_arns), null)
          },
          try(coalesce(config.githubActionsSubject, config.github_actions_subject), null) == null ? {} : {
            github_actions_subject = try(coalesce(config.githubActionsSubject, config.github_actions_subject), null)
          },
        )
      }
      aws = {
        region                            = try(manifest.spec.aws.region, null)
        partition                         = try(manifest.spec.aws.partition, null)
        managed_policy_arns               = try(coalesce(manifest.spec.aws.managedPolicyArns, manifest.spec.aws.managed_policy_arns, manifest.spec.managedPolicyArns, manifest.spec.managed_policy_arns), null)
        github_oidc_provider_host         = try(coalesce(manifest.spec.aws.githubOidcProviderHost, manifest.spec.aws.github_oidc_provider_host, manifest.spec.githubOidcProviderHost, manifest.spec.github_oidc_provider_host), null)
        github_oidc_audience              = try(coalesce(manifest.spec.aws.githubOidcAudience, manifest.spec.aws.github_oidc_audience, manifest.spec.githubOidcAudience, manifest.spec.github_oidc_audience), null)
        tfe_oidc_provider_host            = try(coalesce(manifest.spec.aws.tfeOidcProviderHost, manifest.spec.aws.tfe_oidc_provider_host, manifest.spec.tfeOidcProviderHost, manifest.spec.tfe_oidc_provider_host), null)
        tfe_oidc_audience                 = try(coalesce(manifest.spec.aws.tfeOidcAudience, manifest.spec.aws.tfe_oidc_audience, manifest.spec.tfeOidcAudience, manifest.spec.tfe_oidc_audience), null)
        stack_set_name_prefix             = try(coalesce(manifest.spec.aws.stackSetNamePrefix, manifest.spec.aws.stack_set_name_prefix, manifest.spec.stackSetNamePrefix, manifest.spec.stack_set_name_prefix), null)
        stack_set_permission_model        = try(coalesce(manifest.spec.aws.stackSetPermissionModel, manifest.spec.aws.stack_set_permission_model, manifest.spec.stackSetPermissionModel, manifest.spec.stack_set_permission_model), null)
        stack_set_call_as                 = try(coalesce(manifest.spec.aws.stackSetCallAs, manifest.spec.aws.stack_set_call_as, manifest.spec.stackSetCallAs, manifest.spec.stack_set_call_as), null)
        stack_set_operation_preferences   = try(coalesce(manifest.spec.aws.stackSetOperationPreferences, manifest.spec.aws.stack_set_operation_preferences, manifest.spec.stackSetOperationPreferences, manifest.spec.stack_set_operation_preferences), null)
        retain_stack_instances_on_destroy = try(coalesce(manifest.spec.aws.retainStackInstancesOnDestroy, manifest.spec.aws.retain_stack_instances_on_destroy, manifest.spec.retainStackInstancesOnDestroy, manifest.spec.retain_stack_instances_on_destroy), null)
        tags                              = try(coalesce(manifest.spec.aws.tags, manifest.spec.tags), null)
      }
      workspace = try(coalesce(manifest.spec.workspace, null), null) == null ? {} : {
        name                  = try(manifest.spec.workspace.name, null)
        project_id            = try(coalesce(manifest.spec.workspace.projectId, manifest.spec.workspace.project_id), null)
        project_name          = try(coalesce(manifest.spec.workspace.projectName, manifest.spec.workspace.project_name), null)
        execution_mode        = try(coalesce(manifest.spec.workspace.executionMode, manifest.spec.workspace.execution_mode), null)
        agent_pool_id         = try(coalesce(manifest.spec.workspace.agentPoolId, manifest.spec.workspace.agent_pool_id), null)
        terraform_version     = try(coalesce(manifest.spec.workspace.terraformVersion, manifest.spec.workspace.terraform_version), null)
        auto_apply            = try(coalesce(manifest.spec.workspace.autoApply, manifest.spec.workspace.auto_apply), null)
        queue_all_runs        = try(coalesce(manifest.spec.workspace.queueAllRuns, manifest.spec.workspace.queue_all_runs), null)
        speculative_enabled   = try(coalesce(manifest.spec.workspace.speculativeEnabled, manifest.spec.workspace.speculative_enabled), null)
        working_directory     = try(coalesce(manifest.spec.workspace.workingDirectory, manifest.spec.workspace.working_directory), null)
        tags                  = try(manifest.spec.workspace.tags, null)
        manage_variables      = try(coalesce(manifest.spec.workspace.manageVariables, manifest.spec.workspace.manage_variables), null)
        hcp_terraform_subject = try(coalesce(manifest.spec.workspace.hcpTerraformSubject, manifest.spec.workspace.hcp_terraform_subject), null)
        vcs_repo = try(coalesce(manifest.spec.workspace.vcsRepo, manifest.spec.workspace.vcs_repo), null) == null ? null : {
          branch                     = try(coalesce(manifest.spec.workspace.vcsRepo.branch, manifest.spec.workspace.vcs_repo.branch), null)
          oauth_token_id             = try(coalesce(manifest.spec.workspace.vcsRepo.oauthTokenId, manifest.spec.workspace.vcs_repo.oauth_token_id), null)
          github_app_installation_id = try(coalesce(manifest.spec.workspace.vcsRepo.githubAppInstallationId, manifest.spec.workspace.vcs_repo.github_app_installation_id), null)
          ingress_submodules         = try(coalesce(manifest.spec.workspace.vcsRepo.ingressSubmodules, manifest.spec.workspace.vcs_repo.ingress_submodules), null)
          trigger_patterns           = try(coalesce(manifest.spec.workspace.vcsRepo.triggerPatterns, manifest.spec.workspace.vcs_repo.trigger_patterns), null)
          trigger_prefixes           = try(coalesce(manifest.spec.workspace.vcsRepo.triggerPrefixes, manifest.spec.workspace.vcs_repo.trigger_prefixes), null)
        }
      }
      managed_policy_arns               = try(coalesce(manifest.spec.managedPolicyArns, manifest.spec.managed_policy_arns), null)
      github_oidc_provider_host         = try(coalesce(manifest.spec.githubOidcProviderHost, manifest.spec.github_oidc_provider_host), null)
      github_oidc_audience              = try(coalesce(manifest.spec.githubOidcAudience, manifest.spec.github_oidc_audience), null)
      tfe_oidc_provider_host            = try(coalesce(manifest.spec.tfeOidcProviderHost, manifest.spec.tfe_oidc_provider_host), null)
      tfe_oidc_audience                 = try(coalesce(manifest.spec.tfeOidcAudience, manifest.spec.tfe_oidc_audience), null)
      stack_set_name_prefix             = try(coalesce(manifest.spec.stackSetNamePrefix, manifest.spec.stack_set_name_prefix), null)
      stack_set_permission_model        = try(coalesce(manifest.spec.stackSetPermissionModel, manifest.spec.stack_set_permission_model), null)
      stack_set_call_as                 = try(coalesce(manifest.spec.stackSetCallAs, manifest.spec.stack_set_call_as), null)
      stack_set_operation_preferences   = try(coalesce(manifest.spec.stackSetOperationPreferences, manifest.spec.stack_set_operation_preferences), null)
      retain_stack_instances_on_destroy = try(coalesce(manifest.spec.retainStackInstancesOnDestroy, manifest.spec.retain_stack_instances_on_destroy), null)
      tags                              = try(manifest.spec.tags, null)
    }
  }
}
