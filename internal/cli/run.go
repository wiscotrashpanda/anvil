package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/wiscotrashpanda/anvil/internal/reconcile"
)

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
  anvil reconcile --manifests <path>

Flags:
  --manifests string   Path to a directory containing YAML manifests
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
		return fmt.Errorf("reconcile requires --manifests <path>")
	}

	plan, err := reconcile.Load(*manifestsPath)
	if err != nil {
		return err
	}

	for _, message := range plan.Messages() {
		if _, err := fmt.Fprintln(stdout, message); err != nil {
			return err
		}
	}

	return nil
}
