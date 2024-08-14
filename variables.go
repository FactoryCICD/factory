package factory

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Scope string

const (
	GlobalScope Scope = "global"
	StageScope  Scope = "stage"
)

/*
Current issues with variables. Resolving the variables such as var.foo or resolving objects var.foo.bar

Scopes, there is only 1 global scope but can be many module and resource scopes.
*/
type Variables struct {
	GlobalVariables map[string]cty.Value
	StageVariables  map[string]map[string]cty.Value
}

func NewVariables() *Variables {
	return &Variables{
		GlobalVariables: make(map[string]cty.Value),
		StageVariables:  make(map[string]map[string]cty.Value),
	}
}

func (v *Variables) InsertStage(key string, value *cty.Value, scopeID string) {
	if _, ok := v.StageVariables[scopeID]; !ok {
		// Scope not found, create the scope
		v.StageVariables[scopeID] = make(map[string]cty.Value)
	}
	v.StageVariables[scopeID][key] = *value
}

func (v *Variables) InsertGlobal(key string, value *cty.Value) {
	v.GlobalVariables[key] = *value
}

func decodeGlobalVariableBlock(block *hcl.Block, file *File) hcl.Diagnostics {
	vars, diags := block.Body.JustAttributes()

	for _, attr := range vars {
		name := attr.Name
		value, d := attr.Expr.Value(file.GetEvalContext(nil))
		diags = append(diags, d...)
		file.Variables.InsertGlobal(name, &value)
	}

	return diags
}

func decodeVariableBlock(block *hcl.Block, file *File, scopeID string) hcl.Diagnostics {
	vars, diags := block.Body.JustAttributes()

	for _, attr := range vars {
		name := attr.Name
		value, d := attr.Expr.Value(file.GetEvalContext(&scopeID))
		diags = append(diags, d...)
		file.Variables.InsertStage(name, &value, scopeID)
	}

	return diags
}
