# ADR 0007: Use private YAML manifests for root composition

## Status

Accepted

## Context

Anvil now holds public-safe baseline architecture code. The module code can be public, but the desired state that composes real GitHub repositories, HCP Terraform workspaces, AWS accounts, StackSet roles, and environment settings may contain operational identifiers.

The repository needs a simple authoring boundary that lets the public code stay useful while keeping real environment values out of Git history.

## Decision

The root Terraform configuration reads desired state from YAML files in a root `manifests/` directory. Each YAML file describes one repo-backed Terraform workload with `kind: GitHubTerraformRepository`.

The `manifests/` directory is committed with documentation and placeholders, but real `.yaml` and `.yml` files are ignored by Git. Public examples should stay in documentation or use non-loaded extensions such as `.yaml.example` if needed.

## Consequences

- Root Terraform stays small and manifest-driven.
- The public repo can show the composition contract without committing real account data.
- Private desired-state files can be copied into `manifests/` for local or automation runs.
- The root composition is intentionally a Terraform authoring layer; Forge remains the place where the older manifest/reconcile implementation lives.
