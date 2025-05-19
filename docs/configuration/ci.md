# CI/CD

Можно настроить автоматическую сборку и публикацию релизов для модуля.

Например, представим задачу сделать так, чтобы при публикации нового тега формата `\d.\d.\d` (1.2.3, 23.54.0 и т.д.), производилась сборка и отправка релиза в Маркетплейс, давайте сделаем это на примере **gitlab**.

Ниже приведены два варианта реализации этого процесса с использованием **GitLab CI/CD**:

Пример с [Shell executor](https://docs.gitlab.com/runner/executors/shell/)

```yaml
# .gitlab-ci.yaml в корне репозитория

stages:
  - build # Этап сборки дистрибутива
  - deploy # Этап публикации в Маркетплейс

variables:
  # Все версии можно найти здесь - https://github.com/pixel365/bx/releases
  BX_VERSION: v1.5.1 # Версия BX
  BX_ARCHIVE: bx_v${BX_VERSION}_linux_amd64.tar.gz
  BX_URL: https://github.com/pixel365/bx/releases/download/${BX_VERSION}/${BX_ARCHIVE}
  BX_BIN: /usr/local/bin/bx

# Создадим скрипт, который проверит установлен-ли BX, если нет - скачаем конкретную версию и установим
.check_bx: &check_bx |
  if ! command -v bx >/dev/null 2>&1; then
    mkdir -p .bx_tmp
    curl -sSL "$BX_URL" -o .bx_tmp/"$BX_ARCHIVE"
    tar -xzf .bx_tmp/"$BX_ARCHIVE" -C .bx_tmp
    chmod +x .bx_tmp/bx
    mv .bx_tmp/bx "$BX_BIN"
    rm -rf .bx_tmp
  fi

build:
  stage: build
  # Проверим BX
  before_script:
    - *check_bx
  script:
    - bx build
  rules:
    - if: '$CI_COMMIT_TAG =~ /^\d+\.\d+\.\d+$/'

deploy:
  stage: deploy
  script:
    - bx push --silent
  rules:
    - if: '$CI_COMMIT_TAG =~ /^\d+\.\d+\.\d+$/'
```

Пример с [Docker executor](https://docs.gitlab.com/runner/executors/docker/)

```yaml
# .gitlab-ci.yaml в корне репозитория

stages:
  - build
  - deploy

variables:
  BX_VERSION: v1.5.1
  BX_ARCHIVE: bx_v${BX_VERSION}_linux_amd64.tar.gz
  BX_URL: https://github.com/pixel365/bx/releases/download/${BX_VERSION}/${BX_ARCHIVE}
  BX_BIN: /usr/local/bin/bx

.check_bx: &check_bx |
  if ! command -v bx >/dev/null 2>&1; then
    mkdir -p .bx_tmp
    curl -sSL "$BX_URL" -o .bx_tmp/"$BX_ARCHIVE"
    tar -xzf .bx_tmp/"$BX_ARCHIVE" -C .bx_tmp
    chmod +x .bx_tmp/bx
    mv .bx_tmp/bx "$BX_BIN"
    rm -rf .bx_tmp
  fi

build:
  stage: build
  image: alpine:3.21
  before_script:
    - apk add --no-cache curl tar
    - *check_bx
  script:
    - bx build
  rules:
    - if: '$CI_COMMIT_TAG =~ /^\d+\.\d+\.\d+$/'

deploy:
  stage: deploy
  image: alpine:3.21
  before_script:
    - apk add --no-cache curl tar
    - *check_bx
  script:
    - bx push
  rules:
    - if: '$CI_COMMIT_TAG =~ /^\d+\.\d+\.\d+$/'
```

> **Обратите внимание:**
> - Эти примеры реализуют базовую автоматизацию.
> - Вы можете расширить pipeline: добавить уведомления, тесты и т.д.
> - Убедитесь что в проекте настроен [GitLab Runner](https://docs.gitlab.com/runner/install/), соответствующий выбранному типу executor’а.

Таким образом, можно сфокусироваться на разработке, а остальное автоматизировать.
