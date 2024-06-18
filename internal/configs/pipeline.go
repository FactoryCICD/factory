package configs

type StageDefinition struct {
	Name       string
	DependsOn  []string
	Namespaces []string
}

type Pipeline struct {
	Name    string
	Filters []*Filter
	Stages  []*StageDefinition
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Filters: make([]*Filter, 0),
		Stages:  make([]*StageDefinition, 0),
	}
}
