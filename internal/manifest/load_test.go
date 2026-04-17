package manifest

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDirRecursivelyLoadsGitHubRepositoryManifests(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "nested", "repo.yaml"), `apiVersion: anvil.example.io/v1alpha1
kind: GitHubRepository
metadata:
  name: example-repo
spec:
  owner: example-org
  name: example-repo
`)
	writeTestFile(t, filepath.Join(root, "ignore.txt"), "not yaml")

	result, err := LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir returned error: %v", err)
	}

	if len(result.GitHubRepositories) != 1 {
		t.Fatalf("expected 1 GitHubRepository, got %d", len(result.GitHubRepositories))
	}

	repository := result.GitHubRepositories[0]
	if repository.Manifest.Metadata.Name != "example-repo" {
		t.Fatalf("expected metadata.name to be example-repo, got %q", repository.Manifest.Metadata.Name)
	}

	if repository.Manifest.Spec.Owner != "example-org" {
		t.Fatalf("expected spec.owner to be example-org, got %q", repository.Manifest.Spec.Owner)
	}
}

func TestLoadDirRecursivelyLoadsHCPTerraformWorkspaceManifests(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "nested", "workspace.yaml"), `apiVersion: anvil.example.io/v1alpha1
kind: HCPTerraformWorkspace
metadata:
  name: example-workspace
spec:
  organization: example-org
  name: example-workspace
`)

	result, err := LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir returned error: %v", err)
	}

	if len(result.HCPTerraformWorkspaces) != 1 {
		t.Fatalf("expected 1 HCPTerraformWorkspace, got %d", len(result.HCPTerraformWorkspaces))
	}

	workspace := result.HCPTerraformWorkspaces[0]
	if workspace.Manifest.Metadata.Name != "example-workspace" {
		t.Fatalf("expected metadata.name to be example-workspace, got %q", workspace.Manifest.Metadata.Name)
	}

	if workspace.Manifest.Spec.Organization != "example-org" {
		t.Fatalf("expected spec.organization to be example-org, got %q", workspace.Manifest.Spec.Organization)
	}
}

func TestLoadDirRejectsUnsupportedKind(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "workspace.yaml"), `apiVersion: anvil.example.io/v1alpha1
kind: TerraformWorkspace
metadata:
  name: example-workspace
spec: {}
`)

	_, err := LoadDir(root)
	if err == nil {
		t.Fatal("expected unsupported kind error")
	}

	if !strings.Contains(err.Error(), `unsupported kind "TerraformWorkspace"`) {
		t.Fatalf("expected unsupported kind error, got %v", err)
	}
}

func TestLoadDirRequiresManifestFields(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "repo.yaml"), `kind: GitHubRepository
metadata:
  name: example-repo
spec:
  owner: example-org
  name: example-repo
`)

	_, err := LoadDir(root)
	if err == nil {
		t.Fatal("expected missing apiVersion error")
	}

	if !strings.Contains(err.Error(), "missing apiVersion") {
		t.Fatalf("expected missing apiVersion error, got %v", err)
	}
}
