pipeline "foo-bar" {
    source = "github.com/Cwagne17/factory-templates//pipelines/docker?ref=v1.0.0"
    variables {
        foo = "bar"
    }
}