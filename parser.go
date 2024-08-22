package factory

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/afero"
)

// Parser is the main interface to read configuration files and other related
// files from disk.
//
// It retains a cache of all files that are loaded so that they can be used
// to create source code snippets in diagnostics, etc.
type Parser struct {
	fs afero.Afero
	p  *hclparse.Parser
}

// NewParser creates and returns a new Parser that reads files from the given
// filesystem. If a nil filesystem is passed then the system's "real" filesystem
// will be used, via afero.OsFs.
func NewParser(fs afero.Fs) *Parser {
	if fs == nil {
		fs = afero.OsFs{}
	}

	return &Parser{
		fs: afero.Afero{Fs: fs},
		p:  hclparse.NewParser(),
	}
}

// LoadHCLFile is a low-level method that reads the file at the given path,
// parses it, and returns the hcl.Body representing its root. In many cases
// it is better to use one of the other Load*File methods on this type,
// which additionally decode the root body in some way and return a higher-level
// construct.
//
// If the file cannot be read at all -- e.g. because it does not exist -- then
// this method will return a nil body and error diagnostics. In this case
// callers may wish to ignore the provided error diagnostics and produce
// a more context-sensitive error instead.
//
// The file will be parsed using the HCL native syntax unless the filename
// ends with ".json", in which case the HCL JSON syntax will be used.
func (p *Parser) LoadHCLFile(path string) (hcl.Body, hcl.Diagnostics) {
	src, err := p.fs.ReadFile(path)
	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The file %q could not be read.", path),
			},
		}
	}

	var file *hcl.File
	var diags hcl.Diagnostics

	// Only files ending in ".hcl" are parsed using the HCL native syntax.
	file, diags = p.p.ParseHCL(src, path)

	// If the returned file or body is nil, then we'll return a non-nil empty
	// body so we'll meet our contract that nil means an error reading the file.
	if file == nil || file.Body == nil {
		return hcl.EmptyBody(), diags
	}

	return file.Body, diags
}

// ParseFactoryDirectory parses the "factory/" directory and returns a slice of File pointers and any diagnostics encountered.
//
// Returns:
//   - []*File: A slice of pointers to the parsed files.
//   - hcl.Diagnostics: Any diagnostics encountered during parsing.
func ParseFactoryDirectory() ([]*File, hcl.Diagnostics) {
	return ParseDirectory(".factory/", true)
}

// ParseDirectory parses the directory at the given path and returns a slice of File pointers and any diagnostics encountered.
//
// If the recursive flag is set to true, it will recursively process subdirectories as well.
//
// Parameters:
//   - path: The directory path to process.
//   - recursive: A flag indicating whether to process subdirectories recursively.
//
// Returns:
//   - []*File: A slice of pointers to the parsed files.
//   - hcl.Diagnostics: Any diagnostics encountered during parsing.
func ParseDirectory(path string, recursive bool) ([]*File, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	paths, dirDiags := processDir(path, recursive)
	diags = diags.Extend(dirDiags)
	fs := afero.NewOsFs()
	parser := NewParser(fs)

	files, fileDiags := parser.LoadFiles(paths)
	diags = diags.Extend(fileDiags)

	return files, diags
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
func processDir(path string, recursive bool) ([]string, hcl.Diagnostics) {
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
			if recursive {
				subPaths, subDiags := processDir(subPath, recursive)
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
