# 0003: Publish both release binaries and Docker images

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is primarily intended to run as a CLI inside GitHub Actions and similar automation environments. That makes versioned release binaries a natural distribution mechanism, especially for downstream repositories that want to pin a specific version and execute the tool directly.

At the same time, containerized execution remains useful. A Docker image provides a portable runtime, a consistent execution environment, and a path for future use cases where running Anvil as a container is more convenient than downloading a standalone binary.

The project therefore needed to decide whether to standardize on one distribution format or keep both available.

## Decision

Anvil will publish both:

- versioned release binaries through GitHub Releases
- container images through Docker Hub

The binary release path is the primary distribution path for CI/CD consumption. The Docker image remains a supported distribution option for portability, experimentation, and future flexibility.

## Rationale

- Release binaries are the cleanest fit for the current private-repo consumption model, where a workflow can download a pinned version and run it directly.
- Docker images remain useful for container-based execution, local testing, and future orchestration patterns.
- Maintaining both artifacts keeps the project flexible while the final operating model continues to evolve.
- The incremental cost of publishing both is low once the build pipeline is already producing the binary and the container image.
- Keeping both outputs is also a practical way to explore and demonstrate multiple distribution paths for infrastructure tooling, and, frankly, because we can without adding much complexity.

## Consequences

### Positive

- Downstream workflows can choose the simplest integration path for their environment.
- GitHub Releases become the canonical home for downloadable binaries.
- Docker Hub remains available for users who prefer containerized execution.
- The project keeps optionality for future packaging and runtime decisions.

### Negative

- The CI pipeline has to maintain and verify two distribution paths instead of one.
- Documentation must stay clear about which artifact should be considered the default path.
- Release management becomes slightly more complex because multiple artifact types are published together.

## Alternatives Considered

### Publish only GitHub Release binaries

This is the simplest primary distribution model and likely the most common path for downstream usage. It was not chosen as the exclusive model because keeping the Docker image available provides useful flexibility at relatively low additional cost.

### Publish only Docker images

This would work for container-oriented usage, but it would make simple CLI consumption in GitHub Actions more awkward than downloading and executing a release binary directly.
