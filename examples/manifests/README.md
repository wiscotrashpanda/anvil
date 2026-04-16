# Example Manifests

This directory contains public-safe example manifests for Anvil.

These files are examples only:

- They illustrate manifest shape and expected usage.
- They are not operationally authoritative desired state.
- They must only cover kinds Anvil currently supports.
- They must use sanitized placeholders such as `example-org`, `example-repo`, and `123456789012`.
- They must never include real organization names, repository names, account IDs, credentials, or environment-specific values.

Real manifests belong in separate implementation repositories.

## Local Testing

For local testing against a basic private GitHub repository, start with the smallest possible `GitHubRepository` manifest shape:

```yaml
apiVersion: anvil.example.io/v1alpha1
kind: GitHubRepository
metadata:
  name: example-repo
spec:
  owner: example-org
  name: example-repo
  visibility: private
  description: Example GitHub repository manifest for local testing
  autoInit: true
```

Add optional features only after the base reconcile works for the target repository.

Common GitHub prerequisites for optional features:

- `customProperties` requires the custom property to already exist in the owning organization's GitHub custom property schema.
- `branches[].protection` may be unavailable for private repositories on plans that do not include private branch protection.
- `securityAndAnalysis.advancedSecurity` should be omitted for public repositories because GitHub treats it as already available there.
- `pages` requires GitHub Pages support to be available for the repository and branch/source settings you declare.
