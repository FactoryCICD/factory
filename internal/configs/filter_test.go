package configs

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

func TestDecodeFilterBlock(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
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
			stages = []
	`), "test")

	pipeline, diags := file.Body.Content(pipelineBlockSchema)
	if diags.HasErrors() {
		t.Fatalf("Error decoding filter block: %s", diags)
	}

	filter, d := decodeFilterBlock(pipeline.Blocks[0])
	if d.HasErrors() {
		t.Fatalf("Error decoding filter block: %s", d)
	}

	assert.Equal(t, []string{"hi/*"}, filter.Include.Paths)
	assert.Equal(t, []string{"bye/*"}, filter.Include.Branches)
	assert.Equal(t, []string{"bar/*"}, filter.Exclude.Paths)
	assert.Equal(t, []string{"foo/*"}, filter.Exclude.Branches)
}

func TestDecodeFilterBlockReturnsErrorForIncorrectType(t *testing.T) {
	parser := hclparse.NewParser()
	file, _ := parser.ParseHCL([]byte(`
			filter {
				include {
					paths = [123]
					branches = [123]
				}
				exclude {
					paths = [123]
					branches = [123]
				}
			}
			stages = []
	`), "test")

	pipeline, diags := file.Body.Content(pipelineBlockSchema)
	if diags.HasErrors() {
		t.Fatalf("Error decoding filter block: %s", diags)
	}

	_, d := decodeFilterBlock(pipeline.Blocks[0])
	if !d.HasErrors() {
		t.Error("Expected diags to have error but was empty")
	}

	errs := d.Errs()
	assert.Len(t, errs, 4, "Expected diags to have 4 errors, got %d", len(errs))
	for _, err := range d {
		assert.Equal(t, "Value within paths or branches must be of string type", err.Summary)
		assert.Equal(t, "Invalid type for branch or path", err.Detail)
	}
}
