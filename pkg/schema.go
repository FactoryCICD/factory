package pkg

import "github.com/hashicorp/hcl/v2"

/**
Schema definition for a pipeline block

pipeline "foo-bar" {
  filter {
    include {              # Optional
      paths    = ["foo/*"] # Optional
      branches = ["bar/*"] # Optional
    }
    exclude {              # Optional
      paths    = ["bar/*"] # Optional
      branches = ["foo/*"] # Optional
    }
  }
  stages = [
    { name = "stage1" },
    {
      name       = "stage2",
      depends_on = ["stage1"],
      namespaces = [""]
    }
  ]
}
**/

const includeId = "include"

var (
	includeLabels = []string{}
	includeSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
		Attributes: []hcl.AttributeSchema{
			{Name: "paths", Required: false},
			{Name: "branches", Required: false},
		},
	}
)

const excludeId = "exclude"

var (
	excludeLabels = []string{}
	excludeSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
		Attributes: []hcl.AttributeSchema{
			{Name: "paths", Required: false},
			{Name: "branches", Required: false},
		},
	}
)

const filterId = "filter"

var (
	filterLabels = []string{}
	filterSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: includeId, LabelNames: includeLabels},
			{Type: excludeId, LabelNames: excludeLabels},
		},
		Attributes: []hcl.AttributeSchema{},
	}
)

const pipelineId = "pipeline"

var (
	pipelineLabels = []string{"name"}
	pipelineSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: filterId, LabelNames: filterLabels},
		},
		Attributes: []hcl.AttributeSchema{
			{Name: "stages", Required: true},
		},
	}
)

// const variablesId = "variables"

// const stageId = "stage"

var schema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: pipelineId, LabelNames: pipelineLabels},
	},
}
