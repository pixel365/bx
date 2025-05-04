# Stages

Секция `stages` описывает массив этапов сборки, которые будут использованы в сценариях сборки релиза и/или полной версии модуля, 
а также опционально для произвольных [подкоманд](configuration/run).

- `name` * &mdash; Название этапа. Строковое значение, которое в дальнейшем используется в [builds](configuration/builds) и [run](configuration/run).
- `to` * &mdash; Директория корня дистрибутива модуля.
- `actionIfFileExists` * &mdash; Действие в случае если копируются файлы из разных источников в один путь. Возможные значения:
    - `replace` &mdash; Заменить.
    - `replace_if_newer` &mdash; Заменить если новее.
    - `skip` &mdash; Пропустить.
- `from` * &mdash; Источник из которого нужно скопировать файлы. Это может быть директория, или конечный файл который будет скопирован в `to`.
- `filter` &mdash; Массив шаблонов правил для фильтрации файлов. См. пример ниже.

"*" &mdash; Обязательное поле.

### Примеры

```yaml
name: "module.code"
version: "1.0.0"
account: "test"
buildDirectory: "./dist"
logDirectory: "./logs"

variables:
  longPath: "./long/path/to/some/directory"
  bitrix: "{longPath}/bitrix"
  local: "{longPath}/local"
  
stages:
  - name: "components"
    to: "install/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
```

В данном примере описан базовый, минимальный этап сборки модуля.

По смыслу этого примера видно, 
что этап описывает копирование файлов компонентов в папку `install/components` относительно корня дистрибутива модуля,
из папки `./long/path/to/some/directory/bitrix/components`.

Например, если в папке `./long/path/to/some/directory/bitrix/components` расположен компонент `items.list`,
то содержимое компонента будет скопировано в `./dist/module.code/1.0.0/install/components/items.list`.

```yaml
name: "module.code"
version: "1.0.0"
account: "test"
buildDirectory: "./dist"
logDirectory: "./logs"

variables:
  longPath: "./long/path/to/some/directory"
  bitrix: "{longPath}/bitrix"
  local: "{longPath}/local"
  
stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace_if_newer"
    from:
      - "{bitrix}/components"
      - "{local}/components"
    filter:
      - "**/*.php"
      - "!**/*_test.php"
      - "**/*.js"
      - "**/*.css"
```

В данном примере описан наиболее полный, "продвинутый" этап сборки модуля.

В базовой своей части он повторяет логику из примера выше, но с двумя отличиями:

Во-первых, в массиве поля `from` указано два пути, что означает что если мы копируем содержимое из двух директорий,
и часть его может пересекаться, нужно принять решение как поступить в такой ситуации,
и здесь как раз пригодится поле `actionIfFileExists` значение которого, в данном примере `replace_if_newer`, 
указывает что нужно при пересечении файлов перезаписать тем что новее по дате изменения.

Во-вторых, в этом примере используется поле `filter`, в котором перечислены шаблоны правил того какие файлы будут скопированы.
В данном примере это все `.php`-файлы, кроме тех что заканчиваются на `_test.php`, а также все `.js` и `.css` файлы.
Все другие файлы игнорируются и не будут включены в сборку.

Поле `filter` является хорошим дополнением к секции [ignore](configuration/ignore), 
и позволяет настроить дополнительные правила фильтрации для конкретного этапа сборки.