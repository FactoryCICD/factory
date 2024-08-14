package command

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/factorycicd/factory"
	"github.com/hashicorp/hcl/v2"
)

// ValidateCommand is a Command implementation that validates factory files.
type ValidateCommand struct {
	Meta

	// Path is the relative or absolute path to the factory configuration file
	Path string

	// Recursive is a flag that indicates whether to recursively validate all
	// subdirectories.
	Recursive bool
}

// Run executes the validate command and returns an exit code.
func (c *ValidateCommand) Run(rawArgs []string) int {
	// Parse the command arguments
	cmdFlags := flag.NewFlagSet("validate", flag.ContinueOnError)
	cmdFlags.StringVar(&c.Path, "path", ".", "Path to the factory configuration directory.")
	cmdFlags.BoolVar(&c.Recursive, "recursive", false, "Recursively validate all subdirectories.")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(rawArgs); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to parse validate command arguments: %s\n", err.Error()))
		return 1
	}

	dir, err := filepath.Abs(c.Path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to get absolute path for %s: %s\n", c.Path, err.Error()))
	}

	diags := c.validate(dir)

	c.showDiagnostics(diags)
	if diags.HasErrors() {
		return 2
	}

	return 0
}

// validate validates the given path by processing the directory, loading the files,
// and returning any diagnostics encountered during the process.
func (c *ValidateCommand) validate(path string) hcl.Diagnostics {
	_, diags := factory.ParseDirectory(path, c.Recursive)

	return diags
}

// Help implements cli.Command.
func (*ValidateCommand) Help() string {
	helpText := `
Usage: factory validate [options]

	Validate the configuration files in a directory, referring only to the configuration without execution.

	Validate runs checks that verify whether a configuration is syntactically
	valid. It is primarily intended for verification of configuration files.
	
Options:

  -path <path> Path to the directory to validate. Defaults to the current directory.
  -recursive   Recursively validate all subdirectories as well.
`
	return strings.TrimSpace(helpText)
}

func (*ValidateCommand) Synopsis() string {
	return "Show changes required by the current configuration"
}
