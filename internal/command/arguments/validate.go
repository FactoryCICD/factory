package arguments

// Validate represents the command-line arguments for the validate command.
type Validate struct {
	// Path is the relative or absolute path to the factory configuration file
	Path string
}

// Look into developing our own Diagnotics type
// Refer to internal/tfdiags/diagnostics.go in terraform
func ParseValidate(args []string) (*Validate, error) {
	return &Validate{
		Path: args[0],
	}, nil
}
