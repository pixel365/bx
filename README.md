# Bitrix CLI

Bitrix CLI is a command-line tool for developers of 1C-Bitrix platform modules. 
The project is currently in active development, and its API may change without backward compatibility.

### Features

- Manage developer accounts
- Maintain a module registry
- Build and prepare a module bundle for publication in the 1C-Bitrix Marketplace

### Installation

The installation process will be described later, as the project is still under development.

### Usage

```shell
# Create a new module (default config)
bx create --name my_module
```

```shell
# Check the configuration of a module by name
bx check --name my_module


# Check the configuration of a module by file path
bx check -f module-path/config.yaml
```

```shell
# Build a module by name
bx build --name my_module

# Build a module by file path
bx build -f config.yaml

# Override version
bx build --name my_module --version 1.2.3
```

### Example of default module configuration

```yaml
name: test  # The name of the project or build.
version: 1.0.0  # The version of the project or build.
account: test  # The account associated with the project.
repository: ""  # The repository URL where the project is stored (can be empty if not specified).
buildDirectory: "./dist"  # Directory where the build artifacts will be output.
logDirectory: "./logs"  # Directory where log files will be stored.

mapping:
  - name: "components"  # Name of the mapping, describing what the mapping represents (e.g., components).
    # This can be any name that makes sense for your project, used for your own convenience.
    relativePath: "install/components"  # Relative path in the project to map files to.
    ifFileExists: "replace"  # Action to take if the file already exists (options: replace, skip, copy-new).
    paths:
      - ./examples/structure/bitrix/components  # List of paths to files that will be mapped.
      - ./examples/structure/local/components

  - name: "templates"
    relativePath: "install/templates"
    ifFileExists: "replace"
    paths:
      - ./examples/structure/bitrix/templates
      - ./examples/structure/local/templates

  - name: "rootFiles"
    relativePath: "."
    ifFileExists: "replace"
    paths:
      - ./examples/structure/simple-file.php

  - name: "testFiles"
    relativePath: "test"
    ifFileExists: "replace"
    paths:
      - ./examples/structure/simple-file.php
        
  - name: "some name"
    relativePath: "some path"
    ifFileExists: "skip"
    paths:
      - some-directory
      - some-filepath
      - etc.

ignore:
  - "**/*.log"  # List of files or patterns to ignore during the build or processing (e.g., log files).
```

### Status

The project is under active development. The API is unstable and subject to change.
