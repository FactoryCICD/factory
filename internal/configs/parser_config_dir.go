package configs

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// LoadFiles loads multiple configuration files from the specified paths
// and returns a slice of File objects along with any diagnostics encountered.
func (p *Parser) LoadFiles(paths []string) ([]*File, hcl.Diagnostics) {
	var files []*File
	var diags hcl.Diagnostics

	for _, path := range paths {
		var f *File
		var fDiags hcl.Diagnostics
		log.Printf("[DEBUG] loading config file: %s", path)
		f, fDiags = p.LoadConfigFile(path)
		diags = append(diags, fDiags...)
		if f != nil {
			files = append(files, f)
		}
	}

	return files, diags
}

// This is used to find all the files in a module directory.
// It does this by looking for all files with the .hcl extension
// in the immediate directory.
func (p *Parser) DirFiles(dir string) (primary []string, diags hcl.Diagnostics) {
	infos, err := p.fs.ReadDir(dir)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read module directory",
			Detail:   fmt.Sprintf("Module directory %s does not exist or cannot be read.", dir),
		})
		return
	}

	for _, info := range infos {
		name := info.Name()

		if strings.HasSuffix(name, ".hcl") {
			fullPath := filepath.Join(dir, name)
			primary = append(primary, fullPath)
		}
	}

	return primary, diags
}
