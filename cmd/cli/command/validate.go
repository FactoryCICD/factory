package command

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/factoryci/factory"
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
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
	var diags hcl.Diagnostics

	paths, dirDiags := c.processDir(path)
	diags = diags.Extend(dirDiags)
	log.Printf("[DEBUG] paths found: %s", paths)
	fs := afero.NewOsFs()
	parser := factory.NewParser(fs)

	files, fileDiags := parser.LoadFiles(paths)
	diags = diags.Extend(fileDiags)
	fmt.Println(files)
	return diags
}

// processDir processes the given directory path and returns a list of file paths and any diagnostics encountered.
// If the path is a directory, it can be processed recursively if the `Recursive` flag is set.
// If the path is a file, it will be treated as a single-element slice with the file info.
// Only files with the ".hcl" extension will be included in the returned list of file paths.
//
// Parameters:
//   - path: The directory path to process.
//
// Returns:
//   - paths: A list of file paths found within the directory (including subdirectories if recursive).
//   - diags: Any diagnostics encountered during the processing.
func (c *ValidateCommand) processDir(path string) ([]string, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var paths []string

	log.Printf("[DEBUG] terraform validate: processing directory %s", path)

	info, err := os.Stat(path)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("No file or directory at %s", path),
			Detail:   err.Error(),
		})
		return paths, diags
	}

	var entries []fs.FileInfo
	if info.IsDir() {
		entries, err = ioutil.ReadDir(path)
		if err != nil {
			switch {
			case os.IsNotExist(err):
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("There is no directory at %s", path),
					Detail:   err.Error(),
				})
			default:
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Cannot read directory",
					Detail:   err.Error(),
				})
			}
			return paths, diags
		}
	} else {
		// If the path is a file, create a single-element slice with the file info
		// so that the rest of the code can treat it as paths from a directory.
		entries = []fs.FileInfo{info}
	}

	for _, info := range entries {
		name := info.Name()
		subPath := filepath.Join(path, name)
		if info.IsDir() {
			// If the path is a directory, process
			// it recursively if the flag is set.
			if c.Recursive {
				subPaths, subDiags := c.processDir(subPath)
				paths = append(paths, subPaths...)
				diags = diags.Extend(subDiags)
			}

			// If the directory is not recursive, skip it.
			// This is the default behavior
			continue
		}

		// The rest of this loop only applies to files
		ext := filepath.Ext(name)
		switch ext {
		case ".hcl":
			paths = append(paths, subPath)
		}
	}

	return paths, diags
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
