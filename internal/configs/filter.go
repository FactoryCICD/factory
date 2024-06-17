package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Include struct {
	Paths    cty.Value
	Branches cty.Value
}

type Exclude struct {
	Paths    cty.Value
	Branches cty.Value
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
			include := Include{}
			includeAttributes, incDiags := block.Body.JustAttributes()
			for _, attr := range includeAttributes {
				switch attr.Name {
				case "paths":
					include.Paths, _ = attr.Expr.Value(nil)
				case "branches":
					include.Branches, _ = attr.Expr.Value(nil)
				}
			}
			diags = append(diags, incDiags...)
			filter.Include = include
		case "exclude":
			exclude := Exclude{}
			excludeAttr, exDiag := block.Body.JustAttributes()
			for _, attr := range excludeAttr {
				switch attr.Name {
				case "paths":
					exclude.Paths, _ = attr.Expr.Value(nil)
				case "branches":
					exclude.Branches, _ = attr.Expr.Value(nil)
				}
			}
			diags = append(diags, exDiag...)
			filter.Exclude = exclude
		}
	}

	// Decode stages attribute
	stages := content.Attributes["stages"]
	fmt.Println(stages.Expr)
	return filter, diags
}
