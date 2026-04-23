# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: standalone GitHub repositories, repo-backed HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` is intentionally left unchanged until its next role is clearer.

## Current Layout

- `manifests/` is the root desired-state input directory. Terraform reads `.yaml` and `.yml` files from there.
- `modules/github-repo/` defines the standalone GitHub-only module built from the reduced `GitHubRepository` manifest surface.
- `modules/github-tf-repo/` defines the first extractable Terraform module for one repo-backed Terraform workload.
- `modules/github-tf-repo/examples/basic/` shows a minimal caller shape.

The standalone GitHub module creates:

- one GitHub repository

The Terraform workload module creates:

- one GitHub repository
- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- one AWS IAM provisioner role per environment trusted by GitHub Actions and HCP Terraform

## Direction

For now, this repository is the design and implementation space for the baseline architecture. Once the module contracts settle, `modules/github-repo` and `modules/github-tf-repo` can be extracted into standalone Terraform module repositories with minimal path churn.

Keep public code, module contracts, and sanitized examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.

## Manifests

Create one private YAML file per desired unit under `manifests/`. Files ending in `.yaml` and `.yml` are ignored by Git so real desired state does not get committed to the public repository.

Use `GitHubRepository` for a standalone non-Terraform GitHub repository:

```yaml
apiVersion: anvil.example.io/v1alpha1
kind: GitHubRepository
metadata:
  name: docs-site
spec:
  repository:
    description: Public documentation site.
    visibility: public
    autoInit: true
    defaultBranch: main
    topics:
      - docs
      - website
    features:
      hasIssues: true
      hasProjects: false
      hasWiki: false
    mergePolicy:
      allowSquashMerge: true
      allowMergeCommit: false
      allowRebaseMerge: true
      deleteBranchOnMerge: true
```

Use `GitHubTerraformRepository` for a repo-backed Terraform workload:

```yaml
apiVersion: anvil.emkaytec.dev/v1alpha1
kind: GitHubTerraformRepository
metadata:
  name: sample-service
spec:
  repository:
    description: Terraform-managed sample service.
    visibility: private
  environments:
    dev:
      account_id: "111111111111"
```

Standalone `GitHubRepository` manifests are translated into the same repository settings object used by `modules/github-tf-repo`, so the underlying GitHub repository defaults and behaviors stay aligned across both Terraform module paths.

Run Terraform from the repo root after placing private manifests in `manifests/`:

```bash
terraform init
terraform plan
```

The root module uses explicit `emkaytec` provider aliases for GitHub and HCP Terraform. Set shared provider ownership once in an ignored root `terraform.tfvars` file. Standalone `GitHubRepository` manifests inherit their owner from `github_owner`.

When you are also planning `GitHubTerraformRepository` manifests, add the HCP Terraform organization and shared StackSet role wiring there too:

```hcl
github_owner     = "emkaytec"

tfe_organization = "emkaytec"
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
