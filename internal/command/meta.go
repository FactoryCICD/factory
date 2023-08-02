package command

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/mitchellh/cli"
)

type Meta struct {
	WorkingDir string

	Ui cli.Ui
}

func (m *Meta) showDiagnostics(diags hcl.Diagnostics) {
	for _, diag := range diags {
		switch diag.Severity {
		case hcl.DiagError:
			m.Ui.Error(diag.Error())
		case hcl.DiagWarning:
			m.Ui.Warn(diag.Error())
		default:
			m.Ui.Output(diag.Error())
		}
	}
}
