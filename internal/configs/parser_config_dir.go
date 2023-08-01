package configs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// func (p *Parser) LoadConfigDir(path string) (*Module, hcl.Diagnostics) {
// 	primaryPaths, overridePaths, _, diags := p.dirFiles(path, "")
// 	if diags.HasErrors() {
// 		return nil, diags
// 	}

// 	primary, fDiags := p.loadFiles(primaryPaths, false)
// 	diags = append(diags, fDiags...)
// 	override, fDiags := p.loadFiles(overridePaths, true)
// 	diags = append(diags, fDiags...)

// 	mod, modDiags := NewModule(primary, override)
// 	diags = append(diags, modDiags...)

// 	mod.SourceDir = path

// 	return mod, diags
// }

func (p *Parser) LoadFiles(paths []string) ([]*File, hcl.Diagnostics) {
	var files []*File
	var diags hcl.Diagnostics

	for _, path := range paths {
		var f *File
		var fDiags hcl.Diagnostics

		f, fDiags = p.LoadConfigFile(path)
		diags = append(diags, fDiags...)
		if f != nil {
			files = append(files, f)
		}
	}

	return files, diags
}

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
