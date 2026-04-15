package cli

import (
	"fmt"
	"io"
)

const helpText = `Hello from Anvil.

This is the initial CLI scaffold.

Usage:
  anvil [command]

Available Commands:
  help        Show this message
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
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
