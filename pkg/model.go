package pkg

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type Include struct {
	paths    []string
	branches []string
}

type Exclude struct {
	paths    []string
	branches []string
}

type Stage struct {
	name       string
	depends_on []string
	namespaces []string
}

type Filter struct {
	include Include
	exclude Exclude
}

func (i *Include) Execute(ctx *hcl.EvalContext) error {
	fmt.Printf("Include: %s %s\n", i.paths, i.branches)
	return nil
}

func (e *Exclude) Execute(ctx *hcl.EvalContext) error {
	fmt.Printf("Exclude: %s %s\n", e.paths, e.branches)
	return nil
}

func (s *Stage) Execute(ctx *hcl.EvalContext) error {
	fmt.Printf("Stage: %s %s %s\n", s.name, s.depends_on, s.namespaces)
	return nil
}

func (f *Filter) Execute(ctx *hcl.EvalContext) error {
	f.include.Execute(ctx)
	f.exclude.Execute(ctx)
	return nil
}

type Pipeline struct {
	name   string
	filter Filter
	stages []Stage
}

func (p *Pipeline) Execute(ctx *hcl.EvalContext) error {
	fmt.Printf("Pipeline: %s\n", p.name)
	p.filter.Execute(ctx)
	for _, stage := range p.stages {
		stage.Execute(ctx)
	}
	return nil
}

type Workflow struct {
	pipelines []Pipeline
}

func (w *Workflow) Execute() {
	ctx := &hcl.EvalContext{}

	for _, pipeline := range w.pipelines {
		pipeline.Execute(ctx)
	}
}
