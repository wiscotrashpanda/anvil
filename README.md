# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: GitHub repositories, HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` is intentionally left unchanged until its next role is clearer.

## Current Layout

- `modules/github-tf-repo/` defines the first extractable Terraform module for one repo-backed Terraform workload.
- `modules/github-tf-repo/examples/basic/` shows a minimal caller shape.

The module creates:

- one GitHub repository
- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- two AWS IAM provisioner roles per environment: GitHub Actions and HCP Terraform

## Direction

For now, this repository is the design and implementation space for the baseline architecture. Once the module contract settles, `modules/github-tf-repo` can be extracted into a standalone Terraform module repository, likely named `terraform-aws-github-tf-repo`.

Keep public code, module contracts, and sanitized examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.
