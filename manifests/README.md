# Manifests

Terraform reads desired state from YAML files in this directory.

Files ending in `.yaml` or `.yml` are intentionally ignored by Git because this repository is public-safe. Keep real account IDs, role ARNs, workspace settings, and environment-specific values local or in a private configuration repository.

Use one file per repo-backed Terraform workload:

```yaml
apiVersion: anvil.emkaytec.dev/v1alpha1
kind: GitHubTerraformRepository
metadata:
  name: sample-service
spec:
  github_owner: emkaytec
  tfe_organization: emkaytec

  repository:
    description: Terraform-managed sample service.
    visibility: private
    topics:
      - aws
      - terraform

  environments:
    dev:
      account_id: "111111111111"
    prod:
      account_id: "222222222222"
      region: us-east-2
```

The `metadata.name` value becomes the module key and defaults the GitHub repository name unless `spec.repository.name` is provided.

Shared StackSet role wiring is system-wide and belongs in the ignored root `terraform.tfvars` file:

```hcl
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
