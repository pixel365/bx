# BX - Command-Line Tool for 1C-Bitrix Module Development

BX is a command-line tool for developers working on 1C-Bitrix platform modules. It allows you to declaratively define all stages of project build, as well as validate the module configuration and deploy the final distribution. Build configurations are versioned alongside the project, ensuring consistency and traceability of changes throughout the development process. The configuration file allows you to specify:

- **Variables**: paths to required files and directories.
- **Build Stages**: each stage describes which files to copy, to which directory, and how to handle existing files (e.g., replace).
- **Callbacks**: you can specify actions before and after each stage, such as executing commands or sending HTTP requests.
- **File Exclusions**: you can configure which files or directories should be ignored.

Additionally,
BX supports module configuration validation
to ensure it's properly setup and automatic deployment of the final distribution,
simplifying the deployment process.

For example, during one of the build stages, you can copy components and templates from different directories into the final project structure, while also running additional commands before or after the process. This flexibility allows you to manage the build, validation, and deployment processes efficiently, and versioning the configuration ensures that changes won't affect the project's stability.

### Features

- **Declarative Build Configuration**: Define all build stages and actions (e.g., copying files, handling conflicts) in a simple configuration file.
- **Module Configuration Validation**: Ensure the module configuration is correct and complete before building or deploying.
- **Versioned Build Configurations**: Manage build configurations as part of the project, ensuring changes are tracked and consistent.
- **Automatic Deployment**: Deploy the final module distribution to the 1C-Bitrix Marketplace with a single command.
- **Customizable Build Stages**: Define custom stages for different parts of the build process, such as copying components, templates, or files to specific directories.
- **Pre- and Post-Build Callbacks**: Execute additional commands or HTTP requests before or after each build stage for extra automation.

### Installation

```shell
go install github.com/pixel365/bx@latest
```

### Usage

#### Create module

```shell
# Enter a module name via standard input
bx create
```

```shell
# Create a new module (default config)
bx create --name my_module
```

```shell
# Help
bx create -h
```

#### Validate module configuration

```shell
# Choose a module via standard dialog
bx check
```

```shell
# Check the configuration of a module by name
bx check --name my_module
```

```shell
# Check the configuration of a module by file path
bx check --file module-path/config.yaml
```

```shell
# Help
bx check -h
```

#### Build module

```shell
# Choose a module via standard dialog
bx build
```

```shell
# Build a module by name
bx build --name my_module
```

```shell
# Build a module by file path
bx build --file config.yaml
```

```shell
# Override version
bx build --name my_module --version 1.2.3
```

```shell
# Build .last_version
bx build --name my_module --last
```

```shell
# Help
bx build -h
```

#### Push module to Marketplace

```shell
# Choose a module via standard dialog
bx push
```

```shell
# Push a module by name
bx push --name my_module
```

```shell
# Push a module by file path
bx push --file config.yaml
```

```shell
# Override version
bx push --name my_module --version 1.2.3
```

```shell
# Help
bx push -h
```

#### Run custom subcommand

```shell
bx run --cmd customCommand --name my_module
```

#### Help

```shell
bx -h
```

### Configuration Fields

- **name** – The name of the module.
- **version** – The version of the module.
- **account** – The account associated with the module.
- **buildDirectory** – Directory where the build artifacts will be output.
- **logDirectory** – Directory where log files will be stored.
- **repository** *(optional)* – Path to a module repository.
- **variables** (optional) – A set of key-value pairs where both keys and values are strings. These variables can be used in the stages section for the name, to, and from fields. Placeholders in curly braces {} will be replaced with their corresponding values.
- **changelog** *(optional)* – Specifies how to automatically generate a changelog from commit history.
  - `from` and `to` – Define the commit range (tags or specific commits).
    - `type` – Allowed values: `tag`, `commit`.
    - `value` – The specific tag or commit hash.
  - `condition` – Criteria for including or excluding commits.
    - `type` – Allowed values: `include`, `exclude`.
    - `value` – Array of regular expressions for filtering commits.
- **stages** – Defines the file copying and processing stages.
  - `name` – Stage name (supports variables).
  - `to` – Target directory (supports variables).
  - `from` – Source directories or files (supports variables).
  - `actionIfFileExists` – How to handle existing files (`replace`, `skip`, `replace_if_newer`).
  - `convertTo1251` *(optional)* – Converts PHP files and `description.ru` to windows-1251 encoding. Default: `false`.
- **callbacks** *(optional)* – Actions executed before (`pre`) or after (`post`) specific stages.
  - `stage` – Associated stage name.
  - `pre`/`post` – Actions executed before/after the stage.
    - `type` – Allowed values: `command`, `external`.
    - `action` – Command to run or URL for external requests.
    - `method` *(for **`external`**)* – HTTP method (`GET`, `POST`, etc.).
    - `parameters` *(optional)* – Arguments for commands or query parameters for requests.
- **builds** – Defines named build presets to group-specific stages for different types of builds.
  - **Profile name** – The name of the build preset (e.g., `release`, `lastVersion`).
  - **Stages list** – A list of stage names to be included in the build process.
- **ignore** *(optional)* – Patterns for files or directories to exclude from processing.

### Variables explanation

```yaml
variables:
    structPath: "./examples/structure"
    install: "install"
    bitrix: "{structPath}/bitrix"
    local: "{structPath}/local"
```
In this case, {bitrix} will expand to ./examples/structure/bitrix, and {install} will be replaced with install when used in stages.
### Changelog explanation
The `changelog` section defines how to automatically generate a changelog from your project's commit history. It consists of the following fields:

- **from** – Defines the starting point of the commit range.
  - **type** – Indicates the type of the reference point (`tag` or `commit`).
  - **value** – The specific tag name or commit hash.

- **to** – Defines the endpoint of the commit range.
  - **type** – Indicates the type of the reference point (`tag` or `commit`).
  - **value** – The specific tag name or commit hash.

- **condition** (optional) – Criteria for selecting commits to include or exclude.
  - **type** – Filtering mode:
    - `include` – Includes only commits matching the specified patterns.
    - `exclude` – Excludes commits matching the specified patterns.
  - **value** – An array of regular expressions used for filtering commit messages.
- **sort** (optional) – Sorting commits (`asc` or `desc`)

#### Example

```yaml
changelog:
  from:
    type: "tag"
    value: "v1.0.0"
  to:
    type: "tag"
    value: "v2.0.0"
  condition:
    type: "include"
    value:
      - '^feat:([\\W\\w]+)$'
      - '^fix:([\\W\\w]+)$'
  sort: "asc"
```

In this example,
the changelog will include only commits between tags v1.0.0 and v2.0.0 that match patterns starting with feat: or fix:.

**The provided patterns and types are customizable, according to your project's requirements.**

### Stages explanation

The stages section defines the steps for copying files. Each stage consists of:

- **name** – The name of the stage. Can use variables.
- **to** – The location where files and directories will be copied, relative to the module's distribution root. Can use variables.
  - For example, if the module's root is /build/1.2.3, then setting to: {install}/components means files will be placed in /build/1.2.3/install/components.
- **from** – The source paths from which files should be copied. Can use variables.
- **actionIfFileExists** – Action to take if the file already exists:
  - replace – Overwrite the existing file.
  - skip – Skip copying if the file exists.
  - replace_if_newer – Overwrite only if the source file is newer.
- **convertTo1251** (optional) - Specifies whether to convert the file contents to windows-1251 encoding. Applies only to *.php files, as well as description.ru. Defaults to false.

#### Example

```yaml
stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
      - "{local}/components"
```

- **components** – Copies component files to {install}/components.

**The stage names provided in the examples are for reference only and can be customized as needed.**

### Callbacks explanation

The callbacks section allows executing additional actions **before** (pre) or **after** (post) a specific stage (stage).

Each callback consists of:
- **stage** – The name of the stage it is associated with.
- **pre** – Action executed **before** the stage starts.
- **post** – Action executed **after** the stage is completed.

#### Supported action types

- **command** – Executes a shell command.
- **external** – Sends an HTTP request to an external service.

#### Example

```yaml
callbacks:
  - stage: "components"
    pre:
      type: "command"
      action: "ls"
      parameters:
        - "-lsa"
    post:
      type: "external"
      action: "http://localhost:80"
      method: "GET"
      parameters:
        - "param1=value1"
        - "param2=value2"
```

In this example:
- Before copying components, the ls -lsa command is executed.
- After copying, a GET request is sent to http://localhost:80 with query parameters.

#### Available parameters

- **type**:
  - command – Runs a shell command.
  - external – Sends an HTTP request.
- **action**:
  - For command: The command to execute.
  - For external: The target URL.
- **method** *(for external only)* – HTTP method (GET, POST, etc.).
- **parameters** *(optional)* – List of arguments for commands or query parameters for requests.

### Builds explanation

The `builds` section defines named build presets that allow grouping specific stages for different types of builds. Instead of executing all stages, you can specify a subset of them using predefined build profiles.

Each build profile consists of:

- **Profile name** – The name of the build preset (`release`, `lastVersion` (optional)).
- **Stages list** – A list of stage names to be included in the build process.

#### Example

```yaml
builds:
  release:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"
  lastVersion:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"
```

In this example:

- The `release` build includes the `components`, `templates`, `rootFiles`, and `testFiles` stages.
- The `lastVersion` build includes the same stages but can be used to distinguish a specific build variant.

Using the `builds` section allows for greater flexibility
by enabling different build configurations without modifying the core `stages` definition.

### Run explanation

The `run` section provides a way to arbitrarily group any stages from the `stages` section into custom subcommands, allowing them to be executed independently of the main distribution build process. This enables users to define reusable commands tailored to specific workflows without modifying the core build configuration.

Each subcommand consists of:

- **Subcommand name** – The name of the custom command (e.g., `customCommand`).
- **Stages list** – A list of stage names to be executed when the command is run.

#### Example

```yaml
run:
  customCommand:
    - "components"
    - "anotherTestFiles"
```

In this example:

- Running `bx run --cmd customCommand` will execute the components and anotherTestFiles stages.
- This approach allows users to create tailored commands for different workflows without interfering with the primary build process.
- The `run` section is optional — if not specified, no custom subcommands will be available.

This provides a flexible way to execute specific tasks or workflows without triggering a full module build,
making automation and iterative development more efficient.

### Full example of default module configuration
```yaml
name: "test"
version: "1.0.0"
account: "test"
buildDirectory: "./dist/test"
logDirectory: "./logs/test"
repository: "."

variables:
  structPath: "./examples/structure"
  install: "install"
  bitrix: "{structPath}/bitrix"
  local: "{structPath}/local"
  
changelog:
  from:
    type: "tag"
    value: "v1.0.0"
  to:
    type: "tag"
    value: "v2.0.0"
  condition:
    type: "include"
    value:
      - '^feat:([\W\w]+)$'
      - '^fix:([\W\w]+)$'
  sort: "asc"

stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
      - "{local}/components"
  - name: "templates"
    to: "{install}/templates"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/templates"
      - "{local}/templates"
  - name: "rootFiles"
    to: .
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
  - name: "anotherTestFiles"
    to: "another-test"
    actionIfFileExists: "replace"
    from:
      - "./examples/structure/simple-file.php"
    convertTo1251: false
    
callbacks:
  - stage: "components"
    pre:
      type: "command"
      action: "ls"
      parameters:
        - "-lsa"
    post:
      type: "external"
      action: "http://localhost:80"
      method: "GET"
      parameters:
        - "param1=value1"
        - "param2=value2"

builds:
  release:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"
  lastVersion:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"

run:
  customCommand:
    - "components"
    - "anotherTestFiles"

ignore:
  - "**/*.log"
```

### Status

The project is under active development.
The API is unstable and subject to change, so it may lack backward compatibility in future versions.
The current release corresponds to the latest stable branch,
and breaking changes will be documented with each new version.
