package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

var runBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "command"},
		{Name: "file"},
	},
}

type RunBlock struct {
	Name     string
	Commands []string
	File     string
}

func decodeRunBlock(block *hcl.Block, file *File) (RunBlock, hcl.Diagnostics) {
	// Decode Run block
	run, diags := block.Body.Content(runBlockSchema)

	runBlock := RunBlock{
		Name: block.Labels[0],
	}

	if command, ok := run.Attributes["command"]; ok {
		val, d := command.Expr.Value(file.Variables.GetVariableContext(block.Labels[0]))
		diags = append(diags, d...)
		fmt.Println(val)
		runBlock.Commands = append(runBlock.Commands, val.AsString())
	}

	if f, ok := run.Attributes["file"]; ok {
		val, d := f.Expr.Value(file.Variables.GetVariableContext(block.Labels[0]))
		diags = append(diags, d...)
		runBlock.File = val.AsString()
	}

	return runBlock, diags
}
