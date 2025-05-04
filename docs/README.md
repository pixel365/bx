# BX &mdash; инструмент для сборки модулей для 1С-Битрикс

BX &mdash; позволяет декларативно описывать все этапы сборки проекта, 
а также проверять конфигурацию модуля и развертывать финальный дистрибутив. 
Конфигурации сборки версионируются вместе с проектом, что обеспечивает согласованность и прозрачность изменений 
на протяжении всего процесса разработки и сопровождения модуля.

BX &mdash; делает процесс сборки и публикации релизов простым как никогда раньше.

* [Установка](installation.md)
* [Команды](usage/)
  * [create: Новый модуль](usage/create.md)
  * [check: Проверка конфигурации](usage/check.md)
  * [build: Сборка дистрибутива](usage/build.md)
  * [run: Запуск кастомных команд](usage/run.md)
  * [push: Публикация релиза](usage/push.md)
  * [version: Версия BX](usage/version.md)
* [Настройка](configuration/)
  * [Основные поля](configuration/main.md)
  * [Переменные](configuration/variables.md)
  * [Генерация описания](configuration/changelog.md)
  * [Этапы сборки](configuration/stages.md)
  * [Коллбеки](configuration/callbacks.md)
  * [Сценарии сборки](configuration/builds.md)
  * [Кастомные команды](configuration/run.md)
  * [Настройка исключений](configuration/ignore.md)
  * [Полный пример конфигурации](configuration/example.md)
* [Внести вклад в разработку BX](contribution.md)
