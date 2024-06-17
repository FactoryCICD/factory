package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Scope string

const (
	GlobalScope   Scope = "global"
	ModuleScope   Scope = "module"
	ResourceScope Scope = "resource"
)

type Variables struct {
	GlobalVariables   map[string]cty.Value
	ModuleVariables   map[string]cty.Value
	ResourceVariables map[string]cty.Value
}

func (v *Variables) Resolve(variable string, scope Scope) (cty.Value, bool) {
	// Check the current scope
	switch scope {
	case GlobalScope:
		return v.resolveGlobalScope(variable)
	case ModuleScope:
		return v.resolveModuleScope(variable)
	case ResourceScope:
		return v.resolveResourceScope(variable)
	}

	panic(fmt.Sprintf("%s scope is not a valid scope", scope))
}

func (v *Variables) resolveGlobalScope(variable string) (cty.Value, bool) {
	vari, ok := v.GlobalVariables[variable]
	return vari, ok
}

func (v *Variables) resolveModuleScope(variable string) (cty.Value, bool) {
	// Resolve the module, if not found, check global
	vari, ok := v.ModuleVariables[variable]
	if !ok {
		return v.resolveGlobalScope(variable)
	}
	return vari, ok
}

func (v *Variables) resolveResourceScope(variable string) (cty.Value, bool) {
	// check resource scope, if not found, go up
	vari, ok := v.ResourceVariables[variable]
	if !ok {
		return v.resolveModuleScope(variable)
	}

	return vari, ok
}

func (v *Variables) Insert(key string, value *cty.Value, scope Scope) {
	switch scope {
	case GlobalScope:
		v.GlobalVariables[key] = *value
	case ModuleScope:
		v.ModuleVariables[key] = *value
	case ResourceScope:
		v.ResourceVariables[key] = *value
	}
	panic(fmt.Sprintf("%s is not a valid scope!\n", key))
}

func decodeVariableBlock(body hcl.Body, file *File, scope Scope) (*Variables, hcl.Diagnostics) {
	variables := &Variables{}

	vars, diags := body.JustAttributes()

	for _, attr := range vars {
		name := attr.Name
		var variable cty.Value
		// Convert the attr to map[string]cty.Value
		if len(attr.Expr.Variables()) > 0 {
			// Attribute expression has variables, will need to look at scope to see
			// if all variables are defined
			val, _ := attr.Expr.Value(&hcl.EvalContext{})
			variable = val
		} else {
			// No Variables, constant value
			val, _ := attr.Expr.Value(nil)
			variable = val
		}

		file.Variables.Insert(name, &variable, scope)
	}

	return variables, diags
}
