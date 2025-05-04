# Callbacks

Секция `callbacks` описывает массив действий, которые должны быть выполнены до или после выполнения конкретного этапа сборки.

Это может быть либо вызов произвольной команды, или HTTP-запрос. Каждый элемент должен содержать либо `pre`, либо `post`, или и то, и другое.

- `stage` * &mdash; Название этапа сборки.
- `pre` &mdash;
  - `type` ** &mdash; Возможные значения: `command`, `external`.
  - `action` ** &mdash; Для `command` &mdash; это команда, для `external` &mdash; это URL.
  - `method` &mdash; Только для `external`. Например: `GET`, `POST` и т. д.
  - `parameters` Для `command` &mdash; это массив аргументов, для `external` &mdash; это query-параметры.
- `post` &mdash;
    - `type` ** &mdash; Возможные значения: `command`, `external`.
    - `action` ** &mdash; Для `command` &mdash; это команда, для `external` &mdash; это URL.
    - `method` &mdash; Только для `external`. Например: `GET`, `POST` и т. д.
    - `parameters` Для `command` &mdash; это массив аргументов, для `external` &mdash; это query-параметры.
  

"*" &mdash; Обязательное поле.

"**" &mdash; Обязательное поле если присутствует родитель.

### Пример

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
      action: "http://localhost:8080"
      method: "GET"
      parameters:
        - "param1=value1"
        - "param2=value2"
```

В данном примере, для этапа сборки с именем `components` будут применены следующие коллбеки:

- **до** начала этапа сборки, будет выполнена команда `ls -lsa`
- **после** этапа сборки, будет отправлен HTTP-запрос GET: `http://localhost:8008?param1=value1&param2=value2`
