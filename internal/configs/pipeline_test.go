package configs

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

func TestDecodePipelineBlock(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
		pipeline "test" {
			filter {
				include {
					paths = ["hi/*"]
					branches = ["bye/*"]
				}
				exclude {
					paths = ["bar/*"]
					branches = ["foo/*"]
				}
			}
			stages = [
				{ name = "stage1" },
				{
					name       = "stage2",
					depends_on = ["stage1"],
					namespaces = ["nm1"]
				}
			]
		}
`), "test")

	configFile, diags := file.Body.Content(configFileSchema)
	if diags.HasErrors() {
		t.Fatalf("Error decoding pipeline block: %s", diags)
	}

	pipeline, d := decodePipelineBlock(configFile.Blocks[0])
	if d.HasErrors() {
		t.Fatalf("Error decoding pipeline block: %s", d)
	}

	assert.NotNil(t, pipeline.Filter, "Expected pipeline.Filters to be non-nil but was nil")
	assert.NotNil(t, pipeline.Stages, "Expected pipeline.Stages to be non-nil but was nil")
	assert.Equal(t, "test", pipeline.Name, "Expected pipeline.Name to be \"name\" got: %s", pipeline.Name)
	assert.Len(t, pipeline.Stages, 2, "Expected pipeline.Stages to have len 2 got %d", pipeline.Stages)
	assert.Equal(t, "stage1", pipeline.Stages[0].Name, "Expected name of first stage to be 'stage1' got %s", pipeline.Stages[0].Name)
	assert.Nil(t, pipeline.Stages[0].DependsOn, "Expected depends on for first stage to be nil got %+v", pipeline.Stages[0].DependsOn)
	assert.Nil(t, pipeline.Stages[0].Namespaces, "Expected Namespaces for first stage to be nil got %+v", pipeline.Stages[0].Namespaces)
	assert.Equal(t, "stage2", pipeline.Stages[1].Name, "Expected 2nd stage name to be 'stage2' got %s", pipeline.Stages[1].Name)
	assert.Equal(t, []string{"stage1"}, pipeline.Stages[1].DependsOn, "Expected 2nd stage depends_on to be []string{'stage1'} got %+v", pipeline.Stages[1].DependsOn)
	assert.Equal(t, []string{"nm1"}, pipeline.Stages[1].Namespaces, "Expected 2nd stage namespaces to be []string{'nm1'} got %+v", pipeline.Stages[1].DependsOn)
}
