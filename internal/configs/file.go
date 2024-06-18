package configs

import (
	"fmt"
	"strings"
)

type File struct {
	Pipelines []*Pipeline
	Variables *Variables
	Stages    []*Stage
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
		sb.WriteString(fmt.Sprintf("%s %s: {\n", indent(2), pipeline.Name))
		sb.WriteString(fmt.Sprintf("%s Filters: [\n", indent(3)))
		for _, filter := range f.Pipelines[0].Filters {
			sb.WriteString(fmt.Sprintf("%s Filter: {\n", indent(4)))
			sb.WriteString(fmt.Sprintf("%s Exclude: {\n", indent(5)))
			sb.WriteString(fmt.Sprintf("%s Paths: [\n", indent(6)))
			for _, path := range filter.Exclude.Paths {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent(7), path))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(6)))
			sb.WriteString(fmt.Sprintf("%s Branches: [\n", indent(6)))
			for _, branch := range filter.Exclude.Branches {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent(7), branch))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(6)))
			sb.WriteString(fmt.Sprintf("%s }\n", indent(5)))
			sb.WriteString(fmt.Sprintf("%s Include {\n", indent(5)))
			sb.WriteString(fmt.Sprintf("%s Paths: [\n", indent(6)))
			for _, path := range filter.Include.Paths {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent(7), path))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(6)))
			sb.WriteString(fmt.Sprintf("%s Branches: [\n", indent(6)))
			for _, branch := range filter.Include.Branches {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent(7), branch))
			}
			sb.WriteString(fmt.Sprintf("%s ]\n", indent(6)))
			sb.WriteString(fmt.Sprintf("%s }\n", indent(5)))
			sb.WriteString(fmt.Sprintf("%s }\n", indent(4)))
		}
		sb.WriteString(fmt.Sprintf("%s ]\n", indent(3)))
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
