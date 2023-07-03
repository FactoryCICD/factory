package command

import (
	"log"
	"strings"
)

// ValidateCommand is a Command implementation that validates factory files.
type ValidateCommand struct {
	Meta
}

func (c *ValidateCommand) Run(rawArgs []string) int {
	log.Printf("Running validate command with args: %v", rawArgs)
	return 0
}

// Help implements cli.Command.
func (*ValidateCommand) Help() string {
	helpText := `
Usage: factory validate [options]

	Validate the configuration files in a directory, referring only to the configuration without execution.

	Validate runs checks that verify whether a configuration is syntactically
	valid. It is primarily intended for verification of configuration files.
	
Options:

  TBD
`
	return strings.TrimSpace(helpText)
}

func (*ValidateCommand) Synopsis() string {
	return "Show changes required by the current configuration"
}
