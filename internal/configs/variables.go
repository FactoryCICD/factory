package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Scope string

const (
	GlobalScope Scope = "global"
	StageScope  Scope = "module"
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

func (v *Variables) Insert(key string, value *cty.Value, scope Scope, scopeId string) {
	fmt.Printf("inserting %s -> %s", key, value)
	switch scope {
	case GlobalScope:
		v.GlobalVariables[key] = *value
	case StageScope:
		if _, ok := v.StageVariables[scopeId]; !ok {
			// SCope not found, create the scope
			v.StageVariables[scopeId] = make(map[string]cty.Value)
		}
		v.StageVariables[scopeId][key] = *value
	default:
		panic(fmt.Sprintf("%s is not a valid scope!", scope))
	}
}

func (v *Variables) GetVariableContext(scopeID string) *hcl.EvalContext {
	stageScope, ok := v.StageVariables[scopeID]
	if !ok {
		// Stage scope was not found, just return the global scope
		return &hcl.EvalContext{
			Variables: v.GlobalVariables,
		}
	}

	// Combine the stage Scope with the global scope, overriding global variables
	scope := make(map[string]cty.Value)
	// First add the global variables
	for k, v := range v.GlobalVariables {
		scope[k] = v
	}
	// Add the stage scope
	for k, v := range stageScope {
		scope[k] = v
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.MapVal(scope),
		},
	}
}

func decodeVariableBlock(body hcl.Body, file *File, scope Scope, scopeID string) hcl.Diagnostics {
	vars, diags := body.JustAttributes()

	for _, attr := range vars {
		name := attr.Name
		value, _ := attr.Expr.Value(file.Variables.GetVariableContext(scopeID))

		file.Variables.Insert(name, &value, scope, scopeID)
	}

	return diags
}
