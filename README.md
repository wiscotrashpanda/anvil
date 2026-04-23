# Anvil

Anvil is being repurposed as the public-safe home for Emkaytec baseline architecture: GitHub repositories, HCP Terraform workspaces, and shared AWS provisioning roles.

The earlier manifest authoring, schema, and reconciliation work remains in `emkaytec/forge`. That path works, and it is useful context, but the active direction for baseline cloud setup is now Terraform-first instead of manifest/reconcile-first. `emkaytec/alloy` is intentionally left unchanged until its next role is clearer.

## Current Layout

- `manifests/` is the root desired-state input directory. Terraform reads `.yaml` and `.yml` files from there.
- `modules/github-tf-repo/` defines the first extractable Terraform module for one repo-backed Terraform workload.
- `modules/github-tf-repo/examples/basic/` shows a minimal caller shape.

The module creates:

- one GitHub repository
- one HCP Terraform workspace per environment
- one CloudFormation StackSet per environment
- two AWS IAM provisioner roles per environment: GitHub Actions and HCP Terraform

## Direction

For now, this repository is the design and implementation space for the baseline architecture. Once the module contract settles, `modules/github-tf-repo` can be extracted into a standalone Terraform module repository, likely named `terraform-aws-github-tf-repo`.

Keep public code, module contracts, and sanitized examples here. Real account IDs, operational manifests, tokens, and environment-specific values belong in private configuration.

## Manifests

Create one private YAML file per repo-backed Terraform workload under `manifests/`. Files ending in `.yaml` and `.yml` are ignored by Git so real desired state does not get committed to the public repository.

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
  environments:
    dev:
      account_id: "111111111111"
```

Run Terraform from the repo root after placing private manifests in `manifests/`:

```bash
terraform init
terraform plan
```

Set shared StackSet role wiring once in an ignored root `terraform.tfvars` file:

```hcl
stack_set_administration_role_arn = "arn:aws:iam::999999999999:role/AWSCloudFormationStackSetAdministrationRole"
stack_set_execution_role_name     = "AWSCloudFormationStackSetExecutionRole"
```
