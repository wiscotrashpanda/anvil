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
      try(manifest.spec.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.allowMergeCommit, null) != null || try(manifest.spec.allow_merge_commit, null) != null || try(manifest.spec.repository.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.repository.allowMergeCommit, null) != null || try(manifest.spec.repository.allow_merge_commit, null) != null ? {
        allow_merge_commit = coalesce(
          try(manifest.spec.repository.allow_merge_commit, null),
          try(manifest.spec.repository.allowMergeCommit, null),
          try(manifest.spec.repository.mergePolicy.allowMergeCommit, null),
          try(manifest.spec.allow_merge_commit, null),
          try(manifest.spec.allowMergeCommit, null),
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
        github_repository = try(coalesce(try(manifest.spec.githubRepository, null), try(manifest.spec.github_repository, null)), null)
        environment       = try(manifest.spec.environment, null)
        aws = {
          account_id                        = try(coalesce(try(manifest.spec.aws.accountId, null), try(manifest.spec.aws.account_id, null)), null)
          region                            = try(manifest.spec.aws.region, null)
          partition                         = try(manifest.spec.aws.partition, null)
          managed_policy_arns               = try(coalesce(try(manifest.spec.aws.managedPolicyArns, null), try(manifest.spec.aws.managed_policy_arns, null), try(manifest.spec.managedPolicyArns, null), try(manifest.spec.managed_policy_arns, null)), null)
          github_actions_subject            = try(coalesce(try(manifest.spec.aws.githubActionsSubject, null), try(manifest.spec.aws.github_actions_subject, null), try(manifest.spec.githubActionsSubject, null), try(manifest.spec.github_actions_subject, null)), null)
          github_oidc_provider_host         = try(coalesce(try(manifest.spec.aws.githubOidcProviderHost, null), try(manifest.spec.aws.github_oidc_provider_host, null), try(manifest.spec.githubOidcProviderHost, null), try(manifest.spec.github_oidc_provider_host, null)), null)
          github_oidc_audience              = try(coalesce(try(manifest.spec.aws.githubOidcAudience, null), try(manifest.spec.aws.github_oidc_audience, null), try(manifest.spec.githubOidcAudience, null), try(manifest.spec.github_oidc_audience, null)), null)
          tfe_oidc_provider_host            = try(coalesce(try(manifest.spec.aws.tfeOidcProviderHost, null), try(manifest.spec.aws.tfe_oidc_provider_host, null), try(manifest.spec.tfeOidcProviderHost, null), try(manifest.spec.tfe_oidc_provider_host, null)), null)
          tfe_oidc_audience                 = try(coalesce(try(manifest.spec.aws.tfeOidcAudience, null), try(manifest.spec.aws.tfe_oidc_audience, null), try(manifest.spec.tfeOidcAudience, null), try(manifest.spec.tfe_oidc_audience, null)), null)
          stack_set_name_prefix             = try(coalesce(try(manifest.spec.aws.stackSetNamePrefix, null), try(manifest.spec.aws.stack_set_name_prefix, null), try(manifest.spec.stackSetNamePrefix, null), try(manifest.spec.stack_set_name_prefix, null)), null)
          stack_set_permission_model        = try(coalesce(try(manifest.spec.aws.stackSetPermissionModel, null), try(manifest.spec.aws.stack_set_permission_model, null), try(manifest.spec.stackSetPermissionModel, null), try(manifest.spec.stack_set_permission_model, null)), null)
          stack_set_call_as                 = try(coalesce(try(manifest.spec.aws.stackSetCallAs, null), try(manifest.spec.aws.stack_set_call_as, null), try(manifest.spec.stackSetCallAs, null), try(manifest.spec.stack_set_call_as, null)), null)
          stack_set_operation_preferences   = try(coalesce(try(manifest.spec.aws.stackSetOperationPreferences, null), try(manifest.spec.aws.stack_set_operation_preferences, null), try(manifest.spec.stackSetOperationPreferences, null), try(manifest.spec.stack_set_operation_preferences, null)), null)
          retain_stack_instances_on_destroy = try(coalesce(try(manifest.spec.aws.retainStackInstancesOnDestroy, null), try(manifest.spec.aws.retain_stack_instances_on_destroy, null), try(manifest.spec.retainStackInstancesOnDestroy, null), try(manifest.spec.retain_stack_instances_on_destroy, null)), null)
          tags                              = try(coalesce(try(manifest.spec.aws.tags, null), try(manifest.spec.tags, null)), null)
        }
        workspace = {
          name                  = try(manifest.spec.workspace.name, null)
          project_id            = try(coalesce(try(manifest.spec.workspace.projectId, null), try(manifest.spec.workspace.project_id, null)), null)
          project_name          = try(coalesce(try(manifest.spec.workspace.projectName, null), try(manifest.spec.workspace.project_name, null)), null)
          execution_mode        = try(coalesce(try(manifest.spec.workspace.executionMode, null), try(manifest.spec.workspace.execution_mode, null)), null)
          agent_pool_id         = try(coalesce(try(manifest.spec.workspace.agentPoolId, null), try(manifest.spec.workspace.agent_pool_id, null)), null)
          terraform_version     = try(coalesce(try(manifest.spec.workspace.terraformVersion, null), try(manifest.spec.workspace.terraform_version, null)), null)
          auto_apply            = try(coalesce(try(manifest.spec.workspace.autoApply, null), try(manifest.spec.workspace.auto_apply, null)), null)
          queue_all_runs        = try(coalesce(try(manifest.spec.workspace.queueAllRuns, null), try(manifest.spec.workspace.queue_all_runs, null)), null)
          speculative_enabled   = try(coalesce(try(manifest.spec.workspace.speculativeEnabled, null), try(manifest.spec.workspace.speculative_enabled, null)), null)
          working_directory     = try(coalesce(try(manifest.spec.workspace.workingDirectory, null), try(manifest.spec.workspace.working_directory, null)), null)
          tags                  = try(manifest.spec.workspace.tags, null)
          manage_variables      = try(coalesce(try(manifest.spec.workspace.manageVariables, null), try(manifest.spec.workspace.manage_variables, null)), null)
          hcp_terraform_subject = try(coalesce(try(manifest.spec.workspace.hcpTerraformSubject, null), try(manifest.spec.workspace.hcp_terraform_subject, null)), null)
          vcs_repo = try(coalesce(try(manifest.spec.workspace.vcsRepo, null), try(manifest.spec.workspace.vcs_repo, null)), null) == null ? null : {
            branch                     = try(coalesce(try(manifest.spec.workspace.vcsRepo.branch, null), try(manifest.spec.workspace.vcs_repo.branch, null)), null)
            oauth_token_id             = try(coalesce(try(manifest.spec.workspace.vcsRepo.oauthTokenId, null), try(manifest.spec.workspace.vcs_repo.oauth_token_id, null)), null)
            github_app_installation_id = try(coalesce(try(manifest.spec.workspace.vcsRepo.githubAppInstallationId, null), try(manifest.spec.workspace.vcs_repo.github_app_installation_id, null)), null)
            ingress_submodules         = try(coalesce(try(manifest.spec.workspace.vcsRepo.ingressSubmodules, null), try(manifest.spec.workspace.vcs_repo.ingress_submodules, null)), null)
            trigger_patterns           = try(coalesce(try(manifest.spec.workspace.vcsRepo.triggerPatterns, null), try(manifest.spec.workspace.vcs_repo.trigger_patterns, null)), null)
            trigger_prefixes           = try(coalesce(try(manifest.spec.workspace.vcsRepo.triggerPrefixes, null), try(manifest.spec.workspace.vcs_repo.trigger_prefixes, null)), null)
          }
        }
      },
      try(coalesce(try(manifest.spec.managedPolicyArns, null), try(manifest.spec.managed_policy_arns, null)), null) == null ? {} : {
        managed_policy_arns = try(coalesce(try(manifest.spec.managedPolicyArns, null), try(manifest.spec.managed_policy_arns, null)), null)
      },
      try(coalesce(try(manifest.spec.githubActionsSubject, null), try(manifest.spec.github_actions_subject, null)), null) == null ? {} : {
        github_actions_subject = try(coalesce(try(manifest.spec.githubActionsSubject, null), try(manifest.spec.github_actions_subject, null)), null)
      },
      try(coalesce(try(manifest.spec.githubOidcProviderHost, null), try(manifest.spec.github_oidc_provider_host, null)), null) == null ? {} : {
        github_oidc_provider_host = try(coalesce(try(manifest.spec.githubOidcProviderHost, null), try(manifest.spec.github_oidc_provider_host, null)), null)
      },
      try(coalesce(try(manifest.spec.githubOidcAudience, null), try(manifest.spec.github_oidc_audience, null)), null) == null ? {} : {
        github_oidc_audience = try(coalesce(try(manifest.spec.githubOidcAudience, null), try(manifest.spec.github_oidc_audience, null)), null)
      },
      try(coalesce(try(manifest.spec.tfeOidcProviderHost, null), try(manifest.spec.tfe_oidc_provider_host, null)), null) == null ? {} : {
        tfe_oidc_provider_host = try(coalesce(try(manifest.spec.tfeOidcProviderHost, null), try(manifest.spec.tfe_oidc_provider_host, null)), null)
      },
      try(coalesce(try(manifest.spec.tfeOidcAudience, null), try(manifest.spec.tfe_oidc_audience, null)), null) == null ? {} : {
        tfe_oidc_audience = try(coalesce(try(manifest.spec.tfeOidcAudience, null), try(manifest.spec.tfe_oidc_audience, null)), null)
      },
      try(coalesce(try(manifest.spec.stackSetNamePrefix, null), try(manifest.spec.stack_set_name_prefix, null)), null) == null ? {} : {
        stack_set_name_prefix = try(coalesce(try(manifest.spec.stackSetNamePrefix, null), try(manifest.spec.stack_set_name_prefix, null)), null)
      },
      try(coalesce(try(manifest.spec.stackSetPermissionModel, null), try(manifest.spec.stack_set_permission_model, null)), null) == null ? {} : {
        stack_set_permission_model = try(coalesce(try(manifest.spec.stackSetPermissionModel, null), try(manifest.spec.stack_set_permission_model, null)), null)
      },
      try(coalesce(try(manifest.spec.stackSetCallAs, null), try(manifest.spec.stack_set_call_as, null)), null) == null ? {} : {
        stack_set_call_as = try(coalesce(try(manifest.spec.stackSetCallAs, null), try(manifest.spec.stack_set_call_as, null)), null)
      },
      try(coalesce(try(manifest.spec.stackSetOperationPreferences, null), try(manifest.spec.stack_set_operation_preferences, null)), null) == null ? {} : {
        stack_set_operation_preferences = try(coalesce(try(manifest.spec.stackSetOperationPreferences, null), try(manifest.spec.stack_set_operation_preferences, null)), null)
      },
      try(coalesce(try(manifest.spec.retainStackInstancesOnDestroy, null), try(manifest.spec.retain_stack_instances_on_destroy, null)), null) == null ? {} : {
        retain_stack_instances_on_destroy = try(coalesce(try(manifest.spec.retainStackInstancesOnDestroy, null), try(manifest.spec.retain_stack_instances_on_destroy, null)), null)
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
        try(manifest.spec.repository.mergePolicy.allowMergeCommit, null) != null || try(manifest.spec.repository.allowMergeCommit, null) != null || try(manifest.spec.repository.allow_merge_commit, null) != null ? {
          allow_merge_commit = coalesce(
            try(manifest.spec.repository.allow_merge_commit, null),
            try(manifest.spec.repository.allowMergeCommit, null),
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
              account_id             = try(coalesce(try(config.aws.accountId, null), try(config.aws.account_id, null)), null)
              region                 = try(config.aws.region, null)
              partition              = try(config.aws.partition, null)
              managed_policy_arns    = try(coalesce(try(config.aws.managedPolicyArns, null), try(config.aws.managed_policy_arns, null), try(config.managedPolicyArns, null), try(config.managed_policy_arns, null)), null)
              github_actions_subject = try(coalesce(try(config.aws.githubActionsSubject, null), try(config.aws.github_actions_subject, null), try(config.githubActionsSubject, null), try(config.github_actions_subject, null)), null)
            }
            workspace = try(config.workspace, null) == null ? null : {
              name                  = try(config.workspace.name, null)
              project_id            = try(coalesce(try(config.workspace.projectId, null), try(config.workspace.project_id, null)), null)
              project_name          = try(coalesce(try(config.workspace.projectName, null), try(config.workspace.project_name, null)), null)
              execution_mode        = try(coalesce(try(config.workspace.executionMode, null), try(config.workspace.execution_mode, null)), null)
              agent_pool_id         = try(coalesce(try(config.workspace.agentPoolId, null), try(config.workspace.agent_pool_id, null)), null)
              terraform_version     = try(coalesce(try(config.workspace.terraformVersion, null), try(config.workspace.terraform_version, null)), null)
              auto_apply            = try(coalesce(try(config.workspace.autoApply, null), try(config.workspace.auto_apply, null)), null)
              queue_all_runs        = try(coalesce(try(config.workspace.queueAllRuns, null), try(config.workspace.queue_all_runs, null)), null)
              speculative_enabled   = try(coalesce(try(config.workspace.speculativeEnabled, null), try(config.workspace.speculative_enabled, null)), null)
              working_directory     = try(coalesce(try(config.workspace.workingDirectory, null), try(config.workspace.working_directory, null)), null)
              tags                  = try(config.workspace.tags, null)
              manage_variables      = try(coalesce(try(config.workspace.manageVariables, null), try(config.workspace.manage_variables, null)), null)
              hcp_terraform_subject = try(coalesce(try(config.workspace.hcpTerraformSubject, null), try(config.workspace.hcp_terraform_subject, null)), null)
              vcs_repo = try(coalesce(try(config.workspace.vcsRepo, null), try(config.workspace.vcs_repo, null)), null) == null ? null : {
                branch                     = try(coalesce(try(config.workspace.vcsRepo.branch, null), try(config.workspace.vcs_repo.branch, null)), null)
                oauth_token_id             = try(coalesce(try(config.workspace.vcsRepo.oauthTokenId, null), try(config.workspace.vcs_repo.oauth_token_id, null)), null)
                github_app_installation_id = try(coalesce(try(config.workspace.vcsRepo.githubAppInstallationId, null), try(config.workspace.vcs_repo.github_app_installation_id, null)), null)
                ingress_submodules         = try(coalesce(try(config.workspace.vcsRepo.ingressSubmodules, null), try(config.workspace.vcs_repo.ingress_submodules, null)), null)
                trigger_patterns           = try(coalesce(try(config.workspace.vcsRepo.triggerPatterns, null), try(config.workspace.vcs_repo.trigger_patterns, null)), null)
                trigger_prefixes           = try(coalesce(try(config.workspace.vcsRepo.triggerPrefixes, null), try(config.workspace.vcs_repo.trigger_prefixes, null)), null)
              }
            }
          },
          try(coalesce(try(config.managedPolicyArns, null), try(config.managed_policy_arns, null)), null) == null ? {} : {
            managed_policy_arns = try(coalesce(try(config.managedPolicyArns, null), try(config.managed_policy_arns, null)), null)
          },
          try(coalesce(try(config.githubActionsSubject, null), try(config.github_actions_subject, null)), null) == null ? {} : {
            github_actions_subject = try(coalesce(try(config.githubActionsSubject, null), try(config.github_actions_subject, null)), null)
          },
        )
      }
      aws = {
        region                            = try(manifest.spec.aws.region, null)
        partition                         = try(manifest.spec.aws.partition, null)
        managed_policy_arns               = try(coalesce(try(manifest.spec.aws.managedPolicyArns, null), try(manifest.spec.aws.managed_policy_arns, null), try(manifest.spec.managedPolicyArns, null), try(manifest.spec.managed_policy_arns, null)), null)
        github_oidc_provider_host         = try(coalesce(try(manifest.spec.aws.githubOidcProviderHost, null), try(manifest.spec.aws.github_oidc_provider_host, null), try(manifest.spec.githubOidcProviderHost, null), try(manifest.spec.github_oidc_provider_host, null)), null)
        github_oidc_audience              = try(coalesce(try(manifest.spec.aws.githubOidcAudience, null), try(manifest.spec.aws.github_oidc_audience, null), try(manifest.spec.githubOidcAudience, null), try(manifest.spec.github_oidc_audience, null)), null)
        tfe_oidc_provider_host            = try(coalesce(try(manifest.spec.aws.tfeOidcProviderHost, null), try(manifest.spec.aws.tfe_oidc_provider_host, null), try(manifest.spec.tfeOidcProviderHost, null), try(manifest.spec.tfe_oidc_provider_host, null)), null)
        tfe_oidc_audience                 = try(coalesce(try(manifest.spec.aws.tfeOidcAudience, null), try(manifest.spec.aws.tfe_oidc_audience, null), try(manifest.spec.tfeOidcAudience, null), try(manifest.spec.tfe_oidc_audience, null)), null)
        stack_set_name_prefix             = try(coalesce(try(manifest.spec.aws.stackSetNamePrefix, null), try(manifest.spec.aws.stack_set_name_prefix, null), try(manifest.spec.stackSetNamePrefix, null), try(manifest.spec.stack_set_name_prefix, null)), null)
        stack_set_permission_model        = try(coalesce(try(manifest.spec.aws.stackSetPermissionModel, null), try(manifest.spec.aws.stack_set_permission_model, null), try(manifest.spec.stackSetPermissionModel, null), try(manifest.spec.stack_set_permission_model, null)), null)
        stack_set_call_as                 = try(coalesce(try(manifest.spec.aws.stackSetCallAs, null), try(manifest.spec.aws.stack_set_call_as, null), try(manifest.spec.stackSetCallAs, null), try(manifest.spec.stack_set_call_as, null)), null)
        stack_set_operation_preferences   = try(coalesce(try(manifest.spec.aws.stackSetOperationPreferences, null), try(manifest.spec.aws.stack_set_operation_preferences, null), try(manifest.spec.stackSetOperationPreferences, null), try(manifest.spec.stack_set_operation_preferences, null)), null)
        retain_stack_instances_on_destroy = try(coalesce(try(manifest.spec.aws.retainStackInstancesOnDestroy, null), try(manifest.spec.aws.retain_stack_instances_on_destroy, null), try(manifest.spec.retainStackInstancesOnDestroy, null), try(manifest.spec.retain_stack_instances_on_destroy, null)), null)
        tags                              = try(coalesce(try(manifest.spec.aws.tags, null), try(manifest.spec.tags, null)), null)
      }
      workspace = {
        name                  = try(manifest.spec.workspace.name, null)
        project_id            = try(coalesce(try(manifest.spec.workspace.projectId, null), try(manifest.spec.workspace.project_id, null)), null)
        project_name          = try(coalesce(try(manifest.spec.workspace.projectName, null), try(manifest.spec.workspace.project_name, null)), null)
        execution_mode        = try(coalesce(try(manifest.spec.workspace.executionMode, null), try(manifest.spec.workspace.execution_mode, null)), null)
        agent_pool_id         = try(coalesce(try(manifest.spec.workspace.agentPoolId, null), try(manifest.spec.workspace.agent_pool_id, null)), null)
        terraform_version     = try(coalesce(try(manifest.spec.workspace.terraformVersion, null), try(manifest.spec.workspace.terraform_version, null)), null)
        auto_apply            = try(coalesce(try(manifest.spec.workspace.autoApply, null), try(manifest.spec.workspace.auto_apply, null)), null)
        queue_all_runs        = try(coalesce(try(manifest.spec.workspace.queueAllRuns, null), try(manifest.spec.workspace.queue_all_runs, null)), null)
        speculative_enabled   = try(coalesce(try(manifest.spec.workspace.speculativeEnabled, null), try(manifest.spec.workspace.speculative_enabled, null)), null)
        working_directory     = try(coalesce(try(manifest.spec.workspace.workingDirectory, null), try(manifest.spec.workspace.working_directory, null)), null)
        tags                  = try(manifest.spec.workspace.tags, null)
        manage_variables      = try(coalesce(try(manifest.spec.workspace.manageVariables, null), try(manifest.spec.workspace.manage_variables, null)), null)
        hcp_terraform_subject = try(coalesce(try(manifest.spec.workspace.hcpTerraformSubject, null), try(manifest.spec.workspace.hcp_terraform_subject, null)), null)
        vcs_repo = try(coalesce(try(manifest.spec.workspace.vcsRepo, null), try(manifest.spec.workspace.vcs_repo, null)), null) == null ? null : {
          branch                     = try(coalesce(try(manifest.spec.workspace.vcsRepo.branch, null), try(manifest.spec.workspace.vcs_repo.branch, null)), null)
          oauth_token_id             = try(coalesce(try(manifest.spec.workspace.vcsRepo.oauthTokenId, null), try(manifest.spec.workspace.vcs_repo.oauth_token_id, null)), null)
          github_app_installation_id = try(coalesce(try(manifest.spec.workspace.vcsRepo.githubAppInstallationId, null), try(manifest.spec.workspace.vcs_repo.github_app_installation_id, null)), null)
          ingress_submodules         = try(coalesce(try(manifest.spec.workspace.vcsRepo.ingressSubmodules, null), try(manifest.spec.workspace.vcs_repo.ingress_submodules, null)), null)
          trigger_patterns           = try(coalesce(try(manifest.spec.workspace.vcsRepo.triggerPatterns, null), try(manifest.spec.workspace.vcs_repo.trigger_patterns, null)), null)
          trigger_prefixes           = try(coalesce(try(manifest.spec.workspace.vcsRepo.triggerPrefixes, null), try(manifest.spec.workspace.vcs_repo.trigger_prefixes, null)), null)
        }
      }
      managed_policy_arns               = try(coalesce(try(manifest.spec.managedPolicyArns, null), try(manifest.spec.managed_policy_arns, null)), null)
      github_oidc_provider_host         = try(coalesce(try(manifest.spec.githubOidcProviderHost, null), try(manifest.spec.github_oidc_provider_host, null)), null)
      github_oidc_audience              = try(coalesce(try(manifest.spec.githubOidcAudience, null), try(manifest.spec.github_oidc_audience, null)), null)
      tfe_oidc_provider_host            = try(coalesce(try(manifest.spec.tfeOidcProviderHost, null), try(manifest.spec.tfe_oidc_provider_host, null)), null)
      tfe_oidc_audience                 = try(coalesce(try(manifest.spec.tfeOidcAudience, null), try(manifest.spec.tfe_oidc_audience, null)), null)
      stack_set_name_prefix             = try(coalesce(try(manifest.spec.stackSetNamePrefix, null), try(manifest.spec.stack_set_name_prefix, null)), null)
      stack_set_permission_model        = try(coalesce(try(manifest.spec.stackSetPermissionModel, null), try(manifest.spec.stack_set_permission_model, null)), null)
      stack_set_call_as                 = try(coalesce(try(manifest.spec.stackSetCallAs, null), try(manifest.spec.stack_set_call_as, null)), null)
      stack_set_operation_preferences   = try(coalesce(try(manifest.spec.stackSetOperationPreferences, null), try(manifest.spec.stack_set_operation_preferences, null)), null)
      retain_stack_instances_on_destroy = try(coalesce(try(manifest.spec.retainStackInstancesOnDestroy, null), try(manifest.spec.retain_stack_instances_on_destroy, null)), null)
      tags                              = try(manifest.spec.tags, null)
    }
  }
}
