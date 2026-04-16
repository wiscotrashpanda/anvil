# 0003: Publish release binaries and GHCR images

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is primarily intended to run as a CLI inside GitHub Actions and similar automation environments. That makes versioned release binaries a natural distribution mechanism, especially for downstream repositories that want to pin a specific version and execute the tool directly.

At the same time, containerized execution remains useful for some automation environments and local testing. A Docker image provides a portable runtime and a consistent execution environment without changing the core CLI model.

The project therefore needed to decide whether to keep distribution binary-only or support both binaries and container images. If images are published, the repository also needs a registry choice.

## Decision

Anvil will publish both:

- versioned release binaries through GitHub Releases
- container images through GitHub Container Registry (`ghcr.io`)

The binary release path is the primary distribution path for CI/CD consumption. The Docker image remains a supported secondary distribution option for portability and container-based execution.

## Rationale

- Release binaries are the cleanest fit for the current private-repo consumption model, where a workflow can download a pinned version and run it directly.
- Docker images remain useful for container-based execution, local testing, and future orchestration patterns.
- GitHub Container Registry keeps the container artifact colocated with the source repository, release automation, and GitHub-native permissions model.
- Using GHCR avoids the extra account, secret, and external registry management overhead that came with Docker Hub.
- Maintaining both artifacts keeps the project flexible while staying close to the existing GitHub-centered operating model.

## Consequences

### Positive

- Downstream workflows can choose the simplest integration path for their environment.
- GitHub Releases become the canonical home for downloadable binaries.
- GitHub Container Registry provides an official image for users who prefer containerized execution.
- The project keeps optionality for future packaging and runtime decisions without introducing a separate external registry dependency.

### Negative

- The CI pipeline has to maintain and verify two distribution paths instead of one.
- Documentation must stay clear about which artifact should be considered the default path.
- Release management becomes slightly more complex because multiple artifact types are published together.

## Alternatives Considered

### Publish only GitHub Release binaries

This is the simplest primary distribution model and likely the most common path for downstream usage. It was not chosen as the exclusive model because keeping the Docker image available provides useful flexibility for container-oriented consumers.

### Publish only Docker images

This would work for container-oriented usage, but it would make simple CLI consumption in GitHub Actions more awkward than downloading and executing a release binary directly.

### Publish images through Docker Hub

This would also work technically, but GHCR is a better fit for the current repository because it keeps authentication, permissions, and artifact ownership inside GitHub rather than depending on a separate registry account and secret set.
