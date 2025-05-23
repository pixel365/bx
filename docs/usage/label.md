# Установить метку версии

Для изменения метки релиза используется команда `label`.

```bash
bx label <alpha|beta|stable> [flags]
```

### Флаги

- `--name`, `-n` &mdash; Код модуля.
- `--file`, `-f` &mdash; Абсолютный или относительный путь до файла, если команда вызывается за пределами местоположения файлов конфигурации по-умолчанию.
- `--version`, `-v` &mdash; Версия модуля. Используется если нужно переопределить версию указанную в файле конфигурации.
- `--password`, `-p` &mdash; Пароль от аккаунта к которому привязан модуль. (См. [пароль в переменной окружения](configuration/password.md))
- `--silent`, `-s` &mdash; "Тихий режим", не выводит статус аутентификации.

### Использование

*Команда `label` автоматически применяется при выполнении команды [push](usage/push.md).*

В отличие от других команд в BX, `label` требует в качестве первого аргумента собственно метку.

Возможные значения: `alpha`, `beta`, `stable`.

Вызов `bx label` без флага `--name` инициирует выбор модуля.

Поиск модулей производится в директории `.bx` текущего контекста.

В случае если в каталоге `.bx` только один модуль &mdash; выбор модуля из списка не будет инициирован,
вместо этого будет выбран единственный доступный модуль.

[Исходный код команды](https://github.com/pixel365/bx/blob/main/cmd/label/label.go) на GitHub.