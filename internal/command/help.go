package command

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// helpFunc is a cli.HelpFunc that can be used to output the help CLI instructions for Terraform.
func HelpFunc(commands map[string]cli.CommandFactory) string {
	helpText := fmt.Sprintf(`
Usage: factory <subcommand> [args]

The available commands for execution are listed below.

Subcommands:
%s
`, listCommands(Commands))

	return strings.TrimSpace(helpText)
}

// listCommands just lists the commands in the map.
func listCommands(commands map[string]cli.CommandFactory) string {
	var buf bytes.Buffer

	for command := range commands {
		buf.WriteString(fmt.Sprintf("  %s\n", command))
	}

	return buf.String()
}
