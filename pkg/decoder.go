package pkg

import (
	"errors"

	"github.com/hashicorp/hcl/v2"
)

func Decode(body hcl.Body) (*Workflow, error) {
	ctx := &hcl.EvalContext{}
	workflow := &Workflow{
		pipelines: make([]Pipeline, 0),
	}
	bc, _ := body.Content(schema)

	if len(bc.Blocks) == 0 {
		return nil, errors.New("at least one pipeline block is required")
	}
	blocks := bc.Blocks.ByType()
	for blockId := range blocks {
		switch blockId {
		case pipelineId:
			for _, block := range blocks[blockId] {
				pipeline := new(Pipeline)
				err := pipeline.FromHCLBlock(block, ctx)
				if err != nil {
					return nil, err
				}
				workflow.pipelines = append(workflow.pipelines, *pipeline)
			}
		}
	}
	return workflow, nil
}

func (p *Pipeline) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	bc, d := block.Body.Content(pipelineSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}
	for _, subBlock := range bc.Blocks {
		switch subBlock.Type {
		case filterId:
			filter := new(Filter)
			if err := filter.FromHCLBlock(subBlock, ctx); err != nil {
				return err
			}
			p.filter = *filter
		}
	}
	p.name = block.Labels[0]
	return nil
}

func (f *Filter) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	bc, d := block.Body.Content(filterSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}
	for _, subBlock := range bc.Blocks {
		switch subBlock.Type {
		case includeId:
			include := new(Include)
			if err := include.FromHCLBlock(subBlock, ctx); err != nil {
				return err
			}
			f.include = *include
		case excludeId:
			// exclude := new(Exclude)
			// if err := exclude.FromHCLBlock(subBlock, ctx); err != nil {
			// 	return err
			// }
			// f.exclude = *exclude
		}
	}
	return nil
}

func (i *Include) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	bc, d := block.Body.Content(includeSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}
	if attr, ok := bc.Attributes["paths"]; ok {
		paths, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return d.Errs()[0]
		}
		for _, path := range paths.AsValueSlice() {
			i.paths = append(i.paths, path.AsString())
		}
		println(i.paths)
	}
	if attr, ok := bc.Attributes["branches"]; ok {
		branches, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return d.Errs()[0]
		}
		for _, branch := range branches.AsValueSlice() {
			i.branches = append(i.branches, branch.AsString())
		}
		println(i.branches)
	}
	return nil
}
