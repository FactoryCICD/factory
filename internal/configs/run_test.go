package configs

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

func TestDecodeRunBlock(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
		run "Install Docker" {
			command = "apt-get install docker"
		}

		run "Push Docker Image" {
			file = "push.sh"
		}
	`), "test")

	expected := []struct {
		Name    string
		Command []string
		File    string
	}{
		{
			Name:    "Install Docker",
			Command: []string{"apt-get install docker"},
			File:    "",
		},
		{
			Name:    "Push Docker Image",
			Command: nil,
			File:    "push.sh",
		},
	}

	stage, diags := file.Body.Content(stageBlockSchema)
	if diags.HasErrors() {
		t.Fatalf("Error decoding stage block: %s", diags)
	}

	f := NewFile()
	for i, block := range stage.Blocks {
		runBlock, d := decodeRunBlock(block, f, "stage")
		if d.HasErrors() {
			t.Fatalf("Error decoding run block: %s", d)
		}

		assert.Equal(t, expected[i].Name, runBlock.Name, "Expected %s got %s", expected[i].Name, runBlock.Name)
		assert.Equal(t, expected[i].Command, runBlock.Commands, "expected commands to be %+v got %+v", expected[i].Command, runBlock.Commands)
		assert.Equal(t, expected[i].File, runBlock.File, "Expected runBlock.File to be %s, got %s", expected[i].File, runBlock.File)
	}
}
