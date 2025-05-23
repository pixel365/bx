# Установка

Вы можете установить последнюю версию `bx` через `go install`

```bash
# установка последней доступной версии
go install github.com/pixel365/bx@latest
```

```bash
# установка конкретной версии
go install github.com/pixel365/bx@v1.2.3
```

Также вы можете [скачать актуальную версию](https://github.com/pixel365/bx/releases/latest) `bx` для нужной платформы.

Если вы используете Fedora, CentOS Stream или другой RPM-based дистрибутив, вы можете установить `bx` из COPR:

```bash
sudo dnf copr enable pixel365/bx
sudo dnf install bx
```

*COPR-пакет содержит готовый бинарник bx, собранный из GitHub-релиза.*

Полный [список](https://github.com/pixel365/bx/releases) версий доступен на GitHub.