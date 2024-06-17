package configs

import (
	"fmt"
	"strings"
)

type File struct {
	Filters []*Filter
	// Variables map[Scope]*Variable
	Variables *Variables
	// Stages    []*Stage
}

func NewFile() *File {
	return &File{
		Filters:   make([]*Filter, 0),
		Variables: NewVariables(),
	}
}

func (f *File) String() string {
	var sb strings.Builder
	indent := " "
	sb.WriteString("File {\n")
	sb.WriteString(fmt.Sprintf("%s Filters: [\n", indent))
	indent += " "
	for _, filter := range f.Filters {
		sb.WriteString(fmt.Sprintf("%s Filter: {\n", indent))
		sb.WriteString(fmt.Sprintf("%s  Config: %s\n", indent, filter.Config))
		sb.WriteString(fmt.Sprintf("%s }\n", indent))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	sb.WriteString(fmt.Sprintf("%s Global Variables: [\n", indent))
	indent += " "
	for k, v := range f.Variables.GlobalVariables {
		sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent, k, v.AsString()))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	sb.WriteString(fmt.Sprintf("%s Module Variables: [\n", indent))
	indent += " "
	for k, v := range f.Variables.ModuleVariables {
		sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent, k, v.AsString()))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	sb.WriteString(fmt.Sprintf("%s Resource Variables: [\n", indent))
	indent += " "
	for k, v := range f.Variables.ResourceVariables {
		sb.WriteString(fmt.Sprintf("%s {%s: %s}\n", indent, k, v.AsString()))
	}
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s ]\n", indent))
	indent = indent[:len(indent)-1]
	sb.WriteString(fmt.Sprintf("%s }\n\n", indent))
	return sb.String()
}
