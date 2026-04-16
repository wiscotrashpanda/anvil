# 0003: Publish release binaries only for now

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is primarily intended to run as a CLI inside GitHub Actions and similar automation environments. That makes versioned release binaries a natural distribution mechanism, especially for downstream repositories that want to pin a specific version and execute the tool directly.

The project briefly considered publishing Docker images alongside the binaries, but the current operating model does not require a container runtime. Adding container packaging later would be straightforward if downstream consumers ever need it.

The project therefore needed to decide whether to keep distribution focused on binaries now or invest in publishing multiple artifact types immediately.

## Decision

Anvil will publish versioned release binaries through GitHub Releases.

Docker images are explicitly deferred for now and can be added later if a real distribution or runtime need emerges.

## Rationale

- Release binaries are the cleanest fit for the current private-repo consumption model, where a workflow can download a pinned version and run it directly.
- The current GitHub Actions-based runtime does not require container packaging to be useful or usable.
- Keeping the release process binary-only reduces CI complexity, credentials management, and packaging surface area while the product is still taking shape.
- Docker images are easy enough to add later if containerized execution becomes valuable, so deferring that work does not meaningfully close off future options.

## Consequences

### Positive

- Downstream workflows can choose the simplest integration path for their environment.
- GitHub Releases become the canonical home for downloadable binaries.
- The release pipeline stays focused on the artifact type that matches the current execution model.
- The project still keeps optionality for future packaging and runtime decisions because Docker can be introduced later without changing the core CLI architecture.

### Negative

- Users who prefer containerized execution will not have an official image available yet.
- If Docker packaging becomes important later, the repository will need a follow-up change to add and document it.

## Alternatives Considered

### Publish only GitHub Release binaries

This is the chosen approach. It matches the current runtime model, keeps the release process simple, and preserves the option to add Docker later if a concrete need appears.

### Publish only Docker images

This would work for container-oriented usage, but it would make simple CLI consumption in GitHub Actions more awkward than downloading and executing a release binary directly.
