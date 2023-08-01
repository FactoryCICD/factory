package arguments

import (
	"flag"
	"os"

	"github.com/hashicorp/hcl/v2"
)

// Validate represents the command-line arguments for the validate command.
type Validate struct {
	// Path is the relative or absolute path to the factory configuration file
	Path string

	// Recursive is a flag that indicates whether to recursively validate all
	// subdirectories.
	Recursive bool
}

// Look into developing our own Diagnotics type
// Refer to internal/tfdiags/diagnostics.go in terraform
func ParseValidate(args []string) (*Validate, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	validate := &Validate{}

	var recursiveOutput bool
	cmdFlags := flag.NewFlagSet("validate", flag.ContinueOnError)

	cmdFlags.StringVar(&validate.Path, "path", ".", "Path to the factory configuration directory.")
	cmdFlags.BoolVar(&recursiveOutput, "recursive", false, "Recursively validate all subdirectories.")

	if err := cmdFlags.Parse(args); err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse validate command arguments",
			Detail:   err.Error(),
		})
	}

	args = cmdFlags.Args()
	if len(args) != 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect Usage",
			Detail:   "factory validate [options]",
		})
	}

	// Check that the path is a directory
	fileInfo, err := os.Stat(validate.Path)
	if os.IsNotExist(err) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid Path",
			Detail:   "The path provided does not exist.",
		})
		return nil, diags
	} else if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid Path",
			Detail:   "An error occurred while attempting to read the path provided.",
		})
		return nil, diags
	}
	if !fileInfo.IsDir() {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid Path",
			Detail:   "The path provided is not a directory.",
		})
	}

	validate.Recursive = recursiveOutput

	return validate, diags
}
