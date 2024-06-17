package configs

import "github.com/hashicorp/hcl/v2"

type Step struct {
	Command string
	File    string
}

type Stage struct {
	Variable []*Variables
	Steps    []*Step
}

var stageBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "variables"},
		{Type: "run", LabelNames: []string{"name"}},
	},
}
