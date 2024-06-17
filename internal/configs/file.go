package configs

import (
	"fmt"
	"strings"
)

type File struct {
	Filters []*Filter
	// Variables map[Scope]*Variable
	Variables *Variables
	Stages    []*Stage
}

func NewFile() *File {
	return &File{
		Filters:   make([]*Filter, 0),
		Variables: NewVariables(),
		Stages:    make([]*Stage, 0),
	}
}

func (f *File) String() string {
	var sb strings.Builder
	indent := " "
	sb.WriteString("File {\n")
	// Filters
	sb.WriteString(fmt.Sprintf("%s Filters: [\n", indent))
	indent += " "
	for _, filter := range f.Filters {
		sb.WriteString(fmt.Sprintf("%s Filter: {\n", indent))
		sb.WriteString(fmt.Sprintf("%s  Exclude: %s\n", indent, filter.Exclude))
		sb.WriteString(fmt.Sprintf("%s  Include: %s\n", indent, filter.Include))
		sb.WriteString(fmt.Sprintf("%s }\n", indent))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	// Variables
	sb.WriteString(fmt.Sprintf("%s Global Variables: [\n", indent))
	indent += " "
	for k, v := range f.Variables.GlobalVariables {
		sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent, k, v.AsString()))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	sb.WriteString(fmt.Sprintf("%s Stage Variables: [\n", indent))
	indent += " "
	for k, v := range f.Variables.StageVariables {
		sb.WriteString(fmt.Sprintf("%s %s: [\n", indent, k))
		indent += " "
		for variable, value := range v {
			sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent, variable, value.AsString()))
		}
		indent = indent[:len(indent)-1]
		sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	// Stages
	sb.WriteString(fmt.Sprintf("%s Stages: [\n", indent))
	indent += " "
	for _, stage := range f.Stages {
		sb.WriteString(fmt.Sprintf("%s Stage - %s: [\n", indent, stage.Name))
		indent += " "
		for _, runBlock := range stage.RunBlocks {
			sb.WriteString(fmt.Sprintf("%s %s: [\n", indent, runBlock.Name))
			indent += " "
			for _, command := range runBlock.Commands {
				sb.WriteString(fmt.Sprintf("%s %s\n", indent, command))
			}
			indent = indent[:len(indent)-1]
			sb.WriteString(fmt.Sprintf("%s ]\n", indent))
		}
		indent = indent[:len(indent)-1]
		sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s }\n\n", indent))
	return sb.String()
}
