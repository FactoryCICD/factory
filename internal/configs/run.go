package configs

import "github.com/hashicorp/hcl/v2"

var runBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "command", Required: true},
	},
}

type RunBlock struct {
	Commands []string
}
