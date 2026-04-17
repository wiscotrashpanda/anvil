package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	ghapi "github.com/emkaytec/anvil/internal/github"
	hcpapi "github.com/emkaytec/anvil/internal/hcpterraform"
	"github.com/emkaytec/anvil/internal/reconcile"
)

var executePlan = func(ctx context.Context, plan reconcile.Plan) ([]string, error) {
	var githubClient *ghapi.Client
	if plan.HasGitHubRepositories() {
		client, err := reconcile.NewGitHubClient()
		if err != nil {
			return nil, err
		}
		githubClient = client
	}

	var hcpTerraformClient *hcpapi.Client
	if plan.HasHCPTerraformWorkspaces() {
		client, err := reconcile.NewHCPTerraformClient()
		if err != nil {
			return nil, err
		}
		hcpTerraformClient = client
	}

	return plan.Run(ctx, githubClient, hcpTerraformClient)
}

const helpText = `Hello from Anvil.

This is the current Anvil CLI.

Usage:
  anvil [command]

Available Commands:
  help        Show this message
  reconcile   Reconcile manifests from a directory
`

const reconcileHelpText = `Reconcile manifests from a directory.

Usage:
  anvil reconcile [--manifests <path>]

Flags:
  --manifests string   Path to a directory containing YAML manifests (defaults to current directory)
`

// Run executes the CLI against the provided arguments.
func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprint(stdout, helpText)
		return err
	}

	switch args[0] {
	case "help", "--help", "-h":
		_, err := fmt.Fprint(stdout, helpText)
		return err
	case "reconcile":
		return runReconcile(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runReconcile(args []string, stdout io.Writer) error {
	return runReconcileWithWorkingDir(args, stdout, os.Getwd)
}

func runReconcileWithWorkingDir(args []string, stdout io.Writer, getwd func() (string, error)) error {
	for _, arg := range args {
		if arg == "help" || arg == "--help" || arg == "-h" {
			_, err := fmt.Fprint(stdout, reconcileHelpText)
			return err
		}
	}

	flagSet := flag.NewFlagSet("reconcile", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	manifestsPath := flagSet.String("manifests", "", "Path to a directory containing YAML manifests")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		return fmt.Errorf("unexpected arguments: %v", flagSet.Args())
	}

	if *manifestsPath == "" {
		cwd, err := getwd()
		if err != nil {
			return fmt.Errorf("resolve current working directory: %w", err)
		}

		*manifestsPath = cwd
	}

	plan, err := reconcile.Load(*manifestsPath)
	if err != nil {
		return err
	}

	clientMessages, err := executePlan(context.Background(), plan)
	if err != nil {
		return err
	}

	for _, message := range clientMessages {
		if _, err := fmt.Fprintln(stdout, message); err != nil {
			return err
		}
	}

	return nil
}
