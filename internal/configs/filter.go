package configs

import "github.com/hashicorp/hcl/v2"

type Filter struct {
	Config hcl.Body
}

func decodeFilterBlock(block *hcl.Block) (*Filter, hcl.Diagnostics) {
	return &Filter{
		Config: block.Body,
	}, nil
}
