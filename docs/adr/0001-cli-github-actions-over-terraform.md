# 0001: Use a GitOps-style CLI with GitHub Actions instead of Terraform as the core foundation

- Status: Accepted
- Date: 2026-04-15

## Context

Anvil is intended to act as a reconciliable register of infrastructure resources expressed as individual YAML manifests. Early exploration considered Terraform as the primary implementation model.

Terraform would have been viable for provisioning and managing infrastructure through stacks, modules, and consolidated deployments. That model works well when the main problem is composing infrastructure into larger deployment units and applying those units together.

That was not the direction ultimately chosen for Anvil. The core problem here is closer to GitOps-style reconciliation of discrete resources from a manifest set than to plan-and-apply orchestration of grouped infrastructure changes.

## Decision

Anvil will be built as a CLI that runs in GitHub Actions and reconciles manifest-defined resources directly, instead of using Terraform as the foundational control model.

Terraform is not being rejected as a useful tool. It remains a likely part of the broader ecosystem around Anvil, and future Terraform repositories may be built on top of the foundation Anvil provides. The decision is specifically that Terraform should not be the primary abstraction for Anvil itself.

## Rationale

- The desired operating model is a reconciliable register of resources, not a stack-based deployment unit.
- Individual manifests map more naturally to explicit resource reconcilers than to Terraform modules or larger stateful stacks.
- A small CLI executed by GitHub Actions keeps the runtime model simple, visible, and easy to reason about.
- The product benefits from direct control over reconciliation behavior rather than inheriting Terraform's planning, state, and module-oriented workflow as the core mental model.
- This approach aligns better with standard GitOps infrastructure-as-code patterns for declarative desired state stored in Git.

## Consequences

### Positive

- The product can stay focused on manifest-driven reconciliation of discrete resources.
- Repository structure and execution flow remain simple: manifests in Git, reconciliation in CI, explicit logs as output.
- Anvil can define resource-specific behavior directly rather than forcing everything through Terraform state and module boundaries.
- Terraform can still be layered on top later where it is the right abstraction.

### Negative

- Anvil must implement and own its own reconciliation logic rather than delegating more of that behavior to Terraform.
- Some infrastructure patterns that are convenient in Terraform modules or stacks will need different treatment in Anvil.
- Operators familiar with Terraform may need to adjust to a different workflow and mental model.

## Alternatives Considered

### Terraform as the foundation

This would have worked, especially for grouped infrastructure deployments, but it pulled the design toward stacks, modules, and state-centric workflows. That was a mismatch for the narrower goal of a manifest-driven reconciliable resource register.
