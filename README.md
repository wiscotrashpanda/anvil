# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: standalone GitHub repositories, repo-backed HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` now hosts reusable Terraform modules published through the private HCP Terraform module registry.

## Current Layout

- `.forge/` is the root desired-state input directory. Terraform reads `.yaml` and `.yml` files from there. The directory is gitignored and intentionally never committed to this public repository; it is supplied at plan/apply time from a private configuration repository (see [Manifests](#manifests)).
- The root Terraform translates `.forge/` manifests into calls to the private HCP Terraform module registry module `app.terraform.io/emkaytec/repository/github`, pinned at version `0.3.0`.

The private registry repository module always creates:

- one GitHub repository

When Terraform workspaces are enabled, it also creates:

- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- one AWS IAM provisioner role per environment trusted by GitHub Actions and HCP Terraform
- one GitHub environment with AWS variables per environment

## Direction

For now, this repository is the public-safe composition layer for Emkaytec baseline architecture. The reusable Terraform module implementation and versioning live outside this repo in Emkaytec's private HCP Terraform module registry.

Keep public composition code and sanitized manifest examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.

## Manifests

Terraform reads desired state from YAML files in the root `.forge/` directory. The directory is gitignored and is not part of this repository; it is populated from a private manifests repository that owns all real account IDs, role ARNs, workspace settings, and other environment-specific values.

In GitHub Actions, check out the private manifests repo at the same level so its `.forge/` directory lands at the root of this working tree and Terraform can read it in place. Locally, populate `.forge/` the same way (for example, by cloning or symlinking the private repo so its `.forge/` is visible at the repo root).

The intended authoring path is the Emkaytec Forge CLI (`emkaytec/forge`), which can generate these manifests into `.forge/`. Adding or editing the YAML files manually also works.

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
    topics:
      - docs
      - website
    autoInit: true
    defaultBranch: main
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
  name: complete-service
spec:
  createTerraformWorkspaces: true
  repository:
    description: Complete example Terraform-backed GitHub repository.
    visibility: private
    homepage: https://example.com/platform/complete-service
    topics:
      - aws
      - terraform
      - platform
    autoInit: true
    defaultBranch: main
    features:
      hasIssues: true
      hasProjects: false
      hasWiki: false
      hasDiscussions: true
    mergePolicy:
      allowMergeCommit: true
      allowSquashMerge: true
      allowRebaseMerge: false
      deleteBranchOnMerge: true

  aws:
    region: us-east-1
    partition: aws
    managedPolicyArns:
      - arn:aws:iam::aws:policy/ReadOnlyAccess
    githubOidcProviderHost: token.actions.githubusercontent.com
    githubOidcAudience: sts.amazonaws.com
    tfeOidcProviderHost: app.terraform.io
    tfeOidcAudience: aws.workload.identity
    stackSetNamePrefix: complete-service
    stackSetPermissionModel: SELF_MANAGED
    stackSetCallAs: SELF
    stackSetOperationPreferences:
      failure_tolerance_count: 0
      max_concurrent_count: 1
      region_concurrency_type: SEQUENTIAL
      region_order:
        - us-east-2
    retainStackInstancesOnDestroy: false
    tags:
      managed-by: terraform
      owner: platform
      service: complete-service

  environments:
    admin:
      aws:
        accountId: "111111111111"
        region: us-east-2
        partition: aws
        managedPolicyArns:
          - arn:aws:iam::aws:policy/PowerUserAccess
        githubActionsSubject: repo:emkaytec/complete-service:environment:admin
      workspace:
        name: complete-service-admin
        projectId: prj-example
        projectName: platform
        executionMode: agent
        agentPoolId: apool-example
        terraformVersion: "1.14.8"
        autoApply: false
        queueAllRuns: true
        speculativeEnabled: true
        workingDirectory: terraform
        tags:
          environment: admin
          owner: platform
          service: complete-service
        vcsRepo:
          branch: main
          githubAppInstallationId: 12345678
          ingressSubmodules: false
          triggerPatterns:
            - terraform/**/*.tf
          triggerPrefixes:
            - terraform/
        variables:
          - key: CUSTOM_ADMIN_ENV
            value: admin-value
            type: env
            sensitive: false
            description: Example environment variable scoped to the admin workspace.
        manageVariables: true
        hcpTerraformSubject: organization:emkaytec:project:platform:workspace:complete-service-admin:run_phase:*

  workspace:
    projectName: platform
    terraformVersion: "1.14.8"
    executionMode: remote
    autoApply: false
    queueAllRuns: true
    speculativeEnabled: true
    workingDirectory: terraform
    tags:
      owner: platform
      service: complete-service
    vcsRepo:
      branch: main
      githubAppInstallationId: 12345678
      ingressSubmodules: false
      triggerPatterns:
        - terraform/**/*.tf
      triggerPrefixes:
        - terraform/
    variables:
      - key: custom_workspace_variable
        value: example-value
        type: terraform
        sensitive: true
        description: Example Terraform variable created in every managed workspace.
    manageVariables: true
```

`GitHubRepository` manifests inherit their owner from the root `github_owner` value because the Terraform GitHub provider selects one owner per provider configuration. `metadata.name` becomes the Terraform module key and also defaults `spec.repository.name` when omitted.

`workspace.variables` and `environments.*.workspace.variables` create custom HCP Terraform workspace variables. Each entry must include `key` and `value`; `type` defaults in the module to `terraform` and may be `terraform` or `env`, `sensitive` defaults to `false`, and `description` is optional. Environment-specific variables are appended to shared workspace variables, and each final workspace must have unique `key`/`type` pairs.

## Running Terraform

Run Terraform from the repo root once `.forge/` has been populated from the private manifests repository. `terraform init` downloads the pinned private module from HCP Terraform, so the environment needs credentials for `app.terraform.io` through Terraform CLI login or `TF_TOKEN_app_terraform_io`.

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
