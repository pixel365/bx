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

### Status

The project is under active development.
The API is unstable and subject to change, so it may lack backward compatibility in future versions.
The current release corresponds to the latest stable branch,
and breaking changes will be documented with each new version.
