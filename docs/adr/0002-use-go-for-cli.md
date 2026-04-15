# 0002: Use Go for the Anvil CLI

- Status: Accepted
- Date: 2026-04-15

## Context

The primary implementation language options considered for Anvil were Python and Go.

Python is the author's primary development language, with JavaScript and TypeScript also being familiar options. Python would have been a reasonable choice for CLI development and would likely have allowed fast iteration early on.

At the same time, Anvil is intended to ship as a versioned CLI artifact that runs predictably in GitHub Actions and other automation contexts. The project direction also has some conceptual overlap with Kubernetes-style manifest reconciliation, and future adjacent work may include Go-based tooling such as Terraform providers.

## Decision

Anvil will be implemented in Go.

## Rationale

- Go produces a simple executable binary, which is a strong fit for CI-driven installation and execution.
- A compiled binary simplifies release, distribution, and version pinning compared with a Python-based runtime and dependency environment.
- The manifest-and-reconciler model has a natural affinity with the ecosystem around Kubernetes and related infrastructure tools, where Go is a common implementation language.
- Using Go keeps the door open for future adjacent tooling, including potential provider-style integrations, without forcing a later language migration.
- The operational simplicity of a single deployable artifact matters more here than optimizing for the author's default scripting language.
- Using Go for this project is also a deliberate way to build deeper familiarity with a language that is widely used in infrastructure, platform, and cloud tooling.

## Consequences

### Positive

- Installation in GitHub Actions can be straightforward and reproducible.
- Versioned binary releases fit the public-repo/private-manifests model cleanly.
- Strong typing and a compiled distribution model should help keep the CLI predictable as it grows.
- The implementation language aligns well with common infrastructure and control-plane-adjacent tooling.

### Negative

- The implementation language is not the author's primary day-to-day language, so development may be slower at first.
- Some early iteration tasks may be less fluid than they would be in Python.
- The project takes on the learning and maintenance cost of building in Go from the start.

## Alternatives Considered

### Python

Python would have been a defensible choice and likely the fastest path to an initial prototype. It was not selected because distribution as a standalone binary artifact and long-term fit with infrastructure tooling mattered more than short-term familiarity.
