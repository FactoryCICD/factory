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
