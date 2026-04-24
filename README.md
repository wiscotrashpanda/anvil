# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: standalone GitHub repositories, repo-backed HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` is intentionally left unchanged until its next role is clearer.

## Current Layout

- `.forge/` is the root desired-state input directory. Terraform reads `.yaml` and `.yml` files from there. The directory is gitignored and intentionally never committed to this public repository; it is supplied at plan/apply time from a private configuration repository (see [Manifests](#manifests)).
- `modules/github-repo/` defines the shared GitHub repository creation module.
- `modules/hcp-tf-workspace/` defines the reusable workspace-and-role module for one repo-backed HCP Terraform environment.
- `modules/github-tf-repo/` defines the first extractable Terraform module for one repo-backed Terraform workload.
- `modules/github-tf-repo/examples/basic/` shows a minimal caller shape.

The standalone GitHub module creates:

- one GitHub repository

The standalone HCP Terraform environment module creates:

- one HCP Terraform workspace
- one CloudFormation StackSet
- one AWS IAM provisioner role
- one GitHub environment with AWS variables for that environment

The Terraform workload module creates:

- one GitHub repository
- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- one AWS IAM provisioner role per environment trusted by GitHub Actions and HCP Terraform

## Direction

For now, this repository is the design and implementation space for the baseline architecture. Once the module contracts settle, `modules/github-repo`, `modules/hcp-tf-workspace`, and `modules/github-tf-repo` can be extracted into standalone Terraform module repositories with minimal path churn.

Keep public code, module contracts, and sanitized examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.

## Manifests

Terraform reads desired state from YAML files in the root `.forge/` directory. The directory is gitignored and is not part of this repository; it is populated from a private manifests repository that owns all real account IDs, role ARNs, workspace settings, and other environment-specific values.

In GitHub Actions, check out the private manifests repo at the same level so its `.forge/` directory lands at the root of this working tree and Terraform can read it in place. Locally, populate `.forge/` the same way (for example, by cloning or symlinking the private repo so its `.forge/` is visible at the repo root).

Use one file per desired unit.

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

Standalone `GitHubRepository` manifests inherit their owner from the root `github_owner` value because the Terraform GitHub provider selects one owner per provider configuration. `metadata.name` becomes the Terraform module key for either manifest kind and also defaults `spec.repository.name` for standalone repos when omitted.

Under the hood, Anvil translates standalone `GitHubRepository` manifests into the same repository settings object used by `GitHubTerraformRepository.spec.repository`, so both Terraform module paths share one repository-creation interface and one set of defaults.

Use `HCPTerraformWorkspace` when the GitHub repository already exists and you only want to add the Terraform workspace, GitHub environment wiring, and AWS provisioner role for one environment:

```yaml
apiVersion: anvil.emkaytec.dev/v1alpha1
kind: HCPTerraformWorkspace
metadata:
  name: sample-service-dev
spec:
  githubRepository: emkaytec/sample-service
  environment: dev
  aws:
    accountId: "111111111111"
  workspace:
    terraformVersion: "1.10.0"
```

`spec.githubRepository` must use the same owner configured by the root `github_owner` value because the Terraform GitHub provider selects one owner per provider configuration.

## Running Terraform

Run Terraform from the repo root once `.forge/` has been populated from the private manifests repository:

```bash
terraform init
terraform plan
```

Provider ownership is configured through the root module's explicit `default` provider aliases. `github_owner` is always required. `tfe_organization` and the shared StackSet role wiring are needed when you are planning `GitHubTerraformRepository` or `HCPTerraformWorkspace` manifests, and they belong in the ignored root `terraform.tfvars` file:

```hcl
github_owner     = "emkaytec"

tfe_organization = "emkaytec"
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
