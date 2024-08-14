package factory

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestDecodeGlobalVariableBlock(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
		variables {
			foo = "global"
			bar = "hello"
		}
		stage "stage1" {
			variables {
				foo = "bar"
			}
			run "Push Docker Image" {
				command = "this is a test ${var.foo}-${var.bar}"
			}
		}
	`), "test")

	config, diags := file.Body.Content(configFileSchema)
	if diags.HasErrors() {
		t.Fatalf("Error decoding config: %s", diags)
	}

	f := NewFile()

	// Decode the variables block
	d := decodeGlobalVariableBlock(config.Blocks[0], f)
	if d.HasErrors() {
		t.Fatalf("Error decoding global variable block: %s", d)
	}

	assert.Len(t, f.Variables.GlobalVariables, 2)
	assert.Equal(t, cty.StringVal("global"), f.Variables.GlobalVariables["foo"])
	assert.Equal(t, cty.StringVal("hello"), f.Variables.GlobalVariables["bar"])

	// Decode stage and ensure that the 'foo' global var was overwritten by the local one
	stageBlock, sd := decodeStageBlock(config.Blocks[1], f)
	if sd.HasErrors() {
		t.Fatalf("Error decoding stage block: %s", sd)
	}

	assert.Equal(t, "this is a test bar-hello", stageBlock.RunBlocks[0].Commands[0])
}
