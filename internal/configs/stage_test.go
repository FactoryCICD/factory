package configs

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestDecodeStageBlock(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
		stage "stage1" {
			variables {
				foo = "bar"
			}
			run "Install Docker" {
				command = "apt-get install docker"
			}
			run "Push Docker Image" {
				file = "push.sh"
			}
		}
	`), "test")

	config, d := file.Body.Content(configFileSchema)
	if d.HasErrors() {
		t.Fatalf("Error decoding config file: %s", d)
	}

	f := NewFile()
	stageBlock, sd := decodeStageBlock(config.Blocks[0], f)
	if sd.HasErrors() {
		t.Fatalf("Error decoding stage block: %s", sd)
	}

	assert.Equal(t, "stage1", stageBlock.Name, "Expected stage name to be 'stage1' got %s", stageBlock.Name)
	assert.Equal(t, cty.StringVal("bar"), f.Variables.StageVariables["stage1"]["foo"], "Expected variable foo to be bar. got %s", f.Variables.StageVariables["stage1"]["foo"])
	assert.Len(t, stageBlock.RunBlocks, 2, "Expected run blocks to be len 2 got %d", len(stageBlock.RunBlocks))
}
