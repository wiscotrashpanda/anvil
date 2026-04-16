# 0004: Separate `anvil`, `smyth`, and shared schema code into distinct repositories

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is currently being built as a manifest-driven reconciliation CLI. Its job is to read YAML manifests, dispatch by `kind`, and converge external systems in a GitHub Actions workflow.

At the same time, there is a likely need for an additional operator-facing tool that helps create and update those manifests in the private implementation repository. That authoring workflow may include higher-level UX such as prompts, defaults, naming conventions, repository-specific layout decisions, and eventually pushing generated manifests to the implementation repository.

Those responsibilities are adjacent to reconciliation, but they are not the same concern. Reconciliation needs to stay predictable, explicit, and easy to debug in CI. Authoring and admin workflows are more likely to evolve toward interactive commands, opinionated generation, and repository mutation.

The project therefore needs to decide whether to keep both responsibilities inside one repository or split them into separate repositories while still avoiding duplicated schema and validation logic.

## Decision

The project will separate reconciliation, authoring, and shared schema code into distinct repositories from the start.

- `anvil` remains the reconciliation CLI.
- A separate admin or authoring repository, `smyth`, may handle manifest generation, editing, and implementation-repository-oriented workflows.
- A separate shared Go module repository, `alloy`, holds common manifest schema code consumed by both `anvil` and `smyth`.
- `anvil` and `smyth` should both depend on `alloy` for shared manifest types and validation rather than copying structs between repositories.

The initial extraction scope for `alloy` is intentionally small and schema-focused:

- manifest envelope types such as `apiVersion`, `kind`, and `metadata`
- versioned manifest structs such as `GitHubRepositorySpec` and `GitHubRepositoryManifest`
- schema constants such as supported `apiVersion` and kind names
- basic manifest-shape validation that is purely schema-oriented

The following should stay out of `alloy`:

- reconciliation planning or provider-specific apply logic
- CLI command handling
- filesystem walking and repo-specific loading behavior
- runtime concerns that belong specifically to `anvil` or `smyth`

## Rationale

- Reconciliation and authoring have different operational personalities. `anvil` should stay boring and deterministic in CI, while the authoring tool can be more interactive and operator-oriented.
- Keeping `anvil` focused on reading manifests and reconciling external state aligns with the product direction that execution should stay simple and explicit.
- Starting with separate repositories avoids the future migration cost and churn of carving a mixed codebase apart after both tools have already grown.
- A separate authoring repository provides room for higher-level UX without forcing those concerns into the runtime tool used in GitHub Actions.
- Shared manifest structs and validation help ensure that authoring and reconciliation agree on manifest shape instead of drifting into duplicated schema definitions.
- A dedicated shared module gives the common schema a stable home and versioned dependency path instead of treating it as a build byproduct.
- This approach supports a future where one higher-level repo-centric command may generate several atomic Anvil manifests without forcing humans to hand-author many low-level files.

## Consequences

### Positive

- `anvil` can remain narrowly focused on manifest consumption and reconciliation behavior.
- The authoring repository can optimize for operator experience without complicating the reconcile path.
- `alloy` provides a single source of truth for struct shape and validation logic.
- Starting with separate repositories avoids a later extraction/refactor project when the tools have already diverged operationally.
- The design leaves room for repo-centric generation workflows while preserving Anvil's atomic-manifest model internally.
- A versioned Go module gives both tools a cleaner and more durable dependency contract than CI artifacts or copied files.

### Negative

- The project will need to manage coordination across three repositories instead of one.
- Shared packages must be designed carefully so they remain useful common code rather than becoming a vague abstraction layer.
- Versioning and release discipline for `alloy` will matter, because both `anvil` and `smyth` will depend on it.
- There is some risk that `alloy` grows beyond schema concerns unless its scope is protected deliberately.

## Alternatives Considered

### One repository with both reconcile and authoring CLIs

This would reduce the number of repositories and may feel simpler at first. It was not chosen because it would still defer the eventual separation work until later, after both tools had already accumulated assumptions about living together. That increases the risk of a more painful extraction when the operational boundary is already clear today.

### Separate repositories with shared schema code still housed inside `anvil`

This would still leave `anvil` as the effective owner of code that is meant to be shared equally by both tools. It was not chosen because it keeps the dependency contract physically embedded inside one application repository and makes the eventual extraction to a neutral shared module more awkward.

### Separate repositories with CI artifacts as the schema handoff

This would work as a temporary bridge, but it was not chosen as the intended design because workflow artifacts are run-scoped and awkward as a long-term cross-repository dependency mechanism. A dedicated shared Go module provides a clearer and more stable integration point.
