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

  github_tf_repo_manifests = {
    for file_name, manifest in local.manifests_by_file :
    try(manifest.metadata.name, trimsuffix(trimsuffix(basename(file_name), ".yaml"), ".yml")) => manifest
    if try(manifest.kind, "") == "GitHubTerraformRepository"
  }
}
