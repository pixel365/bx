# Основные поля

- `name` * &mdash; Код модуля в формате `developer.module`
- `version` * &mdash; Версия модуля в формате `x.x.x`. Например: `1.2.3`
- `label` &mdash; Метка версии. Возможные значения: `alpha`, `beta`, `stable`. По-умолчанию &mdash; `alpha`.
- `account` * &mdash; Аккаунт (логин) в 1С-Битрикс Маркетплейс, к которому привязан модуль.
- `buildDirectory` * &mdash; Полный или относительный путь до директории в которой будет сохранён дистрибутив модуля.
- `repository` &mdash; Полный или относительный путь до корня репозитория модуля.
- ~~`logDirectory`~~ &mdash; Устарел (см. [настройка лога](configuration/log.md))

"*" &mdash; Обязательное поле.

### Пример

```yaml
name: "module.code"
version: "1.0.0"
label: "stable"
account: "test"
buildDirectory: "./dist"
repository: "."
```

В данном примере заполнены основные поля, а также поле `repository`, 
что позволит в процессе сборки автоматически сгенерировать список изменений для описания версии на основе коммитов ([см. changelog](configuration/changelog))
