package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	if err := Run([]string{"help"}, &stdout); err != nil {
		t.Fatalf("Run(help) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Hello from Anvil.") {
		t.Fatalf("expected help output to contain greeting, got %q", output)
	}
}

func TestRunNoArgsShowsHelp(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	if err := Run(nil, &stdout); err != nil {
		t.Fatalf("Run(nil) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Fatalf("expected help output to contain usage, got %q", output)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	err := Run([]string{"nope"}, &stdout)
	if err == nil {
		t.Fatal("expected unknown command error")
	}

	if !strings.Contains(err.Error(), "unknown command: nope") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func TestRunReconcileDefaultsToCurrentDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeCLIFile(t, filepath.Join(root, "nested", "repo.yaml"), `apiVersion: anvil.example.io/v1alpha1
kind: GitHubRepository
metadata:
  name: example-repo
spec:
  owner: example-org
  name: example-repo
`)

	var stdout bytes.Buffer

	err := runReconcileWithWorkingDir([]string{}, &stdout, func() (string, error) {
		return root, nil
	})
	if err != nil {
		t.Fatalf("runReconcileWithWorkingDir returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Reconciling GitHubRepository example-org/example-repo") {
		t.Fatalf("expected reconcile output to contain repository message, got %q", output)
	}
}

func TestRunReconcilePrintsDryRunMessages(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeCLIFile(t, filepath.Join(root, "nested", "repo.yaml"), `apiVersion: anvil.example.io/v1alpha1
kind: GitHubRepository
metadata:
  name: example-repo
spec:
  owner: example-org
  name: example-repo
`)

	var stdout bytes.Buffer

	err := Run([]string{"reconcile", "--manifests", root}, &stdout)
	if err != nil {
		t.Fatalf("Run(reconcile) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Reconciling GitHubRepository example-org/example-repo") {
		t.Fatalf("expected reconcile output to contain repository message, got %q", output)
	}

	if !strings.Contains(output, "Dry run only: no external changes applied") {
		t.Fatalf("expected reconcile output to contain dry-run message, got %q", output)
	}
}

func TestRunReconcileHelpDescribesDefaultDirectory(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	err := Run([]string{"reconcile", "--help"}, &stdout)
	if err != nil {
		t.Fatalf("Run(reconcile --help) returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "anvil reconcile [--manifests <path>]") {
		t.Fatalf("expected reconcile help usage to show optional manifests flag, got %q", output)
	}

	if !strings.Contains(output, "defaults to current directory") {
		t.Fatalf("expected reconcile help to mention current directory default, got %q", output)
	}
}

func writeCLIFile(t *testing.T, path string, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", path, err)
	}
}
