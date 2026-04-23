# Manifests

Terraform reads desired state from YAML files in this directory.

Files ending in `.yaml` or `.yml` are intentionally ignored by Git because this repository is public-safe. Keep real account IDs, role ARNs, workspace settings, and environment-specific values local or in a private configuration repository.

Use one file per desired unit.

For a standalone non-Terraform GitHub repository, author a reduced `GitHubRepository` manifest:

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

For a repo-backed Terraform workload, use `GitHubTerraformRepository`:

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
      account_id: "111111111111"
    prod:
      account_id: "222222222222"
      region: us-east-2
```

Standalone `GitHubRepository` manifests inherit their owner from the root `github_owner` value because the Terraform GitHub provider selects one owner per provider configuration. `metadata.name` becomes the Terraform module key for either manifest kind and also defaults `spec.repository.name` for standalone repos when omitted.

Provider ownership is configured through the root module's explicit `emkaytec` provider aliases. `github_owner` is always required. `tfe_organization` and the shared StackSet role wiring are needed only when you are planning `GitHubTerraformRepository` manifests, and they belong in the ignored root `terraform.tfvars` file:

```hcl
github_owner     = "emkaytec"

tfe_organization = "emkaytec"
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
