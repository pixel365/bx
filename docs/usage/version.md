# Вывод текущей версии bx

Команда `version` позволяет получить информацию о текущей установленной версии `bx`.

```bash
bx version
```

### Флаги

- `--verbose`, `-n` &mdash; Расширенный вывод.

### Использование

```bash
# v1.2.3
bx version
```

```bash
# Version: v1.2.3
# Commit: xxxxxxxx
# Date: xxx
# Go: 1.24
bx version --verbose 
```

[Исходный код команды](https://github.com/pixel365/bx/blob/main/cmd/version/version.go) на GitHub.