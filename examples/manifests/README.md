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

Common GitHub notes for the currently supported example surface:

- `topics` are supported, but it is still easier to start with a minimal repository manifest and add them after the base create/update flow works.

For local testing against HCP Terraform, start with a small `HCPTerraformWorkspace` manifest and add tags, variables, remote-state consumers, and variable-set assignments after the base workspace reconcile works:

```yaml
apiVersion: anvil.example.io/v1alpha1
kind: HCPTerraformWorkspace
metadata:
  name: example-workspace
spec:
  organization: example-org
  name: example-workspace
  projectID: prj-123456
  description: Example HCP Terraform workspace manifest for local testing
  terraformVersion: 1.14.8
  executionMode: remote
```

The current HCP Terraform reconciler intentionally does not manage these fields yet:

- `sshKeyID`
- `runTriggers`
- `teamAccess`
- `notifications`
