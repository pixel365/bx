# BX

BX is a command-line tool for developers of 1C-Bitrix platform modules. 
The project is currently in active development, and its API may change without backward compatibility.

### Features

- Manage developer accounts
- Maintain a module registry
- Build and prepare a module bundle for publication in the 1C-Bitrix Marketplace

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

#### Help

```shell
bx -h
```

### Example of default module configuration

```yaml
name: "test"
version: "1.0.0"
account: "test"
buildDirectory: "./dist/test"
logDirectory: "./logs/test"

variables:
  structPath: "./examples/structure"
  install: "install"
  bitrix: "{structPath}/bitrix"
  local: "{structPath}/local"
  
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

ignore:
  - "**/*.log"
```

### Configuration Fields

- **name** – The name of the module.
- **version** – The version of the module.
- **account** – The account associated with the module.
- **buildDirectory** – Directory where the build artifacts will be output.
- **logDirectory** – Directory where log files will be stored.
- **variables** (optional) – A set of key-value pairs where both keys and values are strings. These variables can be used in the `stages` section for the `name`, `to`, and `from` fields. Placeholders in curly braces `{}` will be replaced with their corresponding values.

  Example:

  ```yaml
  variables:
    structPath: "./examples/structure"
    install: "install"
    bitrix: "{structPath}/bitrix"
    local: "{structPath}/local"
  ```

  In this case, `{bitrix}` will expand to `./examples/structure/bitrix`, and `{install}` will be replaced with `install` when used in `stages`.

### Stages

The `stages` section defines the steps for copying files. Each stage consists of:

- **name** – The name of the stage. Can use variables.
- **to** – The location where files and directories will be copied, relative to the module's distribution root. Can use variables.
  - For example, if the module's root is `/build/1.2.3`, then setting `to: {install}/components` means files will be placed in `/build/1.2.3/install/components`.
- **from** – The source paths from which files should be copied. Can use variables.
- **actionIfFileExists** – Action to take if the file already exists:
  - `replace` – Overwrite the existing file.
  - `skip` – Skip copying if the file exists.
  - `replace_if_newer` – Overwrite only if the source file is newer.
- **convertTo1251** (optional) - Specifies whether to convert the file contents to windows-1251 encoding. Applies only to *.php files, as well as description.ru. Defaults to false.

#### Example

- **components** – Copies component files to `{install}/components`.
- **templates** – Copies template files to `{install}/templates`.
- **rootFiles** – Copies specific files to the root directory (`.`).
- **testFiles** – Copies test files to `test/`.

**The stage names provided in the examples are for reference only and can be customized as needed.**

### Callbacks

The `callbacks` section allows executing additional actions **before** (`pre`) or **after** (`post`) a specific stage (`stage`).

Each callback consists of:
- **stage** – The name of the stage it is associated with.
- **pre** – Action executed **before** the stage starts.
- **post** – Action executed **after** the stage is completed.

### Supported action types

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
- Before copying `components`, the `ls -lsa` command is executed.
- After copying, a GET request is sent to `http://localhost:80` with query parameters.

### Available parameters

- **type**:
  - `command` – Runs a shell command.
  - `external` – Sends an HTTP request.
- **action**:
  - For `command`: The command to execute.
  - For `external`: The target URL.
- **method** *(for `external` only)* – HTTP method (`GET`, `POST`, etc.).
- **parameters** *(optional)* – List of arguments for commands or query parameters for requests.

### Ignored Files

The `ignore` section defines file patterns to be excluded from processing.  
For example:

```yaml
ignore:
  - "**/*.log"  # Exclude all log files.
```

### Status

The project is under active development. The API is unstable and subject to change.
