package configs

import (
	"log"

	"github.com/hashicorp/hcl/v2"
)

// check out terraoform\internal\config\parser_config.go line 51
func (p *Parser) loadConfigFile(path string) (*File, hcl.Diagnostics) {
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
			log.Printf("[INFO] Pipeline block found, decoding in progress")

			content, contentDiags := block.Body.Content(pipelineBlockSchema)
			diags = append(diags, contentDiags...)

			for _, innerBlock := range content.Blocks {
				switch innerBlock.Type {

				case "filter":
					log.Printf("[INFO] Filter block found, decoding in progress")

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
			log.Printf("[INFO] Variables block found, decoding not yet implemented")

		case "stage":
			log.Printf("[INFO] Stage block found, decoding not yet implemented")

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
