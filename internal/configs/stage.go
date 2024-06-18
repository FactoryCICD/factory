package configs

import "github.com/hashicorp/hcl/v2"

type Stage struct {
	Name      string
	RunBlocks []RunBlock
}

var stageBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "variables"},
		{Type: "run", LabelNames: []string{"name"}},
	},
}

func decodeStageBlock(block *hcl.Block, file *File) (*Stage, hcl.Diagnostics) {
	content, diags := block.Body.Content(stageBlockSchema)
	stage := &Stage{
		Name: block.Labels[0],
	}

	for _, inner := range content.Blocks {
		switch inner.Type {
		case "variables":
			varDiags := decodeVariableBlock(inner.Body, file, StageScope, block.Labels[0])
			diags = append(diags, varDiags...)
		case "run":
			runBlock, rbDiags := decodeRunBlock(inner, file)
			diags = append(diags, rbDiags...)
			stage.RunBlocks = append(stage.RunBlocks, runBlock)
		}
	}

	return stage, diags
}
