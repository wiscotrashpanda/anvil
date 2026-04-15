# 0004: Separate manifest authoring and reconciliation into distinct CLIs with shared code

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is currently being built as a manifest-driven reconciliation CLI. Its job is to read YAML manifests, dispatch by `kind`, and converge external systems in a GitHub Actions workflow.

At the same time, there is a likely need for an additional operator-facing tool that helps create and update those manifests in the private implementation repository. That authoring workflow may include higher-level UX such as prompts, defaults, naming conventions, repository-specific layout decisions, and eventually pushing generated manifests to the implementation repository.

Those responsibilities are adjacent to reconciliation, but they are not the same concern. Reconciliation needs to stay predictable, explicit, and easy to debug in CI. Authoring and admin workflows are more likely to evolve toward interactive commands, opinionated generation, and repository mutation.

The project therefore needs to decide whether to keep both responsibilities inside a single CLI surface or split them into separate executables while still avoiding duplicated schema and validation logic.

## Decision

The project will separate reconciliation and manifest authoring into distinct CLI tools.

- `anvil` remains the reconciliation CLI.
- A separate admin or authoring CLI, currently expected to be named `smith`, may handle manifest generation, editing, and implementation-repository-oriented workflows.
- The two tools should share common Go packages for manifest types, schema-aligned structs, validation, and related helper logic where that reduces duplication without blurring the executable boundary.

This is a tool-boundary decision, not necessarily a repository-boundary decision. The CLIs may still live in the same repository while the design is young and the shared packages are being established.

## Rationale

- Reconciliation and authoring have different operational personalities. `anvil` should stay boring and deterministic in CI, while the authoring tool can be more interactive and operator-oriented.
- Keeping `anvil` focused on reading manifests and reconciling external state aligns with the product direction that execution should stay simple and explicit.
- Separating the executables reduces the risk that reconciliation grows a large surface area of admin commands, prompts, and generation logic that are unrelated to its core runtime job.
- A separate authoring CLI provides room for higher-level UX without forcing those concerns into the runtime tool used in GitHub Actions.
- Shared manifest structs and validation help ensure that authoring and reconciliation agree on manifest shape instead of drifting into duplicated schema definitions.
- This approach supports a future where one higher-level repo-centric command may generate several atomic Anvil manifests without forcing humans to hand-author many low-level files.

## Consequences

### Positive

- `anvil` can remain narrowly focused on manifest consumption and reconciliation behavior.
- The authoring CLI can optimize for operator experience without complicating the reconcile path.
- Shared manifest packages provide a single source of truth for struct shape and validation logic.
- The design leaves room for repo-centric generation workflows while preserving Anvil's atomic-manifest model internally.
- Keeping the tools in one repository initially still allows fast iteration and straightforward code sharing.

### Negative

- The project will need to name, package, test, and document two executables instead of one.
- Shared packages must be designed carefully so they remain useful common code rather than becoming a vague abstraction layer.
- There is some risk of confusion for contributors if the responsibilities of `anvil` and the authoring CLI are not documented clearly.
- Release and distribution decisions may become slightly more complex if both tools are ultimately published.

## Alternatives Considered

### One CLI with both reconcile and authoring commands

This would reduce the number of executables and may feel simpler at first. It was not chosen because it would mix two different concerns into one command surface: a deterministic CI reconciler and an operator-facing authoring/admin tool. That would make it easier for the runtime CLI to drift into a broader and less focused product.

### Separate repositories for separate tools immediately

This would maximize isolation, but it introduces extra overhead too early. The tools are still young, and shared manifest structs and validation are likely to change together. Keeping them in one repository for now preserves fast iteration while still respecting the boundary at the executable level.
