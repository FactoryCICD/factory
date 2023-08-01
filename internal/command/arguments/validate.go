package arguments

import (
	"flag"

	"github.com/hashicorp/hcl/v2"
)

// Validate represents the command-line arguments for the validate command.
type Validate struct {
	// Path is the relative or absolute path to the factory configuration file
	Path string
}

// Look into developing our own Diagnotics type
// Refer to internal/tfdiags/diagnostics.go in terraform
func ParseValidate(args []string) (*Validate, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	validate := &Validate{
		Path: ".",
	}

	cmdFlags := flag.NewFlagSet("validate", flag.ContinueOnError)

	if err := cmdFlags.Parse(args); err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse validate command arguments",
			Detail:   err.Error(),
		})
	}

	args = cmdFlags.Args()
	if len(args) > 1 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Too many arguments",
			Detail:   "The validate command accepts at most one argument.",
		})
	}

	if len(args) > 0 {
		validate.Path = args[0]
	}

	return validate, diags
}
