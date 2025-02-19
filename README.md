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
name: test
version: 1.0.0
account: test
repository: ""
buildDirectory: "./dist"
logDirectory: "./logs"
stages:
  - name: "components"
    to: "install/components"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/components
      - ./examples/structure/local/components

  - name: "templates"
    to: "install/templates"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/templates
      - ./examples/structure/local/templates

  - name: "rootFiles"
    to: "."
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

  - name: "another stage name"
    to: "any-path"
    actionIfFileExists: "skip"
    from:
      - ./examples/structure/simple-file.php

ignore:
  - "**/*.log"
```

### Configuration Fields  

- **name** – The name of the module.  
- **version** – The version of the module.  
- **account** – The account associated with the module.  
- **repository** – The repository URL where the project is stored (can be empty if not specified).  
- **buildDirectory** – Directory where the build artifacts will be output.  
- **logDirectory** – Directory where log files will be stored.  

### Stages  

The `stages` section defines the steps for copying files. Each stage consists of:  

- **name** – The name of the stage.  
- **to** – The location where files and directories will be copied, relative to the module's distribution root.
  - For example, if the module's root is `/build/1.2.3`, then setting `to: install/components` means files will be placed in `/build/1.2.3/install/components`.
- **from** – The source paths from which files should be copied.  
- **actionIfFileExists** – Action to take if the file already exists:  
  - `replace` – Overwrite the existing file.  
  - `skip` – Skip copying if the file exists.  
  - `replace_if_newer` – Overwrite only if the source file is newer.  

### Example Stages  

- **components** – Copies component files to `install/components`.  
- **templates** – Copies template files to `install/templates`.  
- **rootFiles** – Copies specific files to the root directory (`.`).  
- **testFiles** – Copies test files to `test/`.  

**The stage names provided in the examples are for reference only and can be customized as needed.**

### Ignored Files  

The `ignore` section defines file patterns to be excluded from processing.  
For example:  

```yaml
ignore:
  - "**/*.log"  # Exclude all log files.
```

### Status

The project is under active development. The API is unstable and subject to change.
