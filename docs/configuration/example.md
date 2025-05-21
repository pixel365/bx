# Полный пример конфигурации

```yaml
name: "test"
version: "1.0.0"
label: "stable"
account: "test"
buildDirectory: "./dist/test"
repository: "."

log:
  dir: "./logs"
  maxSize: 10
  maxBackups: 5
  maxAge: 30
  localTime: true
  compress: true

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
  footerTemplate: >
    Warning: This is a template message that is added 
    after the version description or commit list.

stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
      - "{local}/components"
    filter:
      - "**/*.php"
      - "!**/*_test.php"
      - "**/*.js"
      - "**/*.css"
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
