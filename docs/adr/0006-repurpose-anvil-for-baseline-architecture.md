# ADR 0006: Repurpose Anvil for baseline architecture

## Status

Accepted

## Context

Anvil was previously the likely home for carving the Forge reconciler into its own command line utility backed by shared schema from Alloy. Since then, Forge has grown working manifest generation, schema, and reconciliation packages. That implementation does what it set out to do, but it is no longer the preferred active path for baseline cloud setup.

A Terraform module can express the same "one declaration fans out into a repo stack" idea more directly for GitHub, HCP Terraform, and AWS shared resources. It also avoids carrying a bespoke reconcile loop for infrastructure Terraform already manages well.

## Decision

Repurpose `emkaytec/anvil` as the public-safe baseline architecture repository.

The first module lives in `modules/github-tf-repo` and creates a GitHub repository, HCP Terraform workspaces, and per-environment AWS IAM provisioner roles through CloudFormation StackSets. It stays inside Anvil while the contract settles, then can be extracted into a standalone Terraform module repository.

Forge keeps the existing manifest and reconciliation implementation as working context. It may be marked deprecated or legacy later, but it should not be deleted just because this direction is taking over. Alloy is left unchanged until there is a clearer shared-schema need.

## Consequences

- Anvil's active purpose changes from reconciler CLI extraction to baseline architecture definition.
- The module is isolated under `modules/` so extraction later is mechanical.
- The old Forge manifest/reconcile path remains available for reference and comparison.
- Public examples must stay sanitized; real operational values remain outside this repository.
