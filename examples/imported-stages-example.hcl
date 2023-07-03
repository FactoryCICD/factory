pipeline "foo-bar" {
    description = "foo-bar" # Optional
    filter {
        include { # Optional
            paths = ["foo/*"] # Optional
            branches = ["bar/*"] # Optional
        }
        exclude { # Optional
            paths = ["bar/*"] # Optional
            branches = ["foo/*"] # Optional
        }
    }
    
    stages {
        stage1 { # Uses imported stage1 definition
            # Will look for a job named "stage1" in the docker directory of the factory-templates repository
            source = "github.com/Cwagne17/factory-templates//stages/docker?ref=v1.0.0"
            
            depends_on = ["stage1"] # Optional
            namespaces = [""] # Optional Vault namespace/mount to use for this stage
            variables {
                foo = "bar" # Overwrites the stage defined variable for the stage
            }
        }
    }
}