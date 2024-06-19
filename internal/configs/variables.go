package configs

import (
	"fmt"

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

func (v *Variables) Resolve(variable string, scope Scope, scopeID string) (cty.Value, bool) {
	// Check the current scope
	switch scope {
	case GlobalScope:
		return v.resolveGlobalScope(variable)
	case StageScope:
		return v.resolveStageScope(variable, scopeID)
	}

	panic(fmt.Sprintf("%s scope is not a valid scope", scope))
}

func (v *Variables) resolveGlobalScope(variable string) (cty.Value, bool) {
	vari, ok := v.GlobalVariables[variable]
	return vari, ok
}

func (v *Variables) resolveStageScope(variable, scopeID string) (cty.Value, bool) {
	// Resolve the module, if not found, check global
	_, ok := v.StageVariables[scopeID]
	if !ok {
		panic(fmt.Sprintf("Scope with ID: %s was not found.", scopeID))
	}
	vari, ok := v.StageVariables[scopeID][variable]
	if !ok {
		return v.resolveGlobalScope(variable)
	}
	return vari, ok
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
