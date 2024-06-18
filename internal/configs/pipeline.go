package configs

import (
	"log"

	"github.com/hashicorp/hcl/v2"
)

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

type StageDefinition struct {
	Name       string
	DependsOn  []string
	Namespaces []string
}

type Pipeline struct {
	Name   string
	Filter *Filter
	Stages []*StageDefinition
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Stages: make([]*StageDefinition, 0),
	}
}

func decodePipelineBlock(block *hcl.Block) (*Pipeline, hcl.Diagnostics) {
	content, diags := block.Body.Content(pipelineBlockSchema)
	pipeline := NewPipeline()
	pipeline.Name = block.Labels[0]

	for _, innerBlock := range content.Blocks {
		switch innerBlock.Type {
		case "filter":
			log.Printf("[DEBUG] Filter block found, decoding in progress")

			filterCfg, filterDiags := decodeFilterBlock(innerBlock)
			diags = append(diags, filterDiags...)
			pipeline.Filter = filterCfg
		default:
			// Should never happen beacause the above cases should be exhaustive
			// for all block type names in our schema.
			continue
		}
	}
	// Add the stages
	stages := content.Attributes["stages"]

	stagesVal, _ := stages.Expr.Value(nil)
	stageDefs := make([]*StageDefinition, 0)
	for _, el := range stagesVal.AsValueSlice() {
		elMap := el.AsValueMap()
		sd := &StageDefinition{}
		sd.Name = elMap["name"].AsString()
		if dependsOn, ok := elMap["depends_on"]; ok {
			for _, dep := range dependsOn.AsValueSlice() {
				sd.DependsOn = append(sd.DependsOn, dep.AsString())
			}
		}
		if namespaces, ok := elMap["namespaces"]; ok {
			for _, ns := range namespaces.AsValueSlice() {
				sd.Namespaces = append(sd.Namespaces, ns.AsString())
			}
		}
		stageDefs = append(stageDefs, sd)
	}
	pipeline.Stages = stageDefs

	return pipeline, diags
}
