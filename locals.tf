locals {
  manifest_directory = "${path.module}/manifests"

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

  github_tf_repo_manifests = {
    for file_name, manifest in local.manifests_by_file :
    try(manifest.metadata.name, trimsuffix(trimsuffix(basename(file_name), ".yaml"), ".yml")) => manifest
    if try(manifest.kind, "") == "GitHubTerraformRepository"
  }
}
