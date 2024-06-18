package configs

import "github.com/hashicorp/hcl/v2"

var runBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "command"},
		{Name: "file"},
	},
}

type RunBlock struct {
	Name     string
	Commands []string
	File     string
}
