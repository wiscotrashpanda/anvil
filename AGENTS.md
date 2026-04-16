# Anvil

## Product Overview

Anvil is a minimal reconciliation-based infrastructure tool.

It is a Go CLI that reads YAML manifests and idempotently converges external systems when a GitHub Actions workflow runs. Anvil uses explicit, resource-specific reconcilers instead of a generic control-plane framework.

Anvil is not Kubernetes, Crossplane, or a long-running control plane.

## Core Principles

- Manifests are atomic: each manifest describes one resource kind.
- Execution stays simple: load YAML, dispatch by `kind`, and run the matching reconciler.
- Reconciliation must be idempotent: repeated runs should produce no changes when the external system already matches the manifest.
- Status is ephemeral in v1: results live in process output and logs, not in persisted state or manifest writeback.
- GitHub Actions is the runtime for v1: Anvil is not a background service.
- Distribution stays simple: Anvil is shipped as a versioned CLI binary built from the Go codebase for local use and workflow execution.
- The product optimizes for clarity, simplicity, and debuggability over abstraction, extensibility, or future-proofing.

## V1 Scope

The v1 product direction is a manifest-driven reconciliation CLI with explicit reconcilers for supported resource kinds, implemented in Go and distributed as a versioned CLI binary.

Current supported kind for public examples:

- `GitHubRepository`

Possible future kinds may include `AwsProvisionerRole` and `TerraformWorkspace`, but they are not part of the foundational v1 definition.

## Explicit Non-Goals

Anvil v1 does not include:

- Plugin systems
- Generic CRD or resource frameworks
- Persistent control planes
- Event buses
- Dependency graph engines
- Cross-resource orchestration
- State persistence
- Advanced drift engines
- Background services
- Deletion handling

## Public Repository Boundary

This repository is intended to remain public.

- Public example manifests live under `examples/manifests/`.
- Public examples are illustrative only and are not deployable desired state.
- Public examples may only cover kinds that Anvil currently supports.
- Public examples must use sanitized placeholder values such as `example-org`, `example-repo`, and `123456789012`.
- Public examples must never include real organization names, repository names, account IDs, credentials, or operational values.
- Real environment-specific manifests belong in separate implementation repositories.

## Shared Code Boundary

Anvil should keep reconciliation-specific runtime code in this repository and pull shared manifest/schema code from the separate `alloy` project.

- Common manifest structs, schema versions, kind constants, and schema-oriented validation should be added to `alloy`, not redefined locally in `anvil`.
- When Anvil needs new shared types or schema changes, update `alloy` first, then load the new version into this project through the Go module dependency.
- Avoid recreating a local shared `pkg` tree in this repository as a convenience shortcut; `alloy` is the source of truth for shared code.
- Keep `anvil` focused on manifest loading, reconciliation planning, and provider-specific runtime behavior rather than ownership of common schema packages.

## Working Style

- Keep durable project guidance in this file rather than introducing a spec system by default.
- Add focused documentation only when it materially helps contributors understand the product or implementation.
- Prefer direct implementation work over process-heavy planning artifacts.
