package configs

import (
	"github.com/hashicorp/hcl/v2"
)

type Include struct {
	Paths    []string
	Branches []string
}

type Exclude struct {
	Paths    []string
	Branches []string
}

type Filter struct {
	Include Include
	Exclude Exclude
}

var filterBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "stages"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "include"},
		{Type: "exclude"},
	},
}

func decodeFilterBlock(block *hcl.Block) (*Filter, hcl.Diagnostics) {
	content, diags := block.Body.Content(filterBlockSchema)
	filter := &Filter{}

	// Decode blocks
	for _, block := range content.Blocks {
		switch block.Type {
		case "include":
			include, incDiags := decodeIncludeOrExcludeBlock(block)
			diags = append(diags, incDiags...)
			filter.Include = *include.(*Include)
		case "exclude":
			exclude, exDiags := decodeIncludeOrExcludeBlock(block)
			diags = append(diags, exDiags...)
			filter.Exclude = *exclude.(*Exclude)
		}
	}

	return filter, diags
}

func decodeIncludeOrExcludeBlock(block *hcl.Block) (interface{}, hcl.Diagnostics) {
	attributes, diags := block.Body.JustAttributes()
	var result interface{}
	switch block.Type {
	case "include":
		include := Include{}
		result = &include
	case "exclude":
		exclude := Exclude{}
		result = &exclude
	}

	for _, attr := range attributes {
		switch attr.Name {
		case "paths":
			paths := decodeStringSliceAttribute(attr, &diags)
			if block.Type == "include" {
				result.(*Include).Paths = paths
			} else {
				result.(*Exclude).Paths = paths
			}
		case "branches":
			branches := decodeStringSliceAttribute(attr, &diags)
			if block.Type == "include" {
				result.(*Include).Branches = branches
			} else {
				result.(*Exclude).Branches = branches
			}
		}
	}

	return result, diags
}

func decodeStringSliceAttribute(attr *hcl.Attribute, diags *hcl.Diagnostics) []string {
	var result []string
	p, d := attr.Expr.Value(nil)
	*diags = append(*diags, d...)
	for _, val := range p.AsValueSlice() {
		result = append(result, val.AsString())
	}
	return result
}
