package configs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type File struct {
	Pipelines []*Pipeline
	Variables *Variables
	Stages    []*Stage
}

func (f *File) GetEvalContext(scopeID *string) *hcl.EvalContext {
	v := f.Variables
	// No Scope, get the global context
	if scopeID == nil {
		if len(v.GlobalVariables) == 0 {
			return nil
		}

		return &hcl.EvalContext{
			Variables: map[string]cty.Value{
				"var": cty.MapVal(v.GlobalVariables),
			},
		}
	}

	// Combine the stage Scope with the global scope, overriding global variables
	scope := make(map[string]cty.Value)
	// First add the global variables
	for k, v := range v.GlobalVariables {
		scope[k] = v
	}
	if stageScope, ok := v.StageVariables[*scopeID]; ok {
		// Add the stage scope
		for k, v := range stageScope {
			scope[k] = v
		}
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.MapVal(scope),
		},
	}
}

func NewFile() *File {
	return &File{
		Pipelines: make([]*Pipeline, 0),
		Variables: NewVariables(),
		Stages:    make([]*Stage, 0),
	}
}

func (f *File) String() string {
	var sb strings.Builder
	sb.WriteString("File {\n")
	// Filters
	sb.WriteString(fmt.Sprintf("%s Pipeline: [\n", indent(1)))
	for _, pipeline := range f.Pipelines {
		// Filter
		filter := pipeline.Filter
		sb.WriteString(fmt.Sprintf("%s %s: {\n", indent(2), pipeline.Name))
		sb.WriteString(fmt.Sprintf("%s Filter: {\n", indent(3)))
		sb.WriteString(fmt.Sprintf("%s Exclude: {\n", indent(4)))
		sb.WriteString(fmt.Sprintf("%s Paths: [\n", indent(5)))
		for _, path := range filter.Exclude.Paths {
			sb.WriteString(fmt.Sprintf("%s %s\n", indent(6), path))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(5)))
		sb.WriteString(fmt.Sprintf("%s Branches: [\n", indent(5)))
		for _, branch := range filter.Exclude.Branches {
			sb.WriteString(fmt.Sprintf("%s %s\n", indent(6), branch))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(5)))
		sb.WriteString(fmt.Sprintf("%s }\n", indent(4)))
		sb.WriteString(fmt.Sprintf("%s Include {\n", indent(4)))
		sb.WriteString(fmt.Sprintf("%s Paths: [\n", indent(5)))
		for _, path := range filter.Include.Paths {
			sb.WriteString(fmt.Sprintf("%s %s\n", indent(6), path))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(5)))
		sb.WriteString(fmt.Sprintf("%s Branches: [\n", indent(5)))
		for _, branch := range filter.Include.Branches {
			sb.WriteString(fmt.Sprintf("%s %s\n", indent(6), branch))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(5)))
		sb.WriteString(fmt.Sprintf("%s }\n", indent(4)))
		sb.WriteString(fmt.Sprintf("%s }\n", indent(3)))
		// Stages
		sb.WriteString(fmt.Sprintf("%s Stage Definitions: [\n", indent(3)))
		for _, stage := range pipeline.Stages {
			sb.WriteString(fmt.Sprintf("%s {\n", indent(4)))
			sb.WriteString(fmt.Sprintf("%s  Name: %s\n", indent(5), stage.Name))
			sb.WriteString(fmt.Sprintf("%s  Depends On: [\n", indent(5)))
			for _, dp := range stage.DependsOn {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent(6), dp))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(5)))
			sb.WriteString(fmt.Sprintf("%s }\n", indent(4)))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(3)))
	}
	sb.WriteString(fmt.Sprintf("%s ]\n", indent(2)))
	// Variables
	sb.WriteString(fmt.Sprintf("%s Global Variables: [\n", indent(2)))
	for k, v := range f.Variables.GlobalVariables {
		sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent(3), k, v.AsString()))
	}
	sb.WriteString(fmt.Sprintf("%s ]\n", indent(2)))
	sb.WriteString(fmt.Sprintf("%s Stage Variables: [\n", indent(2)))
	for k, v := range f.Variables.StageVariables {
		sb.WriteString(fmt.Sprintf("%s %s: [\n", indent(3), k))
		for variable, value := range v {
			sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent(4), variable, value.AsString()))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(3)))
	}
	sb.WriteString(fmt.Sprintf("%s ]\n", indent(2)))
	// Stages
	sb.WriteString(fmt.Sprintf("%s Stages: [\n", indent(2)))
	for _, stage := range f.Stages {
		sb.WriteString(fmt.Sprintf("%s %s: [\n", indent(3), stage.Name))
		for _, runBlock := range stage.RunBlocks {
			sb.WriteString(fmt.Sprintf("%s %s: [\n", indent(4), runBlock.Name))
			for _, command := range runBlock.Commands {
				sb.WriteString(fmt.Sprintf("%s {\n", indent(5)))
				sb.WriteString(fmt.Sprintf("%s command: %s\n", indent(6), command))
				sb.WriteString(fmt.Sprintf("%s }\n", indent(5)))
			}
			if runBlock.File != "" {
				sb.WriteString(fmt.Sprintf("%s { file: %s }\n", indent(5), runBlock.File))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(4)))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(3)))
	}
	sb.WriteString(fmt.Sprintf("%s ]\n", indent(2)))
	sb.WriteString(fmt.Sprintf("%s }\n\n", indent(1)))
	return sb.String()
}

func indent(spaces int) string {
	return strings.Repeat("  ", spaces)
}
