locals {
  environment_workspace_inputs = {
    for environment, config in var.environments :
    environment => merge(
      var.workspace,
      config.workspace == null ? {} : config.workspace,
      {
        vcs_repo = try(config.workspace.vcs_repo, null) == null && var.workspace.vcs_repo == null ? null : merge(
          var.workspace.vcs_repo == null ? {} : var.workspace.vcs_repo,
          try(config.workspace.vcs_repo, null) == null ? {} : config.workspace.vcs_repo,
          {
            branch = coalesce(
              try(config.workspace.vcs_repo.branch, null),
              try(var.workspace.vcs_repo.branch, null),
              module.github_repo.repository.default_branch,
            )
          },
        )
      },
    )
  }
}
