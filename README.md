# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: standalone GitHub repositories, repo-backed HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` is intentionally left unchanged until its next role is clearer.

## Current Layout

- `.forge/` is the root desired-state input directory. Terraform reads `.yaml` and `.yml` files from there. The directory is gitignored and intentionally never committed to this public repository; it is supplied at plan/apply time from a private configuration repository (see [Manifests](#manifests)).
- `modules/github-repo/` defines the GitHub repository module. It can create a standalone repository, or it can also create repo-backed HCP Terraform workspaces and AWS provisioner roles when `create_terraform_workspaces` is enabled.
- `modules/github-repo/examples/basic/` shows a minimal standalone caller shape.
- `modules/github-repo/examples/complete/` shows the full repo-backed HCP Terraform workspace and AWS provisioner role shape.

The GitHub repository module always creates:

- one GitHub repository

When Terraform workspaces are enabled, it also creates:

- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- one AWS IAM provisioner role per environment trusted by GitHub Actions and HCP Terraform
- one GitHub environment with AWS variables per environment

## Direction

For now, this repository is the design and implementation space for the baseline architecture. Once the module contract settles, `modules/github-repo` can be extracted into a standalone `terraform-github-repository` module repository with minimal path churn.

Keep public code, module contracts, and sanitized examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.

## Manifests

Terraform reads desired state from YAML files in the root `.forge/` directory. The directory is gitignored and is not part of this repository; it is populated from a private manifests repository that owns all real account IDs, role ARNs, workspace settings, and other environment-specific values.

In GitHub Actions, check out the private manifests repo at the same level so its `.forge/` directory lands at the root of this working tree and Terraform can read it in place. Locally, populate `.forge/` the same way (for example, by cloning or symlinking the private repo so its `.forge/` is visible at the repo root).

Use one file per desired unit.

Use `GitHubRepository` for a standalone non-Terraform GitHub repository:

```yaml
apiVersion: anvil.emkaytec.dev/v1alpha1
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

Set `createTerraformWorkspaces: true` on `GitHubRepository` for a repo-backed Terraform workload:

```yaml
apiVersion: anvil.emkaytec.dev/v1alpha1
kind: GitHubRepository
metadata:
  name: sample-service
spec:
  createTerraformWorkspaces: true
  repository:
    description: Terraform-managed sample service.
    visibility: private
    topics:
      - aws
      - terraform

  environments:
    dev:
      aws:
        accountId: "111111111111"
    prod:
      aws:
        accountId: "222222222222"
        region: us-east-2
```

`GitHubRepository` manifests inherit their owner from the root `github_owner` value because the Terraform GitHub provider selects one owner per provider configuration. `metadata.name` becomes the Terraform module key and also defaults `spec.repository.name` when omitted.

## Running Terraform

Run Terraform from the repo root once `.forge/` has been populated from the private manifests repository:

```bash
terraform init
terraform plan
```

Provider ownership is configured through the root module's explicit `default` provider aliases. `github_owner` is always required. `tfe_organization` and the shared StackSet role wiring are needed when any `GitHubRepository` creates Terraform workspaces, and they belong in the ignored root `terraform.tfvars` file:

```hcl
github_owner     = "emkaytec"

tfe_organization = "emkaytec"
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
