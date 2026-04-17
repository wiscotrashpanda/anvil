# 0005: Manage the core GitHub repository surface first and expand intentionally

- Status: Accepted
- Date: 2026-04-16

## Context

Anvil is intended to be a minimal reconciliation-based infrastructure tool with explicit, resource-specific behavior. The initial public resource kind is `GitHubRepository`.

As implementation of the `GitHubRepository` reconciler has progressed, it has become clearer that the total GitHub repository management surface is wider and more irregular than a first pass might suggest. Some settings are simple repository attributes with stable update semantics. Others are conditional on account type, organization policy, product licensing, repository visibility, or separate GitHub feature enablement. Some behaviors are effectively adjacent resources that happen to be associated with a repository rather than plain repository settings.

That creates a product and engineering choice. One option would be to treat `GitHubRepository` as a broad parity target and continue expanding reconciliation coverage toward the full downstream provider surface. Another option would be to define a narrower managed surface for v1, focus on the common baseline attributes that are valuable and predictable, and add more coverage later only when the implementation and management contract are clear.

This question matters both for product quality and for project sustainability. Anvil is a working tool and a portfolio project, not a full clone of Terraform or the GitHub Terraform provider. The project benefits more from a clear and defensible management contract than from chasing every edge case early.

## Decision

Anvil will manage the core `GitHubRepository` surface first and expand intentionally over time instead of pursuing broad provider parity up front.

For v1 direction, `GitHubRepository` should prioritize attributes that are all or most of the following:

- commonly needed as part of an initial repository baseline
- safe to reconcile idempotently with clear drift semantics
- supported through stable repository-level API behavior
- understandable to operators without deep provider-specific caveats

The initial managed surface should favor baseline repository settings such as:

- repository existence and identity
- visibility where safe and intentionally declared
- description and homepage
- default branch
- basic repository feature toggles
- topics
- straightforward merge policy settings

More conditional or irregular behavior should be treated conservatively. Depending on the feature, Anvil may:

- defer support entirely until the management contract is clear
- treat the feature as best-effort with explicit operator-facing messaging
- eventually split the behavior into a separate resource kind rather than keeping it inside `GitHubRepository`

The project does not need to manage every GitHub repository-associated feature in v1. Omitted or explicitly unsupported fields are allowed to remain manual for now.

## Rationale

- A narrow managed surface is more consistent with Anvil's product goal of clarity and explicitness than attempting full provider parity.
- The common repository baseline is where the tool can provide the most value early with the least confusing behavior.
- Stable, idempotent reconciliation is more important than exposing a large manifest surface that only works inconsistently.
- GitHub includes many settings whose behavior depends on plan level, visibility, organization setup, or separate prerequisite resources. Those should not all be treated as routine repository drift correction.
- Deliberate scoping keeps the implementation tractable across `anvil`, `alloy`, and `smyth` instead of forcing all three repositories to evolve in lockstep around a rapidly expanding edge-case surface.
- A smaller, well-documented contract is a better demonstration of engineering judgment than an aspirational but incomplete imitation of a mature provider.
- Expanding later remains possible. The decision is to sequence the work intentionally, not to forbid richer reconciliation forever.

## Consequences

### Positive

- `GitHubRepository` can become reliable sooner for the common baseline use case.
- The schema in `alloy` can stay aligned with a clearly owned management contract rather than growing around speculative support.
- `smyth` can default to generating the most useful and least surprising manifest shape first.
- Documentation can describe supported behavior more confidently and with fewer caveats.
- Future resource additions can proceed without blocking on GitHub edge-case completeness.

### Negative

- Some repository settings will remain manual or partially managed for a time.
- Operators may need to use GitHub directly for advanced or uncommon configuration.
- There is a risk of temporary frustration if users expect the manifest to imply complete repository ownership.
- The project will need to document supported, best-effort, and unsupported behavior clearly to avoid ambiguity.

## Alternatives Considered

### Pursue broad GitHub repository parity immediately

This would aim to make `GitHubRepository` manage most or all repository-associated features from the start. It was not chosen because it would pull the project toward reproducing a mature provider surface before the core management contract is stable. That increases implementation cost, edge-case exposure, and cross-repository coordination burden without proportionate early product value.

### Treat Terraform as the immediate replacement for GitHub repository management

Terraform remains a credible tool for GitHub management and may still become part of the broader ecosystem around Anvil. It was not chosen as the immediate response to this scoping problem because the issue identified here is not that the CLI model is invalid; it is that the current managed surface should be narrower and more intentional. The first corrective move is to tighten scope, not to replace the core foundation.

### Keep every field in the schema but silently ignore unsupported behavior

This would reduce implementation pressure in the short term, but it was not chosen as the intended direction because silent non-ownership creates ambiguity. If Anvil does not own a field reliably, that should be reflected explicitly in schema shape, documentation, or operator messaging rather than hidden behind a broad manifest that suggests stronger control than the tool actually provides.
