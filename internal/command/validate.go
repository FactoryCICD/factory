package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/factorycicd/factory/internal/command/arguments"
	"github.com/factorycicd/factory/internal/configs"
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
)

// ValidateCommand is a Command implementation that validates factory files.
type ValidateCommand struct {
	Meta
}

func (c *ValidateCommand) Run(rawArgs []string) int {
	// Parse and validate flags
	args, diags := arguments.ParseValidate(rawArgs)
	if diags.HasErrors() {
		c.UI.Error(diags.Error())
		cmd := &ValidateCommand{}
		cmd.Help()
		return 1
	}

	dir, err := filepath.Abs(args.Path)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unable to locate module",
			Detail:   err.Error(),
		})
	}

	// If recursive is set to true, walk the directory and validate all subdirectories
	if args.Recursive {
		walkDiags := filepath.Walk(args.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path)
			if info.IsDir() {
				validateDiags := c.validate(path)
				diags = diags.Extend(validateDiags)
			}
			return nil
		})
		if walkDiags != nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unable to validate directory",
				Detail:   walkDiags.Error(),
			})
		}
	} else {
		dirDiags := c.validate(dir)
		diags = diags.Extend(dirDiags)
	}

	if diags.HasErrors() {
		for _, diag := range diags {
			c.UI.Error(diag.Error())
		}
		return 1
	}

	return 0 // have an error parser instead of returning 0
}

func (c *ValidateCommand) validate(dir string) hcl.Diagnostics {
	var diags hcl.Diagnostics

	fs := afero.NewOsFs()
	parser := configs.NewParser(fs)

	paths, dirDiags := parser.DirFiles(dir)
	if dirDiags.HasErrors() {
		diags = diags.Extend(dirDiags)
		return diags
	}

	_, fileDiags := parser.LoadFiles(paths)
	if fileDiags.HasErrors() {
		diags = diags.Extend(fileDiags)
	}

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
