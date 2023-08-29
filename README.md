# Avito Backend Internship Task
Сервис динамического сегментирования пользователей
<details>
  <summary>Содержание</summary>
  <ol>
    <li><a href="#основные-сведения">Основные сведения</a></li>
    <li><a href="#установка-и-запуск">Установка и запуск</a></li>
    <li><a href="#тесты-и-ci">Тесты и CI</a></li>
  </ol>
</details>

## Основные сведения
Данная программа содержит реализацию сервиса динамического сегментирования пользователей. 

Оригинальное задание -> https://github.com/avito-tech/backend-trainee-assignment-2023 

Ниже приведены основные моменты данной реализации:
 1. Реализовано на **Golang**;
 2. Используются следующие библиотеки:

    2.1. https://github.com/georgysavva/scany -- Работа с БД;
    
    2.2. https://github.com/jackc/pgconn -- Работа с БД;
    
    2.3. https://github.com/jackc/pgx/ -- Работа с БД;
    
    2.4. https://pkg.go.dev/go.uber.org/mock -- Создание stub'ов и mock'ов для тестов.
    
 3. **docker-compose** для поднятия и развертывания среды;
 4. **CI** на **GitHub Actions** и тесты для проверки главного функционала.

Структура проекта выглядит следующим образом:
```markdown
├── .github
|   └── workflows
|       └── tests.yml                   - Файл конфигурации CI
├── build
|   └── Dockerfile                      - Dockerfile для создания образа приложения
├── cmd
|   └── main.go                         - Запускающий файл
├── internal
|   ├── app
|   |   ├── config.go                   - Конфигурация сервера
|   |   ├── handlers.go                 - Методы для обработки запросов
|   |   ├── handlers_test.go            - Тесты для проверки работоспособности
|   |   └── server.go                   - Методы для запуска сервера
|   └── pkg
|       ├── db
|       |   ├── mocks
|       |   |   └── mock.go             - Авто-генерация моков для интерфейса DBops
|       |   ├── sql
|       |   |   └── init.sql            - Файл для создания таблиц в базе данных
|       |   ├── client.go               - Создание экземпляра PostgreSQL
|       |   └── database.go             - Методы по работе с базой данных
|       ├── model
|       |   ├── segment.go              - Структуры для работы с JSON
|       |   └── user_segment.go         - Структуры для работы с JSON
|       ├── repository
|       |   ├── postgresql
|       |   |   ├── segments.go         - Реализация паттерна репозиторий для PostgreSQL
|       |   |   └── user_segment.go     - Реализация паттерна репозиторий для PostgreSQL
|       |   ├── repository.go           - Паттерн репозиторий
|       |   └── structs.go              - Структуры для базы данных
|       └── service
|           ├── mocks
|           |   └── mock.go             - Авто-генерация моков для интерфейса Service
|           ├── postgres.go             - Реализация бизнес-логики для PostgreSQL
|           └── service.go              - Интерфес для работы бизнес-логики
├── .gitignore
├── docker-compose.yml                  - Для поднятия в Docker'е
├── go.mod
├── Makefile
└── README.md
```
Реализация модели хранения для базы данных выглядит не замысловато. Так как в задании не указано, что пользователи должны храниться в этом сервисе, то **подразумевается**, что ID пользователя нам известен заранее и подается из вне как запрос на этот сервис.
![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/286731c1-e35d-457b-b4e1-47aa35db91b3)

Таблица **segments** хранит в себе в качестве ключа название сегмента и его описание. Описание может быть пустым.

Таблица **user_segments** хранит в себе в качестве ключа **id** записи в этой таблице, **user_id** хранит в себе уникальный идентификатор пользователя (который мы знаем заранее), **seg_title** хранит в себе название сегмента, в котором состоит пользователь. **seg_title** не является **FK**! Данная зависимость реализована напрямую из кода, дабы избежать ссылок в самих таблицах.   

Таким образом, чтение, добавление и удаление сегментов у пользователя реализовано по следующему алгоритму:
1. Если данный ID пользователя имеется в таблице **user_segments**, то будут выведены все сегменты из данной таблицы, где поле **user_id** равняется нашему запрашиваемому;
2. Если данный ID пользователя в таблице НЕ имеется, то соответственно количество найденных строк будет 0. Поэтому, в качестве ответа, вернется сообщение о том, что данный пользователь не состоит ни в одном сегменте;
3. Если при добавлении пользователю перечислены несуществующие в таблице **segments** названия сегментов, то записи не будут созданы в таблице **user_segments**, а в качестве ответа вернется массив элементов, которые не были добавлены из-за этого несоответсвия. В остальном в таблицу будут добавлены новые записи, содержащие в себе ID пользователя и сегмент (и так для каждого из перечисленных существующих сегментов);
4. Если при удалении у пользователя перечислены несуществующие в таблице **user_segments** названия сегментов, то подходящие записи не найдутся, а следовательно, в качестве ответа, будет возвращен массив элементов, которые удалить не удалось. В остальном у данного пользователя будут убраны соотвествующие сегменты из данной таблицы.

Для взаимодействия с сервером есть 4 способа:
1. /create-segment -- создать сегмент. Принимает JSON в качестве тела запроса;
2. /delete-segment -- удалить сегмент. Принимает JSON в качестве тела запроса;
3. /user-in-segment -- выполнить операции создания и\или удаления для определенного пользователя. Принимает JSON в качестве тела запроса;
4. /get-user-segments -- получить все сегменты, в которых состоит пользователь. Принимает JSON в качестве тела запроса.

В качестве ответа возвращается JSON, который содержит всегда поле "status" и дополнительные.

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/056722ad-7900-4e34-85c3-b02cabbf06cb)

Сами структуры запросов и ответов в этих файлах:
1. https://github.com/Icerzack/avito-backend-internship/blob/main/internal/pkg/model/segment.go
2. https://github.com/Icerzack/avito-backend-internship/blob/main/internal/pkg/model/user_segment.go

Примеры запросов и ответов в Postman:
1. Создать сегмент (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/a1ddc4cb-2f7b-4b20-a883-ea889e420239)

3. Создать сегмент (некорректный запрос, неправильное наименование полей):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/6cc7f481-becf-4758-84dd-ff9e158021ce)

4. Добавить пользователя в сегмент (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/5652c263-a3d7-426e-84c0-4b0dcdb5f527)

5. Добавить пользователя в сегмент (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/32bd4e3f-edf8-4340-9fb1-dd4908d409d7)

6. Удалить пользователя из сегмента (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/a572e180-07b9-4f70-9908-38234a052d72)

7. Удалить пользователя из сегмента (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/be47a3a8-d9e7-4100-a077-6cf0d2c027d9)

8. Одновременное создание и удаление (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/d6fee204-d9a3-480e-b05b-8f739e858c9b)

9. Получить сегменты пользователя (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/daaa5304-94c3-4d12-bfe0-67c90b4f23d7)

10. Получить сегменты пользователя (корректный запрос, но данный пользователь не содержит сегментов):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/a730fd6e-2d29-4a6d-8485-50855db82e69)

## Установка и запуск

Проект содержит **Makefile**, который имеет следующий вид:

```makefile
.PHONY: compose-db-up
compose-db-up:
	docker-compose build
	docker-compose up -d postgres

.PHONY: compose-db-rm
compose-db-rm:
	docker-compose down

.PHONY: compose-app-up
compose-app-up:
	docker build -f build/Dockerfile -t avito-app .
	docker-compose up -d avito-app

.PHONY: compose-app-rm
compose-app-rm:
	docker-compose down

.PHONY: compose-all-up
compose-all-up: compose-db-up compose-app-up

.PHONY: compose-all-rm
compose-all-rm: compose-app-rm compose-db-rm
```
Для запуска конкретно одного из компонентов приложения нужно, например, прописать следующее:
```bash
$ make compose-db-up
```
Для поднятие всего целиком:
```bash
$ make compose-all-up
```

## Тесты и CI
Для тестирования использовалась библиотека **gomock**. С ее помощью были созданы моковые реализации БД и Сервиса на основании их интерфейсов. Данные реализации применяются в файле **(internal/app/handlers_test.go)**, который тестирует главный функционал приложения (корректная обработка URL и добавление в БД).

В проекте был настроен простой CI для запуска кода на тестах.
Файл конфигурации может быть найден по следующему пути: **(.github/workflows/tests.yml)**
и имеет следующий вид:
```yml
name: tests

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout repository
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.19

    - name: Tests
      run: go test -v avito-backend-internship/internal/app
 ```

<p align="right">(<a href="#основные-сведения">К началу</a>)</p>
