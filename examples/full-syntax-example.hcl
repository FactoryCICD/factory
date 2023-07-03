# Build Agent will run on a linux machine
# Parses all .hcl files in the .factory directory
# Allows separation of pipelines and stage declarations into multiple files

# Supports multiple pipeline declarations
pipeline "foo-bar" {
  filter {
    include {              # Optional
      paths    = ["foo/*"] # Optional
      branches = ["bar/*"] # Optional
    }
    exclude {              # Optional
      paths    = ["bar/*"] # Optional
      branches = ["foo/*"] # Optional
    }
  }
  stages = [
    { name = "stage1" },
    {
      name       = "stage2",
      depends_on = ["stage1"],
      namespaces = [""]
    }
  ]
}

# Global Variable definition for all stages
variables {
  foo = "bar" # Optional
}

# Supports multiple "unique" stage declarations
stage "stage1" {
  variables {
    foo = "bar" # Variable overwrites the global variable for this stage
  }

  # Supports multiple "unique" run declarations
  run "Install Docker" { # Run supports command or file
    command = "apt-get install docker"
  }
  run "Build/Tag Docker Image" {
    command = <<EOT
            # Supports multiline commands
            docker build -t my-image .
            docker tag my-image:latest my-image:${var.foo}
        EOT
  }
  run "Push Docker Image" {
    command = "docker push my-image:${var.foo}" # Uses the local variable defined in the stage first then the global variable
  }
}

# Supports multiple "unique" stage declarations
stage "stage2" {
  run "Deploy Application" {
    file = "deploy.sh" # automatically makes the file executable and runs it (file must have a shebang)
  }
}
