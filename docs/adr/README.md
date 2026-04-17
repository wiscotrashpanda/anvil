# Architecture Decision Records

This directory contains Architecture Decision Records for Anvil.

ADRs capture durable engineering and product decisions, especially when there were real alternatives or trade-offs. They are intended to explain why a decision was made, not just what the repository currently contains.

These ADRs were drafted with help from AI coding agents after extended discussions about scope, trade-offs, and underlying concerns. The decisions themselves remain the repository author's; AI assistance is used to help capture the reasoning clearly and completely in writing.

## Index

- [0001: Use a GitOps-style CLI with GitHub Actions instead of Terraform as the core foundation](/Volumes/Bolt/Code/emkaytec/anvil/docs/adr/0001-cli-github-actions-over-terraform.md)
- [0002: Use Go for the Anvil CLI](/Volumes/Bolt/Code/emkaytec/anvil/docs/adr/0002-use-go-for-cli.md)
- [0003: Publish release binaries and GHCR images](/Volumes/Bolt/Code/emkaytec/anvil/docs/adr/0003-publish-release-binaries-and-ghcr-images.md)
- [0004: Separate `anvil`, `smyth`, and shared schema code into distinct repositories](/Volumes/Bolt/Code/emkaytec/anvil/docs/adr/0004-separate-authoring-and-reconciliation-clis.md)
- [0005: Manage the core GitHub repository surface first and expand intentionally](/Volumes/Bolt/Code/emkaytec/anvil/docs/adr/0005-manage-core-github-repository-surface-first.md)
