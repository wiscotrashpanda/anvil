# Anvil

Anvil is a minimal reconciliation-based infrastructure tool.

It is a Go CLI that reads YAML manifests and idempotently converges external systems when a GitHub Actions workflow runs. Anvil uses explicit, resource-specific reconcilers instead of a generic control-plane framework.

This repository is the public product repository for Anvil. It is intended to show the core design, implementation, and example manifest shape without exposing real production data.

## Status

This is an active working repository. The README is intentionally lightweight and will evolve as the CLI, reconciler implementations, packaging, and workflows take shape.

The repository now includes the first Go CLI scaffold under `cmd/anvil` and `internal/cli`.

## Core Principles

- Manifests are atomic: each manifest describes one resource kind.
- Execution stays simple: load YAML, dispatch by `kind`, and run the matching reconciler.
- Reconciliation must be idempotent: repeated runs should produce no changes when the external system already matches the manifest.
- Status is ephemeral in v1: results live in process output and logs, not in persisted state or manifest writeback.
- GitHub Actions is the runtime for v1: Anvil is not a background service.
- Distribution stays simple: Anvil is shipped as a versioned CLI binary built from the Go codebase for local use and workflow execution.

## V1 Scope

The v1 direction is a manifest-driven reconciliation CLI with explicit reconcilers for supported resource kinds.

Current public example support:

- `GitHubRepository`

Possible future kinds may include `AwsProvisionerRole` and `TerraformWorkspace`, but they are not part of the foundational v1 definition.

## Non-Goals

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

## Public Repo Boundary

This repository remains public by design.

- Public example manifests live under `examples/manifests/`.
- Public examples are illustrative only and are not deployable desired state.
- Public examples may only cover kinds that Anvil currently supports.
- Public examples must use sanitized placeholder values such as `example-org`, `example-repo`, and `123456789012`.
- Public examples must never include real organization names, repository names, account IDs, credentials, or operational values.
- Real environment-specific manifests belong in separate implementation repositories.

## Example Manifests

Example manifests live in [examples/manifests](/Volumes/Bolt/Code/wiscotrashpanda/anvil/examples/manifests/README.md).

## Local Development

Run the current CLI scaffold locally with:

```bash
go run ./cmd/anvil --help
```

## Architecture Decisions

Strategic and architectural decisions are tracked as ADRs under [docs/adr](/Volumes/Bolt/Code/wiscotrashpanda/anvil/docs/adr/README.md).

## AI-Assisted Development

AI agents are used in this repository for coding assistance, drafting, and documentation generation.

They are used as tools to accelerate implementation and communication, not as a substitute for engineering judgment or as unreviewed "vibe coding."

All code and documentation committed to this repository are reviewed by the repository author and are expected to be understood before they are kept.

## Working Notes

- `AGENTS.md` is the durable internal guidance file for the repository.
- `README.md` is the public-facing working document and should stay concise.
- As implementation lands, this file should be expanded with setup, usage, release, and workflow documentation.
