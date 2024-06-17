package configs

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// LoadConfigFile loads a configuration file from the specified path and returns
// the parsed file along with any diagnostics encountered during parsing.
// It first loads the HCL file from the given path and then parses its content
// according to the defined schema. The function supports parsing pipeline blocks
// and filter blocks within them. It logs debug messages for each block type found.
// Currently, decoding of variables and stage blocks is not implemented.
// The function returns the parsed file and any encountered diagnostics.
// check out terraoform\internal\config\parser_config.go line 51
func (p *Parser) LoadConfigFile(path string) (*File, hcl.Diagnostics) {
	body, diags := p.LoadHCLFile(path)
	if body == nil {
		return nil, diags
	}

	file := &File{}

	content, contentDiags := body.Content(configFileSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "pipeline":
			log.Printf("[DEBUG] Pipeline block found, decoding in progress")

			content, contentDiags := block.Body.Content(pipelineBlockSchema)
			diags = append(diags, contentDiags...)

			for _, innerBlock := range content.Blocks {
				switch innerBlock.Type {

				case "filter":
					log.Printf("[DEBUG] Filter block found, decoding in progress")

					filterCfg, filterDiags := decodeFilterBlock(innerBlock)
					diags = append(diags, filterDiags...)
					if filterCfg != nil {
						file.Filters = append(file.Filters, filterCfg)
					}

				default:
					// Should never happen beacause the above cases should be exhaustive
					// for all block type names in our schema.
					continue
				}
			}

		// Check out line 493 of internal/configs/named_values.go in terraform
		case "variables":
			log.Printf("[DEBUG] Variables block found, decoding in progress")
			variables, varDiag := decodeVariableBlock(block.Body, file, GlobalScope)
			diags = append(diags, varDiag...)
			file.Variables = append(file.Variables, variables)
		case "stage":
			log.Printf("[DEBUG] Stage block found, decoding not yet implemented")
			content, contDiags := block.Body.Content(stageBlockSchema)
			diags = append(diags, contDiags...)

			for _, inner := range content.Blocks {
				switch inner.Type {
				case "variables":
					vars, varDiags := decodeVariableBlock(inner.Body, file)
					diags = append(diags, varDiags...)
					fmt.Println(vars)

				case "run":
					fmt.Println(inner.Labels[0])
					fmt.Println(inner.Body)
					// Decode Run block
					run, d := inner.Body.Content(runBlockSchema)
					diags = append(diags, d...)

					runBlock := RunBlock{}
					for _, attr := range run.Attributes {
						if len(attr.Expr.Variables()) > 0 {
							fmt.Println("Expression has variables")
							fmt.Println(attr.Expr.Variables())
							val, _ := attr.Expr.Value(&hcl.EvalContext{
								Variables: map[string]cty.Value{
									"var": cty.MapVal(map[string]cty.Value{"foo": cty.StringVal("hello")}),
								},
							})

							fmt.Println(val.AsString())
						} else {
							values, _ := attr.Expr.Value(&hcl.EvalContext{})

							fmt.Println(values.AsString())
						}

						runBlock.Commands = append(runBlock.Commands, attr.Name)
					}
				}
			}
		default:
			// Should never happen beacause the above cases should be exhaustive
			// for all block type names in our schema.
			continue
		}
	}

	return file, diags
}

// configurationFileSchema is the schema for the top-level of a config file. We use
// the low-level HCL API for this level so we can easily deal with each
// block type seaparately with its own decoding logic.
var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "pipeline",
			LabelNames: []string{"name"},
		},
		{
			Type: "variables",
		},
		{
			Type:       "stage",
			LabelNames: []string{"name"},
		},
	},
}

// pipelineBlockSchema is the schema for a top-level "pipeline" block in
// a configuration file.
var pipelineBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "stages", Required: true},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "filter",
		},
	},
}
